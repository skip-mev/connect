package bitstamp

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
// handle messages received from the Bitstamp websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Bitstamp websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Bitstamp PriceWebSocketDataHandler.
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

// HandleMessage handles a message received from the Bitstamp websocket API. There
// are four types of messages that can be received:
//
//  1. HeartbeatEvent: This is a heartbeat event. This event is sent from the server
//     to the client letting the client know that the connection is still alive.
//  2. ReconnectEvent: This is a reconnect event. This event is sent from the server
//     to the client letting the client know that the server is about to restart.
//  3. SubscriptionSucceededEvent: This is a subscription succeeded event. This event
//     is sent from the server to the client letting the client know that the
//     subscription was successful.
//  4. TradeEvent: This is a trade event. This event is sent from the server to the
//     client letting the client know that a trade has occurred. This event contains
//     the price information for the trade.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp types.PriceResponse
		msg  BaseMessage
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message %w", err)
	}

	switch event := EventType(msg.Event); event {
	case HeartbeatEvent:
		h.logger.Debug("received heartbeat event")
		return resp, nil, nil
	case ReconnectEvent:
		h.logger.Debug("received reconnect event")
		updateMessages, err := NewReconnectRequestMessage()
		return resp, updateMessages, err
	case SubscriptionSucceededEvent:
		h.logger.Debug("received subscription succeeded event")

		var subscriptionMsg SubscriptionResponseMessage
		if err := json.Unmarshal(message, &subscriptionMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscription message %w", err)
		}

		h.logger.Debug("successfully subscribed to channel", zap.String("channel", subscriptionMsg.Channel))
		return resp, nil, nil
	case TradeEvent:
		h.logger.Debug("received ticker event")

		var tickerMsg TickerResponseMessage
		if err := json.Unmarshal(message, &tickerMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %w", err)
		}

		// Parse the price information.
		resp, err := h.parseTickerMessage(tickerMsg)
		return resp, nil, err
	default:
		return resp, nil, fmt.Errorf("unknown event type %s", event)
	}
}

// CreateMessages creates the messages to send to the Bitstamp websocket API. The
// messages are used to subscribe to the live trades channel for the specified tickers.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		instruments = append(instruments, fmt.Sprintf("%s%s", TickerChannel, ticker.GetOffChainTicker()))
		h.cache.Add(ticker)
	}

	return NewSubscriptionRequestMessages(instruments)
}

// HeartBeatMessages is used to create the heartbeat messages to send to the Bitstamp
// websocket API.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return NewHeartbeatRequestMessage()
}

// Copy creates a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger: h.logger,
		ws:     h.ws,
		cache:  types.NewProviderTickers(),
	}
}
