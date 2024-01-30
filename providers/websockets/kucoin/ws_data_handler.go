package kucoin

import (
	"encoding/json"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebSocketDataHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the KuCoin websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the KuCoin web socket API.
	cfg config.ProviderConfig

	// sequences is a map of currency pair to sequence number.
	sequences map[oracletypes.CurrencyPair]int64
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for KuCoin.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	cfg config.ProviderConfig,
) (handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %s", err)
	}

	if !cfg.WebSocket.Enabled {
		return nil, fmt.Errorf("web socket is not enabled for provider %s", cfg.Name)
	}

	if cfg.Name != Name {
		return nil, fmt.Errorf("invalid provider name %s", cfg.Name)
	}

	return &WebSocketDataHandler{
		cfg:       cfg,
		logger:    logger.With(zap.String("web_socket_data_handler", Name)),
		sequences: make(map[oracletypes.CurrencyPair]int64),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The KuCoin web
// socket expects the client to send a subscribe message within 10 seconds of the
// connection, with a ping message sent every 10 seconds. There are 4 types of messages
// that can be received from the KuCoin web socket:
//
//  1. WelcomeMessage: This is sent by the KuCoin web socket when the connection is
//     established.
//  2. PongMessage: This is sent by the KuCoin web socket in response to a ping message.
//  3. AckMessage: This is sent by the KuCoin web socket in response to a subscribe
//     message.
//  4. Message: This is sent by the KuCoin web socket when a match happens.
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
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
		return resp, nil, nil
	case AckMessage:
		h.logger.Debug("received ack message; markets were successfully subscribed to")
		return resp, nil, nil
	case Message:
		h.logger.Debug("received price feed message")

		// Parse the message.
		var ticker TickerResponseMessage
		if err := json.Unmarshal(message, &ticker); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker response message %s", err)
		}

		// Parse the price data from the message.
		resp, err := h.parseTickerResponseMessage(ticker)
		if err != nil {
			return resp, nil, err
		}

		return resp, nil, nil
	default:
		h.logger.Debug("received invalid message type", zap.String("message_type", string(msgType)))
		return resp, nil, fmt.Errorf("invalid message type %s", msgType)
	}
}

// CreateMessages is used to create the initial set of subscribe messages to send to the
// KuCoin web socket API. The subscribe messages are created based on the currency pairs
// that are configured for the provider.
func (h *WebSocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Warn("currency pair not found in market configs", zap.String("currency_pair", cp.String()))
			continue
		}

		instruments = append(instruments, market.Ticker)
	}

	return NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is used to create the set of heartbeat messages to send to the KuCoin
// websocket API. Per the KuCoin websocket documentation, the interval between heartbeats
// should be around 10 seconds, however, this is dynamic. As such, the web socket connection
// handler will determine both the credentials and desired ping interval during the pre-dial
// hook.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return NewHeartbeatMessage()
}
