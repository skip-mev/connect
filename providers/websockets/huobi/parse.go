package huobi

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// parseSubscriptionResponse attempts to parse a subscription message.   It returns an error if the message
// cannot be properly parsed.
func (h *WebsocketDataHandler) parseSubscriptionResponse(resp SubscriptionResponse) ([]handlers.WebsocketEncodedMessage, error) {
	if Status(resp.Status) != StatusOk {
		h.logger.Error("unable to create subscription", zap.String("ticker", resp.Subbed))
		// create new message
		msg, err := NewSubscriptionRequest(symbolFromSub(resp.Subbed))
		return []handlers.WebsocketEncodedMessage{msg}, err
	}

	if symbolFromSub(resp.Subbed) == "" {
		return nil, fmt.Errorf("invalid ticker returned")
	}

	h.logger.Debug("successfully subscribed", zap.String("ticker", resp.Subbed))
	return nil, nil
}

// parseTickerStream attempts to parse a ticker stream message.  It returns a providertypes.GetResponse for the
// ticker update.
func (h *WebsocketDataHandler) parseTickerStream(stream TickerStream) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
	)

	ticker := symbolFromSub(stream.Channel)
	if ticker == "" {
		h.logger.Error("incorrectly formatted stream", zap.Any("stream", stream))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("incorrectly formatted stream: %v", stream)
	}

	market, ok := h.cfg.Market.TickerToMarketConfigs[ticker]
	if !ok {
		h.logger.Error("received stream for unknown channel", zap.String("channel", stream.Channel))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("received stream for unknown channel %s", stream.Channel)
	}

	cp := market.CurrencyPair
	price := math.Float64ToBigInt(stream.Tick.LastPrice, cp.Decimals())
	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved), nil
}
