package coinbase

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

	// config is the config for the Coinbase web socket API.
	cfg config.ProviderConfig

	// Sequence is the current sequence number for the Coinbase web socket API per currency pair.
	sequence map[oracletypes.CurrencyPair]int64
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Coinbase.
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
		cfg:      cfg,
		logger:   logger.With(zap.String("web_socket_data_handler", Name)),
		sequence: make(map[oracletypes.CurrencyPair]int64),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Coinbase web
// socket expects the client to send a subscribe message within 5 seconds of the initial connection.
// Otherwise, the connection will be closed. There are two types of messages that can be received
// from the Coinbase web socket API:
//
//  1. SubscriptionsMessage: This is sent by the Coinbase web socket API after a subscribe message
//     is sent. This message contains the list of channels that were successfully subscribed to.
//  2. TickerMessage: This is sent by the Coinbase web socket API when a match happens. This message
//     contains the price of the currency pair.
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

	switch MessageType(msg.Type) {
	case SubscriptionsMessage:
		h.logger.Debug("received subscriptions message")

		var subscriptionsMessage SubscribeResponseMessage
		if err := json.Unmarshal(message, &subscriptionsMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscriptions message %s", err)
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
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %s", err)
		}

		resp, err := h.parseTickerResponseMessage(tickerMessage)
		return resp, nil, err
	default:
		h.logger.Debug("received unknown message type", zap.String("type", msg.Type))
		return resp, nil, fmt.Errorf("invalid message type %s", msg.Type)
	}
}

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given currency pairs. This is called when the connection to the data
// provider is first established.
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

		instruments = append(instruments, market.Ticker)
	}

	return NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is not used for Coinbase.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
