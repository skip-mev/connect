package gate

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Gate.io websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Gate.io websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Gate.io PriceWebSocketDataHandler.
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

// HandleMessage is used to handle a message received from the data provider. The Gate.io
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker stream message. This is sent when a ticker update is received from the
//     Gate.io websocket API.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp         types.PriceResponse
		baseMessage  BaseMessage
		subResponse  SubscribeResponse
		tickerStream TickerStream
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		return resp, nil, err
	}

	switch Event(baseMessage.Event) {
	case EventSubscribe:
		if err := json.Unmarshal(message, &subResponse); err != nil {
			return resp, nil, err
		}

		// handle subscription
		updateMsg, err := h.parseSubscribeResponse(subResponse)
		return resp, updateMsg, err

	case EventUpdate:
		if err := json.Unmarshal(message, &tickerStream); err != nil {
			return resp, nil, err
		}

		// update pair info
		resp, err := h.parseTickerStream(tickerStream)
		return resp, nil, err

	default:
		return resp, nil, fmt.Errorf("unknown message type %s", baseMessage.Event)
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the tickers that are specified in the config are subscribed to. The only channel that is
// subscribed to is the tickers channel - which supports spot markets.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		instruments = append(instruments, ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	return h.NewSubscribeRequest(instruments)
}

// HeartBeatMessages is not used for Gate.io.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger: h.logger,
		ws:     h.ws,
		cache:  types.NewProviderTickers(),
	}
}
