package binance

import (
	"encoding/json"
	"fmt"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"go.uber.org/zap"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to handle
// messages received from the Binance websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Binance websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
	// messageIDs is the current message ID for the Binance websocket API per currency pair(s).
	messageIDs map[int64][]string
}

// NewWebSocketDataHandler returns a new Binance PriceWebSocketDataHandler.
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
		logger:     logger,
		ws:         ws,
		cache:      types.NewProviderTickers(),
		messageIDs: make(map[int64][]string),
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
		resp     types.PriceResponse
		msg      SubscribeMessageResponse
		tradeMsg AggregatedTradeMessageResponse
	)

	// Unmarshal the message. If the message fails to be unmarshaled or is empty, this means
	// that we likely received a price update message.
	if err := json.Unmarshal(message, &msg); err != nil || msg.IsEmpty() {
		if err := json.Unmarshal(message, &tradeMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal message %w", err)
		}

		// Parse the message.
		resp, err := h.parseAggregateTradeMessage(tradeMsg)
		return resp, nil, err
	}

	instruments, ok := h.messageIDs[msg.ID]
	if !ok {
		return resp, nil, fmt.Errorf("failed to find instruments for message ID %d", msg.ID)
	}

	switch {
	case msg.Result != nil:
		// If the result is not nil, this means that the subscription failed to be made. Return
		// an update message with the same subscription.
		h.logger.Debug("failed to make subscription", zap.Any("instruments", msg))
		subscriptionMsg, err := h.NewSubscribeRequestMessage(instruments)
		return resp, subscriptionMsg, err
	default:
		// If the result is nil, this means that the subscription was successful. Return an empty
		// response.
		h.logger.Debug("successfully subscribed to instruments", zap.Any("instruments", instruments))
		return resp, nil, nil
	}
}

// CreateMessages is used to create a message to send to Binance. This is used to subscribe to
// the given tickers. This is called when the connection to the data provider is first established.
// Notably, the tickers have a unique identifier that is used to identify the messages going back
// and forth. This unique identifier is the same one sent in the initial subscription.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		instruments = append(instruments, ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	return h.NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is not used for Binance. Heartbeats are handled on an ad-hoc basis when
// messages are received from the Binance websocket API.
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger:     h.logger,
		ws:         h.ws,
		cache:      types.NewProviderTickers(),
		messageIDs: make(map[int64][]string),
	}
}
