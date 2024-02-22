package bybit

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

// parseSubscribeResponseMessage parses a subscribe response message. The format of the message
// is defined in the messages.go file. There are two cases that are handled:
//
// 1. Successfully subscribed to the channel. In this case, no further action is required.
// 2. Error message. In this case, we attempt to re-subscribe to the channel.
func (h *WebSocketHandler) parseSubscriptionResponse(resp SubscriptionResponse) ([]handlers.WebsocketEncodedMessage, error) {
	// A response with an event type of subscribe means that we have successfully subscribed to the channel.
	if t := Operation(resp.Op); t == OperationSubscribe && resp.Success {
		h.logger.Info("successfully subscribed to channel", zap.String("connection", resp.ConnID))
		return nil, nil
	}

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
	ticker, ok := h.market.OffChainMap[data.Symbol]
	if !ok {
		return types.NewPriceResponse(resolved, unresolved), fmt.Errorf("unknown ticker %s", data.Symbol)
	}

	// Convert the price to a big.Int.
	price, err := math.Float64StringToBigInt(data.LastPrice, ticker.Decimals)
	if err != nil {
		unresolved[ticker] = fmt.Errorf("failed to convert price to big.Int: %w", err)
		return types.NewPriceResponse(resolved, unresolved), nil
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unresolved), nil
}
