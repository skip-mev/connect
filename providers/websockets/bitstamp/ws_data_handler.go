package bitstamp

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
// handle messages received from the Coinbase websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the Coinbase websocket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Bitstamp.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	cfg config.ProviderConfig,
) (handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %w", err)
	}

	if !cfg.WebSocket.Enabled {
		return nil, fmt.Errorf("websocket is not enabled for provider %s", cfg.Name)
	}

	if cfg.Name != Name {
		return nil, fmt.Errorf("invalid provider name %s", cfg.Name)
	}

	return &WebSocketDataHandler{
		cfg:    cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
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
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
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
		h.logger.Debug("received trade event")

		var tickerMsg TickerResponseMessage
		if err := json.Unmarshal(message, &tickerMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %w", err)
		}

		// Parse the price information.
		resp, err := h.parseTickerMessage(tickerMsg)
		return resp, nil, err
	default:
		h.logger.Debug("received unknown event", zap.String("event", string(event)))
		return resp, nil, fmt.Errorf("unknown event type %s", event)
	}
}

// CreateMessages creates the messages to send to the Bitstamp websocket API. The
// messages are used to subscribe to the live trades channel for the specified currency
// pairs.
func (h *WebSocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.Ticker()]
		if !ok {
			return nil, fmt.Errorf("currency pair not found in market configs %s", cp.Ticker())
		}

		instruments = append(instruments, fmt.Sprintf("%s%s", TickerChannel, market.Ticker))
	}

	return NewSubscriptionRequestMessages(instruments)
}

// HeartBeatMessages is used to create the heartbeat messages to send to the Bitstamp
// websocket API.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return NewHeartbeatRequestMessage()
}
