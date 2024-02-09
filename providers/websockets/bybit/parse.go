package bybit

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// parseSubscribeResponseMessage parses a subscribe response message. The format of the message
// is defined in the messages.go file. There are two cases that are handled:
//
// 1. Successfully subscribed to the channel. In this case, no further action is required.
// 2. Error message. In this case, we attempt to re-subscribe to the channel.
func (h *WebsocketDataHandler) parseSubscriptionResponse(resp SubscriptionResponse) ([]handlers.WebsocketEncodedMessage, error) {
	// A response with an event type of subscribe means that we have successfully subscribed to the channel.
	if t := Operation(resp.Op); t == OperationSubscribe && resp.Success {
		h.logger.Info("successfully subscribed to channel", zap.String("connection", resp.ConnID))
		return nil, nil
	}

	if t := Operation(resp.Op); t == OperationSubscribe && !resp.Success {
		h.logger.Error("received error message", zap.String("message", resp.RetMsg))

		// TODO resubscribe ?
		return nil, nil
	}

	h.logger.Error("unable to parse message", zap.Any("message", resp))
	return nil, fmt.Errorf("unable to parse message")
}

// parseTickerUpdate parses a ticker update message. The format of the message is defined
// in the messages.go file. This message contains the latest price data for a set of pairs.
func (h *WebsocketDataHandler) parseTickerUpdate(
	resp TickerUpdateMessage,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[slinkytypes.CurrencyPair]error)
	)

	// The topic must be the tickers topic.
	if !strings.Contains(resp.Topic, string(TickerChannel)) {
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("invalid topic %s", resp.Topic)
	}

	data := resp.Data
	// Iterate through all the tickers and add them to the response.
	market, ok := h.cfg.Market.TickerToMarketConfigs[data.Symbol]
	if !ok {
		h.logger.Debug("currency pair not found for symbol ID", zap.String("symbol", data.Symbol))
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved), nil
	}

	cp := market.CurrencyPair

	// Convert the price to a big.Int.
	price, err := math.Float64StringToBigInt(data.LastPrice, cp.Decimals())
	if err != nil {
		h.logger.Error("failed to convert price to big.Int", zap.Error(err))
		unresolved[cp] = fmt.Errorf("failed to convert price to big.Int: %w", err)
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved), nil
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved), nil
}
