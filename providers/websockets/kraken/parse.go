package kraken

import (
	"encoding/json"
	"fmt"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

// parseBaseMessage will parse message responses from the Kraken websocket API that are
// not related to price updates. There are three types of messages that are handled by
// this function:
//  1. System status response messages. This is used to check if the Kraken system is online.
//     Usually this is the first message that is received after connecting to the websocket.
//  2. Heartbeat response messages. This is used by the Kraken websocket server to notify
//     the client that the connection is still alive.
//  3. Subscription status response messages. This is used to check if the subscription request
//     was successful. If the subscription request was not successful, the handler will attempt
//     to resubscribe to the market.
func (h *WebSocketHandler) parseBaseMessage(message []byte, event Event) ([]handlers.WebsocketEncodedMessage, error) {
	switch event {
	case SystemStatusEvent:
		h.logger.Debug("received system status response message")

		var resp SystemStatusResponseMessage
		if err := json.Unmarshal(message, &resp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal system status response message: %w", err)
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
			return nil, fmt.Errorf("failed to unmarshal subscription status response message: %w", err)
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

			return h.NewSubscribeRequestMessage([]string{resp.Pair})
		default:
			return nil, fmt.Errorf("unknown subscription status %s", status)
		}
	default:
		return nil, fmt.Errorf("received unknown event %s", event)
	}
}

// parseTickerMessage will parse message responses from the Kraken websocket API that are
// related to price updates. The response message is expected to be in the format of a JSON
// array that contains an update for a single ticker. The response message format can be found
// in messages.go.
func (h *WebSocketHandler) parseTickerMessage(
	resp TickerResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// We will only parse messages from the ticker channel.
	if ch := Channel(resp.ChannelName); ch != TickerChannel {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("invalid channel %s", ch)
	}

	// Get the ticker from the instrument.
	ticker, ok := h.cache.FromOffChainTicker(resp.Pair)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("no ticker found for instrument %s", resp.Pair)
	}

	// Ensure that the length of the price update is valid.
	if len(resp.TickerData.VolumeWeightedAveragePrice) != ExpectedVolumeWeightedAveragePriceLength {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("invalid price update length %d", len(resp.TickerData.VolumeWeightedAveragePrice))
	}

	// Parse the price update.
	priceStr := resp.TickerData.VolumeWeightedAveragePrice[TodayPriceIndex]
	price, err := math.Float64StringToBigFloat(priceStr)
	if err != nil {
		wErr := fmt.Errorf("failed to parse price %s: %w", priceStr, err)
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), unResolved[ticker]
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}

// DecodeTickerResponseMessage decodes a ticker response message.
func DecodeTickerResponseMessage(message []byte) (TickerResponseMessage, error) {
	var rawResponse []json.RawMessage
	if err := json.Unmarshal(message, &rawResponse); err != nil {
		return TickerResponseMessage{}, err
	}

	if len(rawResponse) != ExpectedTickerResponseMessageLength {
		return TickerResponseMessage{}, fmt.Errorf(
			"invalid ticker response message; expected length %d, got %d", ExpectedTickerResponseMessageLength, len(rawResponse),
		)
	}

	var response TickerResponseMessage
	if err := json.Unmarshal(rawResponse[ChannelIDIndex], &response.ChannelID); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[TickerDataIndex], &response.TickerData); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[ChannelNameIndex], &response.ChannelName); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[PairIndex], &response.Pair); err != nil {
		return TickerResponseMessage{}, err
	}

	return response, nil
}
