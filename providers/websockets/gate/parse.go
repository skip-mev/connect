package gate

import (
	"fmt"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

// parseSubscribeResponse attempts to parse a SubscribeResponse to see if it was successful.
func (h *WebSocketHandler) parseSubscribeResponse(
	msg SubscribeResponse,
) ([]handlers.WebsocketEncodedMessage, error) {
	if msg.Error.Message != "" {
		return nil, ErrorCode(msg.Error.Code).Error()
	}

	if Status(msg.Result.Status) != StatusSuccess {
		return nil, fmt.Errorf("subscription was not successful: %s", msg.Result.Status)
	}

	return nil, nil
}

// parseTickerStream attempts to parse a TickerStream and translate it to the corresponding
// ticker update.
func (h *WebSocketHandler) parseTickerStream(
	stream TickerStream,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// The channel must be the tickers channel.
	if Channel(stream.Channel) != ChannelTickers {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("invalid channel %s", stream.Channel)
	}

	// Get the ticker from the off-chain representation.
	ticker, ok := h.cache.FromOffChainTicker(stream.Result.CurrencyPair)
	if !ok {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("no currency pair found for symbol %s", stream.Result.CurrencyPair)
	}

	// Parse the price update.
	priceStr := stream.Result.Last
	price, err := math.Float64StringToBigFloat(priceStr)
	if err != nil {
		wErr := fmt.Errorf("failed to parse price %s: %w", priceStr, err)
		unresolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unresolved), unresolved[ticker]
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unresolved), nil
}
