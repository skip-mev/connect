package cryptodotcom

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
	// Name is the name of the Crypto.com provider.
	Name = "crypto_dot_com"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebSocketDataHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Crypto.com websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the Crypto.com web socket API.
	config Config
}

// NewWebSocketDataHandlerFromConfig returns a new WebSocketDataHandler implementation for Crypto.com.
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

	return &WebSocketDataHandler{
		config: cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Crypto.com.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	cfg Config,
) (handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}

	return &WebSocketDataHandler{
		config: cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Crypto.com
// web socket API sends a heartbeat message every 30 seconds. If a heartbeat message is received,
// a heartbeat response message must be sent back to the Crypto.com web socket API, otherwise
// the connection will be closed. If a subscribe message is received, the message must be parsed
// and a response must be returned. No update message is required for subscribe messages.
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []byte, error) {
	var (
		msg  InstrumentResponseMessage
		resp providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message: %s", err)
	}

	// The status code of the message must be 0 for success.
	if StatusCode(msg.Code) != SuccessStatusCode {
		return resp, nil, fmt.Errorf("got unexpected error code %d: %s", msg.Code, string(message))
	}

	// Case on the two supported methods
	switch Method(msg.Method) {
	case HeartBeatRequestMethod:
		h.logger.Debug("received heartbeat")

		// If a heartbeat is received, send a heartbeat response back. This will not include
		// any instrument data.
		heartbeatResp, err := NewHeartBeatResponseMessage(msg.ID)
		if err != nil {
			return resp, nil, err
		}

		return resp, heartbeatResp, nil
	case InstrumentMethod:
		h.logger.Debug("received instrument message")

		// If a subscribe message is received, parse the message and send a response.
		subscribeResp, err := h.parseInstrumentMessage(msg)
		if err != nil {
			return resp, nil, err
		}

		return subscribeResp, nil, nil
	default:
		return resp, nil, fmt.Errorf("unknown method %s", msg.Method)
	}
}

// CreateMessage is used to create a message to send to the data provider. This is used to
// subscribe to the given currency pairs. This is called when the connection to the data
// provider is first established.
func (h *WebSocketDataHandler) CreateMessage(
	cps []oracletypes.CurrencyPair,
) ([]byte, error) {
	instruments := make([]string, 0)

	// Iterate through each currency pair and get the instrument name. The instrument name
	// corresponds to the perpetual contract name on the Crypto.com web socket API. This will
	// only subscribe to price feeds that are configured in the config file.
	for _, cp := range cps {
		instrument, ok := h.config.Cache[cp]
		if !ok {
			h.logger.Debug("no instrument for currency pair", zap.String("currency_pair", cp.ToString()))
			continue
		}

		instruments = append(instruments, fmt.Sprintf(TickerChannel, instrument))
	}

	h.logger.Debug("subscribing to instruments", zap.Strings("instruments", instruments))
	return NewInstrumentMessage(instruments)
}

// Name returns the name of the data provider.
func (h *WebSocketDataHandler) Name() string {
	return Name
}

// URL is used to get the URL of the data provider.
func (h *WebSocketDataHandler) URL() string {
	if h.config.Production {
		return ProductionURL
	}

	return SandboxURL
}
