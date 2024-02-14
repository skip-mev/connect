package huobi

import (
	"fmt"
	"math/big"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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

	return nil, nil
}

// parseTickerStream attempts to parse a ticker stream message. It returns a providertypes.GetResponse for the
// ticker update.
func (h *WebSocketHandler) parseTickerStream(stream TickerStream) (providertypes.GetResponse[mmtypes.Ticker, *big.Int], error) {
	var (
		resolved   = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unresolved = make(map[mmtypes.Ticker]error)
	)

	ticker := symbolFromSub(stream.Channel)
	if ticker == "" {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved),
			fmt.Errorf("incorrectly formatted stream: %v", stream)
	}

	inverted := h.market.Invert()
	market, ok := inverted[ticker]
	if !ok {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved),
			fmt.Errorf("received stream for unknown channel %s", stream.Channel)
	}

	price := math.Float64ToBigInt(stream.Tick.LastPrice, market.Ticker.Decimals)
	resolved[market.Ticker] = providertypes.NewResult[*big.Int](price, time.Now().UTC())

	return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved), nil
}
