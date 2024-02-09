package cryptodotcom

import (
	"encoding/json"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ handlers.WebSocketDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*WebSocketDataHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Crypto.com websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the Crypto.com websocket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Crypto.com.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	cfg config.ProviderConfig,
) (handlers.WebSocketDataHandler[slinkytypes.CurrencyPair, *big.Int], error) {
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

// HandleMessage is used to handle a message received from the data provider. The Crypto.com
// websocket API sends a heartbeat message every 30 seconds. If a heartbeat message is received,
// a heartbeat response message must be sent back to the Crypto.com websocket API, otherwise
// the connection will be closed. If a subscribe message is received, the message must be parsed
// and a response must be returned. No update message is required for subscribe messages.
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		msg  InstrumentResponseMessage
		resp providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message: %w", err)
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

		return resp, []handlers.WebsocketEncodedMessage{heartbeatResp}, nil
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

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given currency pairs. This is called when the connection to the data
// provider is first established.
func (h *WebSocketDataHandler) CreateMessages(
	cps []slinkytypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	// Iterate through each currency pair and get the instrument name. The instrument name
	// corresponds to the perpetual contract name on the Crypto.com websocket API. This will
	// only subscribe to price feeds that are configured in the config file.
	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("no market configuration for currency pair", zap.String("currency_pair", cp.String()))
			continue
		}

		instruments = append(instruments, fmt.Sprintf(TickerChannel, market.Ticker))
	}

	h.logger.Debug("subscribing to instruments", zap.Strings("instruments", instruments))
	return NewInstrumentMessage(instruments)
}

// HeartBeatMessages is not used for Crypto.com.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
