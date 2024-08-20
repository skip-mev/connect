package okx

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the OKX websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the OKX websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new OKX PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	ws config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
	if ws.Name != Name {
		return nil, fmt.Errorf("expected websocket config name %s, got %s", Name, ws.Name)
	}

	if !ws.Enabled {
		return nil, fmt.Errorf("websocket config for %s is not enabled", Name)
	}

	if err := ws.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid websocket config for %s: %w", Name, err)
	}

	return &WebSocketHandler{
		logger: logger,
		ws:     ws,
		cache:  types.NewProviderTickers(),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The OKX
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker response message. This is sent when a ticker update is received from the
//     OKX websocket API.
//
// Heartbeat messages are NOT sent by the OKX websocket. The connection is only closed
// iff no data is received within a 30-second interval or if all subscriptions
// fail. In the case where no data is received within a 30-second interval, the OKX
// will be restarted after the configured restart interval.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp        types.PriceResponse
		baseMessage BaseMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		return resp, nil, err
	}

	eventType := EventType(baseMessage.Event)
	switch {
	case eventType == EventSubscribe || eventType == EventError:
		h.logger.Debug("received subscribe response message")

		var subscribeMessage SubscribeResponseMessage
		if err := json.Unmarshal(message, &subscribeMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
		}

		updateMessage, err := h.parseSubscribeResponseMessage(subscribeMessage)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %w", err)
		}

		return resp, updateMessage, nil
	case eventType == EventTickers:
		h.logger.Debug("received ticker response message")

		var tickerMessage TickersResponseMessage
		if err := json.Unmarshal(message, &tickerMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker response message: %w", err)
		}

		resp, err := h.parseTickerResponseMessage(tickerMessage)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse ticker response message: %w", err)
		}

		return resp, nil, nil
	default:
		return resp, nil, fmt.Errorf("unknown message type %s", baseMessage.Event)
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]SubscriptionTopic, 0)
	for _, ticker := range tickers {
		instruments = append(instruments, SubscriptionTopic{
			Channel:      string(TickersChannel),
			InstrumentID: ticker.GetOffChainTicker(),
		})
		h.cache.Add(ticker)
	}

	return h.NewSubscribeToTickersRequestMessage(instruments)
}

// HeartBeatMessages is not used for okx.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger: h.logger,
		ws:     h.ws,
		cache:  types.NewProviderTickers(),
	}
}
