package gate

import (
	"fmt"
	"math/big"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
) (providertypes.GetResponse[mmtypes.Ticker, *big.Int], error) {
	var (
		resolved   = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unresolved = make(map[mmtypes.Ticker]error)
	)

	// The channel must be the tickers channel.
	if Channel(stream.Channel) != ChannelTickers {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved),
			fmt.Errorf("invalid channel %s", stream.Channel)
	}

	// Get the the ticker from the off-chain representation.
	inverted := h.market.Invert()
	market, ok := inverted[stream.Result.CurrencyPair]
	if !ok {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved),
			fmt.Errorf("no currency pair found for symbol %s", stream.Result.CurrencyPair)
	}

	// Parse the price update.
	priceStr := stream.Result.Last
	price, err := math.Float64StringToBigInt(priceStr, market.Ticker.Decimals)
	if err != nil {
		unresolved[market.Ticker] = fmt.Errorf("failed to parse price %s: %w", priceStr, err)
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved), unresolved[market.Ticker]
	}

	resolved[market.Ticker] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse(resolved, unresolved), nil
}
