package bybit

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
	// Name is the name of the OKX provider.
	Name = "bybit"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the ByBit websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the ByBit web socket API.
	config Config
}

// NewWebSocketDataHandlerFromConfig returns a new WebSocketDataHandler implementation for ByBit
// from a given provider configuration.
func NewWebSocketDataHandlerFromConfig(
	logger *zap.Logger,
	providerCfg config.ProviderConfig,
) (handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int], error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("invalid provider name %s", providerCfg.Name)
	}

	cfg, err := ReadConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %s", providerCfg.Path, err)
	}

	return &WebsocketDataHandler{
		config: cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for ByBit.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	cfg Config,
) (handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}

	return &WebsocketDataHandler{
		config: cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The ByBit
// provider sends three types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker update message. This is sent when a ticker update is received from the
//     ByBit web socket API.
//  3. Heartbeat update messages.  This should be sent every 20 seconds to ensure the
//     connection remains open.
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []byte, error) {
	var (
		resp         providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		baseResponse BaseResponse
		update       TickerUpdateMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseResponse); err != nil {
		h.logger.Error("failed to unmarshal subscribe response message", zap.Error(err))
		return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %s", err)

	}

	opType := Operation(baseResponse.Op)
	switch {
	case opType == OperationSubscribe:
		h.logger.Debug("received subscribe response message")

		var subscribeMessage SubscriptionResponse
		if err := json.Unmarshal(message, &subscribeMessage); err != nil {
			h.logger.Error("failed to unmarshal subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %s", err)
		}

		updateMessage, err := h.parsSubscriptionResponse(subscribeMessage)
		if err != nil {
			h.logger.Error("failed to parse subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %s", err)
		}

		return resp, updateMessage, nil
	case opType == OperationPong:
		h.logger.Debug("received pong response message")

		return resp, nil, nil
	default:
		// if the message is not a base message, then it must be a stream response
		if err := json.Unmarshal(message, &update); err != nil {
			h.logger.Debug("unable to recognize message", zap.Error(err), zap.Binary("message", message))
			return resp, nil, err
		}

		// parse
		resp, err := h.parseTickerUpdate(update)
		if err != nil {
			h.logger.Error("failed to parse ticker update message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse ticker update message: %s", err)
		}

		return resp, nil, nil
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebsocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	pairs := make([]string, 0)

	for _, cp := range cps {
		pair, ok := h.config.Cache[cp]
		if !ok {
			h.logger.Debug("pair ID not found for currency pair", zap.String("currency_pair", cp.ToString()))
			continue
		}

		pairs = append(pairs, string(TickerChannel)+"."+pair)
	}

	h.logger.Debug("subscribing to pairs", zap.Any("pairs", pairs))
	return NewSubscriptionRequestMessage(pairs)
}

// Name returns the name of the provider.
func (h *WebsocketDataHandler) Name() string {
	return Name
}

// URL returns the URL of the provider.
func (h *WebsocketDataHandler) URL() string {
	if h.config.Production {
		return ProductionURL
	}

	return TestnetURL
}
