package bybit

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the ByBit websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the ByBit websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new ByBit PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	ws config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
	if ws.Name != Name {
		return nil, fmt.Errorf("expected websocket config name %s, got %s", Name, ws.Name)
	}

	if !ws.Enabled {
		return nil, fmt.Errorf("websocket config for %s is not enabled", Name)
	}

	if err := ws.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid websocket config for %s: %w", Name, err)
	}

	return &WebSocketHandler{
		logger: logger,
		ws:     ws,
		cache:  types.NewProviderTickers(),
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

	switch Operation(baseResponse.Op) {
	case OperationSubscribe:
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
	case OperationPing:
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
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	pairs := make([]string, 0)

	for _, ticker := range tickers {
		pairs = append(pairs, string(TickerChannel)+"."+ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	return h.NewSubscriptionRequestMessage(pairs)
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

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger: h.logger,
		ws:     h.ws,
		cache:  types.NewProviderTickers(),
	}
}
