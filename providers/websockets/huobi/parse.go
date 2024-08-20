package huobi

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

// parseSubscriptionResponse attempts to parse a subscription message. It returns an error if the message
// cannot be properly parsed.
func (h *WebSocketHandler) parseSubscriptionResponse(resp SubscriptionResponse) ([]handlers.WebsocketEncodedMessage, error) {
	if Status(resp.Status) != StatusOk {
		msg, err := NewSubscriptionRequest(symbolFromSub(resp.Subbed))
		return []handlers.WebsocketEncodedMessage{msg}, err
	}

	if symbolFromSub(resp.Subbed) == "" {
		return nil, fmt.Errorf("invalid ticker returned")
	}

	h.logger.Debug("successfully subscribed", zap.String("ticker", resp.Subbed))
	return nil, nil
}

// parseTickerStream attempts to parse a ticker stream message. It returns a providertypes.GetResponse for the
// ticker update.
func (h *WebSocketHandler) parseTickerStream(stream TickerStream) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	offChainTicker := symbolFromSub(stream.Channel)
	if offChainTicker == "" {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("incorrectly formatted stream: %v", stream)
	}

	ticker, ok := h.cache.FromOffChainTicker(offChainTicker)
	if !ok {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("received stream for unknown channel %s", stream.Channel)
	}

	price := big.NewFloat(stream.Tick.LastPrice)
	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())

	return types.NewPriceResponse(resolved, unresolved), nil
}
