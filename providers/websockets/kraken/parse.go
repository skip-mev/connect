package kraken

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// parseBaseMessage will parse message responses from the Kraken websocket API that are
// not related to price updates. There are three types of messages that are handled by
// this function:
//  1. System status response messages. This is used to check if the Kraken system is online.
//     Usually this is the first message that is received after connecting to the websocket.
//  2. Heartbeat response messages. This is used by the Kraken web socket server to notify
//     the client that the connection is still alive.
//  3. Subscription status response messages. This is used to check if the subscription request
//     was successful. If the subscription request was not successful, the handler will attempt
//     to resubscribe to the market.
func (h *WebSocketDataHandler) parseBaseMessage(message []byte, event Event) ([]handlers.WebsocketEncodedMessage, error) {
	switch event {
	case SystemStatusEvent:
		h.logger.Debug("received system status response message")

		var resp SystemStatusResponseMessage
		if err := json.Unmarshal(message, &resp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal system status response message: %s", err)
		}

		// If the Kraken system is not online, return an error.
		if status := Status(resp.Status); status != OnlineStatus {
			return nil, fmt.Errorf("invalid system status %s", status)
		}

		h.logger.Debug("system status is online")
		return nil, nil
	case HeartbeatEvent:
		h.logger.Debug("received heartbeat response message")
		return nil, nil
	case SubscriptionStatusEvent:
		h.logger.Debug("received subscription status response message")

		var resp SubscribeResponseMessage
		if err := json.Unmarshal(message, &resp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal subscription status response message: %s", err)
		}

		// If the subscription request was successful, return nil. Otherwise, we will attempt to
		// resubscribe to the market.
		switch status := Status(resp.Status); status {
		case SubscribedStatus:
			h.logger.Debug("received successful subscription status response message", zap.String("ticker", resp.Pair))
			return nil, nil
		case ErrorStatus:
			h.logger.Debug(
				"could not successfully subscribe to ticker; attempting to resubscribe",
				zap.String("ticker", resp.Pair),
				zap.String("error", resp.ErrorMessage),
			)

			return NewSubscribeRequestMessage([]string{resp.Pair})
		default:
			return nil, fmt.Errorf("unknown subscription status %s", status)
		}
	default:
		h.logger.Debug("received unknown event", zap.String("event", string(event)))
		return nil, fmt.Errorf("received unknown event %s", event)
	}
}

// parseTickerMessage will parse message responses from the Kraken websocket API that are
// related to price updates. The response message is expected to be in the format of a JSON
// array that contains an update for a single ticker. The response message format can be found
// in messages.go.
func (h *WebSocketDataHandler) parseTickerMessage(
	resp TickerResponseMessage,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unResolved = make(map[oracletypes.CurrencyPair]error)
	)

	// We will only parse messages from the ticker channel.
	if ch := Channel(resp.ChannelName); ch != TickerChannel {
		h.logger.Debug("received price update for unknown channel", zap.String("channel", string(ch)))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("invalid channel %s", ch)
	}

	// Get the currency pair from the instrument.
	h.logger.Debug("received price update", zap.String("instrument", resp.Pair))
	market, ok := h.cfg.Market.TickerToMarketConfigs[resp.Pair]
	if !ok {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("no currency pair found for instrument %s", resp.Pair)
	}

	// Ensure that the length of the price update is valid.
	if len(resp.TickerData.VolumeWeightedAveragePrice) != ExpectedVolumeWeightedAveragePriceLength {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("invalid price update length %d", len(resp.TickerData.VolumeWeightedAveragePrice))
	}

	// Parse the price update.
	cp := market.CurrencyPair
	priceStr := resp.TickerData.VolumeWeightedAveragePrice[TodayPriceIndex]
	price, err := math.Float64StringToBigInt(priceStr, cp.Decimals())
	if err != nil {
		unResolved[cp] = fmt.Errorf("failed to parse price %s: %s", priceStr, err)
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), unResolved[cp]
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
