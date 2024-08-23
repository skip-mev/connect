package kucoin

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the KuCoin websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the KuCoin websocket.
	ws config.WebSocketConfig
	// sequences is a map of currency pair to sequence number.
	sequences map[types.ProviderTicker]int64
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Kucoin PriceWebSocketDataHandler.
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
		logger:    logger,
		ws:        ws,
		sequences: make(map[types.ProviderTicker]int64),
		cache:     types.NewProviderTickers(),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The KuCoin web
// socket expects the client to send a subscribe message within 10 seconds of the
// connection, with a ping message sent every 10 seconds. There are 4 types of messages
// that can be received from the KuCoin websocket:
//
//  1. WelcomeMessage: This is sent by the KuCoin websocket when the connection is
//     established.
//  2. PongMessage: This is sent by the KuCoin websocket in response to a ping message.
//  3. AckMessage: This is sent by the KuCoin websocket in response to a subscribe
//     message.
//  4. Message: This is sent by the KuCoin websocket when a match happens.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp types.PriceResponse
		msg  BaseMessage
	)

	// Determine the type of message received.
	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, err
	}

	switch msgType := MessageType(msg.Type); msgType {
	case WelcomeMessage:
		h.logger.Debug("received welcome message")
		return resp, nil, nil
	case PongMessage:
		h.logger.Debug("received pong message")
		return h.cache.NoPriceChangeResponse(), nil, nil
	case AckMessage:
		h.logger.Debug("received ack message; markets were successfully subscribed to")
		return resp, nil, nil
	case Message:
		h.logger.Debug("received price feed message")

		// Parse the message.
		var ticker TickerResponseMessage
		if err := json.Unmarshal(message, &ticker); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker response message %w", err)
		}

		// Parse the price data from the message.
		resp, err := h.parseTickerResponseMessage(ticker)
		if err != nil {
			return resp, nil, err
		}

		return resp, nil, nil
	default:
		return resp, nil, fmt.Errorf("invalid message type %s", msgType)
	}
}

// CreateMessages is used to create the initial set of subscribe messages to send to the
// KuCoin websocket API. The subscribe messages are created based on the currency pairs
// that are configured for the provider.
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

// HeartBeatMessages is used to create the set of heartbeat messages to send to the KuCoin
// websocket API. Per the KuCoin websocket documentation, the interval between heartbeats
// should be around 10 seconds, however, this is dynamic. As such, the websocket connection
// handler will determine both the credentials and desired ping interval during the pre-dial
// hook.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return NewHeartbeatMessage()
}

// Copy is used to create a copy of the data handler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger:    h.logger,
		ws:        h.ws,
		sequences: make(map[types.ProviderTicker]int64),
		cache:     types.NewProviderTickers(),
	}
}
