package okx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

const (
	// ExpectedErrorPrefix is the prefix of an error message that is returned by the OKX API.
	// Specifically, this is the prefix of the error message that is returned when the user
	// attempts to subscribe to a channel but could not be subscribed.
	ExpectedErrorPrefix = "Invalid request: "

	// ExpectedErrorElements is the number of elements that are expected in the error message.
	ExpectedErrorElements = 2
)

// parseSubscribeResponseMessage parses a subscribe response message. The format of the message
// is defined in the messages.go file. There are two cases that are handled:
//
// 1. Successfully subscribed to the channel. In this case, no further action is required.
// 2. Error message. In this case, we attempt to re-subscribe to the channel.
func (h *WebSocketHandler) parseSubscribeResponseMessage(resp SubscribeResponseMessage) ([]handlers.WebsocketEncodedMessage, error) {
	// A response with an event type of subscribe means that we have successfully subscribed to the channel.
	if t := EventType(resp.Event); t == EventSubscribe {
		h.logger.Debug("successfully subscribed to channel", zap.String("instrument", resp.Arguments.InstrumentID))
		return nil, nil
	}

	// Attempt to re-subscribe to the channel.
	// Format of the message is:
	//  ...
	//	"msg": "Invalid request: {\"op\": \"subscribe\", \"args\":[{ \"channel\" : \"index-tickers\", \"instId\" : \"BTC-USDT\"}]}",
	//  ...
	//
	// The message is an exact copy of the request message, so we can just unmarshal it and re-subscribe.
	h.logger.Debug("received error message", zap.String("message", resp.Message), zap.String("code", resp.Code))
	jsonString := strings.Split(resp.Message, ExpectedErrorPrefix)
	if len(jsonString) != ExpectedErrorElements {
		return nil, fmt.Errorf("unable to parse subscription message from message: %s", resp.Message)
	}

	// Attempt to unmarshal the request.
	var request SubscribeRequestMessage
	if err := json.Unmarshal([]byte(jsonString[1]), &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// Re-subscribe to the channel.
	h.logger.Debug("re-subscribing to channel", zap.Any("instrument", request.Arguments))
	bz, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return []handlers.WebsocketEncodedMessage{bz}, nil
}

// parseTickerResponseMessage parses a ticker response message. The format of the message is defined
// in the messages.go file. This message contains the latest price data for a set of instruments.
func (h *WebSocketHandler) parseTickerResponseMessage(
	resp TickersResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// The channel must be the index tickers channel.
	if Channel(resp.Arguments.Channel) != TickersChannel {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("invalid channel %s", resp.Arguments.Channel)
	}

	// Iterate through all tickers and add them to the response.
	for _, instrument := range resp.Data {
		ticker, ok := h.cache.FromOffChainTicker(instrument.ID)
		if !ok {
			h.logger.Debug("ticker not found for instrument ID", zap.String("instrument_id", instrument.ID))
			continue
		}

		// Convert the price to a big.Float.
		price, err := math.Float64StringToBigFloat(instrument.LastPrice)
		if err != nil {
			wErr := fmt.Errorf("failed to convert price to big.Float: %w", err)
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
			}
			continue
		}

		resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	}

	return types.NewPriceResponse(resolved, unresolved), nil
}
