package bybit

import (
	"fmt"
	"strings"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

// parseSubscriptionResponse parses a subscribe response message. The format of the message
// is defined in the messages.go file. There are two cases that are handled:
//
// 1. Successfully subscribed to the channel. In this case, no further action is required.
// 2. Error message. In this case, we attempt to re-subscribe to the channel.
func (h *WebSocketHandler) parseSubscriptionResponse(resp SubscriptionResponse) ([]handlers.WebsocketEncodedMessage, error) {
	// A response with an event type of subscribe means that we have successfully subscribed to the channel.
	if t := Operation(resp.Op); t == OperationSubscribe && resp.Success {
		h.logger.Debug("successfully subscribed to channel", zap.String("connection", resp.ConnID))
		return nil, nil
	}

	// TODO(david): Add a retry mechanism here.
	if t := Operation(resp.Op); t == OperationSubscribe && !resp.Success {
		return nil, fmt.Errorf("received error message: %s", resp.RetMsg)
	}

	return nil, fmt.Errorf("unable to parse message")
}

// parseTickerUpdate parses a ticker update message. The format of the message is defined
// in the messages.go file. This message contains the latest price data for a set of pairs.
func (h *WebSocketHandler) parseTickerUpdate(
	resp TickerUpdateMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// The topic must be the tickers topic.
	if !strings.Contains(resp.Topic, string(TickerChannel)) {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("invalid topic %s", resp.Topic)
	}

	// Iterate through all the tickers and add them to the response.
	data := resp.Data
	ticker, ok := h.cache.FromOffChainTicker(data.Symbol)
	if !ok {
		return types.NewPriceResponse(resolved, unresolved), fmt.Errorf("unknown ticker %s", data.Symbol)
	}

	// Convert the price to a big.Float.
	price, err := math.Float64StringToBigFloat(data.LastPrice)
	if err != nil {
		wErr := fmt.Errorf("failed to convert price to big.Float: %w", err)
		unresolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unresolved), nil
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unresolved), nil
}
