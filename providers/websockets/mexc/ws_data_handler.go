package mexc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the exchange.
	Name = "mexc"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebSocketDataHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the MEXC websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the MEXC web socket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for MEXC.
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

// HandleMessage is used to handle a message received from the data provider. This is called
// when a message is received from the data provider. There are three types of messages that
// can be received from the data provider:
//
// 1. A message that confirms that the client has successfully subscribed to a channel.
// 2. A message that confirms that the client has successfully pinged the server.
// 3. A message that contains the latest price for a currency pair.
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp      providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		msg       BaseMessage
		tickerMsg TickerResponseMessage
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message %s", err)
	}

	// If the base message is empty, we assume it is a price message.
	if msg.IsEmpty() {
		if err := json.Unmarshal(message, &tickerMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %s", err)
		}

		// Parse the ticker message.
		resp, err := h.parseTickerResponseMessage(tickerMsg)
		return resp, nil, err
	}

	// Otherwise, we assume it is a subscription or pong message.
	switch {
	case strings.HasPrefix(msg.Message, string(MiniTickerChannel)):
		h.logger.Debug("subscribed to ticker channel", zap.String("instruments", msg.Message))
		return resp, nil, nil
	case MethodType(msg.Message) == PongMethod:
		h.logger.Debug("received pong message")
		return resp, nil, nil
	default:
		h.logger.Debug("received unknown message type", zap.String("type", msg.Message))
		return resp, nil, fmt.Errorf("invalid message type %s", msg.Message)
	}
}

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given currency pairs. This is called when the connection to the data
// provider is first established.
func (h *WebSocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	if len(cps) > MaxSubscriptionsPerConnection {
		return nil, fmt.Errorf("cannot subscribe to more than %d currency pairs per connection", MaxSubscriptionsPerConnection)
	}

	instruments := make([]string, 0)
	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			return nil, fmt.Errorf("currency pair %s not found in market configs", cp.String())
		}

		mexcTicker := fmt.Sprintf("%s%s%s", string(MiniTickerChannel), strings.ToUpper(market.Ticker), "@UTC+8")
		instruments = append(instruments, mexcTicker)
	}

	return NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is used by the MEXC handler to send heart beat messages to the data provider.
// This is used to keep the connection alive when no messages are being sent from the data provider.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return NewPingRequestMessage()
}
