package bybit

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the ByBit websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the ByBit API.
	market mmtypes.MarketConfig
	// ws is the config for the ByBit websocket.
	ws config.WebSocketConfig
}

// NewWebSocketDataHandler returns a new ByBit PriceWebSocketDataHandler.
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

// HandleMessage is used to handle a message received from the data provider. The ByBit
// provider sends three types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker update message. This is sent when a ticker update is received from the
//     ByBit websocket API.
//  3. Heartbeat update messages.  This should be sent every 20 seconds to ensure the
//     connection remains open.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp         types.PriceResponse
		baseResponse BaseResponse
		update       TickerUpdateMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseResponse); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
	}

	opType := Operation(baseResponse.Op)
	switch {
	case opType == OperationSubscribe:
		h.logger.Debug("received subscribe response message")

		var subscribeMessage SubscriptionResponse
		if err := json.Unmarshal(message, &subscribeMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
		}

		updateMessage, err := h.parseSubscriptionResponse(subscribeMessage)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %w", err)
		}

		return resp, updateMessage, nil
	case opType == OperationPing:
		h.logger.Debug("received pong response message")

		return resp, nil, nil
	default:
		// If the message is not a base message, then it must be a stream response.
		if err := json.Unmarshal(message, &update); err != nil {
			h.logger.Debug("unable to recognize message", zap.Error(err), zap.Binary("message", message))
			return resp, nil, err
		}

		// Parse the price information.
		resp, err := h.parseTickerUpdate(update)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse ticker update message: %w", err)
		}

		return resp, nil, nil
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the tickers that are specified in the config are subscribed to. The only channel that is
// subscribed to is the index tickers channel - which supports spot markets.
func (h *WebSocketHandler) CreateMessages(
	cps []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	pairs := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.market.TickerConfigs[cp.String()]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", cp.String())
		}

		pairs = append(pairs, string(TickerChannel)+"."+market.OffChainTicker)
	}

	return NewSubscriptionRequestMessage(pairs)
}

// HeartBeatMessages is used to construct heartbeat messages to be sent to the data provider. Note that
// the handler must maintain the necessary state information to construct the heartbeat messages. This
// can be done on the fly as messages as handled by the handler.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	msg, err := json.Marshal(HeartbeatPing{BaseRequest{
		ReqID: time.Now().String(),
		Op:    string(OperationPing),
	}})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal heartbeat ping: %w", err)
	}

	return []handlers.WebsocketEncodedMessage{msg}, nil
}
