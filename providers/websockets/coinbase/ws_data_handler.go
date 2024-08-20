package coinbase

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
// handle messages received from the Coinbase websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Coinbase websocket.
	ws config.WebSocketConfig
	// sequence is the current trade sequence number for the Coinbase websocket API per currency pair.
	sequence map[types.ProviderTicker]int64
	// tradeIDs is the current trade ID for the Coinbase websocket API per currency pair.
	tradeIDs map[types.ProviderTicker]int64
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Coinbase PriceWebSocketDataHandler.
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
		logger:   logger,
		ws:       ws,
		sequence: make(map[types.ProviderTicker]int64),
		tradeIDs: make(map[types.ProviderTicker]int64),
		cache:    types.NewProviderTickers(),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Coinbase web
// socket expects the client to send a subscribe message within 5 seconds of the initial connection.
// Otherwise, the connection will be closed. There are two types of messages that can be received
// from the Coinbase websocket API:
//
//  1. SubscriptionsMessage: This is sent by the Coinbase websocket API after a subscribe message
//     is sent. This message contains the list of channels that were successfully subscribed to.
//  2. TickerMessage: This is sent by the Coinbase websocket API when a match happens. This message
//     contains the price of the ticker.
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

	switch MessageType(msg.Type) {
	case SubscriptionsMessage:
		h.logger.Debug("received subscriptions message")

		var subscriptionsMessage SubscribeResponseMessage
		if err := json.Unmarshal(message, &subscriptionsMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscriptions message %w", err)
		}

		// Log the channels that were successfully subscribed to.
		for _, channel := range subscriptionsMessage.Channels {
			for _, instrument := range channel.Instruments {
				h.logger.Debug("subscribed to ticker channel", zap.String("instrument", instrument), zap.String("channel", channel.Name))
			}
		}

		return resp, nil, nil
	case TickerMessage:
		h.logger.Debug("received ticker message")

		var tickerMessage TickerResponseMessage
		if err := json.Unmarshal(message, &tickerMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %w", err)
		}

		resp, err := h.parseTickerResponseMessage(tickerMessage)
		return resp, nil, err
	case HeartbeatMessage:
		h.logger.Debug("received product heartbeat message")

		var heartbeatMessage HeartbeatResponseMessage
		if err := json.Unmarshal(message, &heartbeatMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal heartbeat message %w", err)
		}

		resp, err := h.parseHeartbeatResponseMessage(heartbeatMessage)
		return resp, nil, err
	default:
		return resp, nil, fmt.Errorf("invalid message type %s", msg.Type)
	}
}

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given tickers. This is called when the connection to the data provider is
// first established.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		instruments = append(instruments, ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	return h.NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is not used for Coinbase.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger:   h.logger,
		ws:       h.ws,
		sequence: make(map[types.ProviderTicker]int64),
		tradeIDs: make(map[types.ProviderTicker]int64),
		cache:    types.NewProviderTickers(),
	}
}
