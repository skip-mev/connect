package binance

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
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
	// nextID is the next message ID to use for the Binance websocket API.
	nextID int64
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
		nextID:     rand.Int63() + 1,
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Binance websocket
// API is expected to handle the following types of messages:
//  1. SubscribeMessageResponse: This is a response to a subscription request. If the subscription
//     was successful, the response will contain a nil result. If the subscription failed, a
//     re-subscription message will be returned.
//  2. StreamMessageResponse: This is a response to a stream message. The stream message contains
//     the latest price of a ticker - either received when a trade is made or an automated price
//     update is received.
//
// Heartbeat messages are handled by default by the gorilla websocket library. The Binance websocket
// API does not require any additional heartbeat messages to be sent. The pong frames are sent
// automatically by the gorilla websocket library.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp      types.PriceResponse
		msg       SubscribeMessageResponse
		streamMsg StreamMessageResponse
	)

	// Unmarshal the message. If the message fails to be unmarshaled or is empty, this means
	// that we likely received a price update message.
	if err := json.Unmarshal(message, &msg); err == nil && !msg.IsEmpty() {
		instruments, ok := h.messageIDs[msg.ID]
		if !ok {
			return resp, nil, fmt.Errorf("failed to find instruments for message ID %d", msg.ID)
		}

		if msg.Result != nil {
			// If the result is not nil, this means that the subscription failed to be made. Return
			// an update message with the same subscription.
			h.logger.Debug("failed to make subscription; attempting to re-subscribe", zap.Any("instruments", msg))
			subscriptionMsgs, err := h.NewSubscribeRequestMessage(instruments)
			return resp, subscriptionMsgs, err
		}

		// If the result is nil, this means that the subscription was successful. Return an empty
		// response.
		h.logger.Debug("successfully subscribed to instruments", zap.Any("instruments", instruments))
		return resp, nil, nil
	}

	// Unmarshal the message as a stream message. If the message fails to be unmarshaled, this means
	// that we received an unknown message type.
	if err := json.Unmarshal(message, &streamMsg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message %w", err)
	}

	switch streamMsg.GetStreamType() {
	case TickerStream:
		// Ticker stream is sent every 1000ms and contains the latest price of a ticker.
		var tickerResp TickerMessageResponse
		if err := json.Unmarshal(message, &tickerResp); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %w", err)
		}

		h.logger.Debug("received ticker message", zap.String("ticker", tickerResp.Data.Ticker))
		resp, err := h.parsePriceUpdateMessage(tickerResp.Data.Ticker, tickerResp.Data.LastPrice)
		return resp, nil, err
	case AggregateTradeStream:
		// Aggregate trade stream is sent when a trade is executed on the Binance exchange.
		var aggTradeResp AggregatedTradeMessageResponse
		if err := json.Unmarshal(message, &aggTradeResp); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal aggregate trade message %w", err)
		}

		h.logger.Debug("received aggregate trade message", zap.String("ticker", aggTradeResp.Data.Ticker))
		resp, err := h.parsePriceUpdateMessage(aggTradeResp.Data.Ticker, aggTradeResp.Data.Price)
		return resp, nil, err
	default:
		return resp, nil, fmt.Errorf("unknown stream type %s", streamMsg.Stream)
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
		nextID:     rand.Int63() + 1,
	}
}
