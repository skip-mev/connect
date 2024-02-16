package cryptodotcom

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Crypto.com websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the Crypto.com API.
	market mmtypes.MarketConfig
	// ws is the config for the Crypto.com websocket.
	ws config.WebSocketConfig
}

// NewWebSocketDataHandler returns a new Crypto.com PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	marketCfg mmtypes.MarketConfig,
	wsCfg config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
	if err := marketCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market config for %s: %w", Name, err)
	}

	if marketCfg.Name != Name {
		return nil, fmt.Errorf("expected market config name %s, got %s", Name, marketCfg.Name)
	}

	if wsCfg.Name != Name {
		return nil, fmt.Errorf("expected websocket config name %s, got %s", Name, wsCfg.Name)
	}

	if !wsCfg.Enabled {
		return nil, fmt.Errorf("websocket config for %s is not enabled", Name)

	}

	if err := wsCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid websocket config for %s: %w", Name, err)
	}

	return &WebSocketHandler{
		logger: logger,
		market: marketCfg,
		ws:     wsCfg,
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Crypto.com
// websocket API sends a heartbeat message every 30 seconds. If a heartbeat message is received,
// a heartbeat response message must be sent back to the Crypto.com websocket API, otherwise
// the connection will be closed. If a subscribe message is received, the message must be parsed
// and a response must be returned. No update message is required for subscribe messages.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		msg  InstrumentResponseMessage
		resp types.PriceResponse
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
// subscribe to the given tickers. This is called when the connection to the data provider
// is first established.
func (h *WebSocketHandler) CreateMessages(
	tickers []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker.String()]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", ticker.String())
		}

		instruments = append(instruments, fmt.Sprintf(TickerChannel, market.OffChainTicker))
	}

	return NewInstrumentMessage(instruments)
}

// HeartBeatMessages is not used for Crypto.com.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
