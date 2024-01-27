package gate

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

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Gate.io websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the OKX web socket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Gate.io
// from a given provider configuration.
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

	return &WebsocketDataHandler{
		cfg:    cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Gate.io
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker stream message. This is sent when a ticker update is received from the
//     Gate.io web socket API.
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp         providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		baseMessage  BaseMessage
		subResponse  SubscribeResponse
		tickerStream TickerStream
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		h.logger.Debug("unable to unmarshal message into base message", zap.Error(err))
		return resp, nil, err
	}

	switch Event(baseMessage.Event) {
	case EventSubscribe:
		if err := json.Unmarshal(message, &subResponse); err != nil {
			h.logger.Debug("unable to unmarshal message into subscribe response", zap.Error(err))
			return resp, nil, err
		}

		// handle subscription
		updateMsg, err := h.parseSubscribeResponse(subResponse)
		return resp, updateMsg, err

	case EventUpdate:
		if err := json.Unmarshal(message, &tickerStream); err != nil {
			h.logger.Debug("unable to unmarshal message into ticker stream", zap.Error(err))
			return resp, nil, err
		}

		// update pair info
		resp, err := h.parseTickerStream(tickerStream)
		return resp, nil, err

	default:
		h.logger.Debug("received unknown message", zap.String("message", string(message)))
		return resp, nil, fmt.Errorf("unknown message type %s", baseMessage.Event)
	}

	return resp, nil, nil
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the tickers channel - which supports spot markets.
func (h *WebsocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	symbols := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("market not found for currency pair", zap.String("currency_pair", cp.String()))
			continue
		}

		symbols = append(symbols, market.Ticker)
	}

	h.logger.Debug("subscribing", zap.Any("symbols", symbols))
	return NewSubscribeRequest(symbols)
}

// HeartBeatMessages is not used for Gate.io.
func (h *WebsocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
