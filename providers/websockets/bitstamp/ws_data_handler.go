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

const (
	// Name is the name of the exchange.
	Name = "bitstamp"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebSocketDataHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Coinbase websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the Coinbase web socket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Bitstamp.
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
		cfg:    cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage handles a message received from the Bitstamp websocket API. There
// are three types of messages that can be received:
//
//  1. HeartbeatEvent: This is a heartbeat event. This event is sent from the server
//     to the client letting the client know that the connection is still alive.
//  2. TickerEvent: This is a ticker event. This event is sent from the server to
//     the client when a ticker update is available.
//  3. ReconnectEvent: This is a reconnect event. This event is sent from the server
//     to the client when the server is about to disconnect the client. The client
//     has a few seconds to reconnect. If the client does not reconnect, the client
//     will be disconnected.
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		msg  BaseMessage
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message %s", err)
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

		var subcriptionMsg SubscriptionResponseMessage
		if err := json.Unmarshal(message, &subcriptionMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscription message %s", err)
		}

		h.logger.Debug("successfully subscribed to channel", zap.String("channel", subcriptionMsg.Channel))
		return resp, nil, nil
	case TradeEvent:
		var tickerMsg TickerResponseMessage
		if err := json.Unmarshal(message, &tickerMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %s", err)
		}

		// Parse the price information.
		resp, err := h.parseTickerMessage(tickerMsg)
		return resp, nil, err
	default:
		return resp, nil, fmt.Errorf("unknown event type %s", event)
	}
}

func (h *WebSocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("currency pair not found in market configs", zap.String("currency_pair", cp.String()))
			continue
		}

		instruments = append(instruments, fmt.Sprintf("%s%s", TickerChannel, market.Ticker))
	}

	return NewSubscriptionRequestMessages(instruments)
}

// HeartBeatMessages is not used for Bitstamp.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
