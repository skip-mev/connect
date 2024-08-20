package huobi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Huobi websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Huobi websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Huobi PriceWebSocketDataHandler.
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

// HandleMessage is used to handle a message received from the data provider. The Huobi
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker response message. This is sent when a ticker update is received from the
//     Huobi websocket API.
//  3. Heartbeat ping message.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp                 types.PriceResponse
		pingMessage          PingMessage
		subscriptionResponse SubscriptionResponse
		tickerStream         TickerStream
	)

	reader, err := gzip.NewReader(bytes.NewReader(message))
	if err != nil {
		return resp, nil, err
	}
	defer func(reader *gzip.Reader) {
		closeErr := reader.Close()
		err = fmt.Errorf("error closing reader: %w. other errors: %w", closeErr, err)
	}(reader)

	var uncompressed bytes.Buffer
	_, err = io.Copy(&uncompressed, reader)
	if err != nil {
		return resp, nil, err
	}

	// attempt to unmarshal to ping and check if field values are not nil
	if err := json.Unmarshal(uncompressed.Bytes(), &pingMessage); err == nil && pingMessage.Ping != 0 {
		h.logger.Debug("received ping message")
		updateMessage, err := NewPongMessage(pingMessage)

		// The receipt of a ping message means that the connection is still alive and that all market's corresponding
		// to the tickers subscribed to are still being tracked. Therefore, the response can include a message to let
		// the provider know that market prices are still valid.
		return h.cache.NoPriceChangeResponse(), updateMessage, err
	}

	// attempt to unmarshal to subscription response message and check if field values are not nil
	if err := json.Unmarshal(uncompressed.Bytes(), &subscriptionResponse); err == nil && subscriptionResponse.ID != "" {
		h.logger.Debug("received subscription response message")

		updateMsg, err := h.parseSubscriptionResponse(subscriptionResponse)
		if err != nil {
			return resp, updateMsg, fmt.Errorf("failed parse subscription message: %w", err)
		}

		return resp, updateMsg, nil
	}

	// attempt to unmarshal to ticker stream message and check if field values are not nil
	if err := json.Unmarshal(uncompressed.Bytes(), &tickerStream); err == nil && tickerStream.Channel != "" {
		h.logger.Debug("received ticker stream message")

		resp, err := h.parseTickerStream(tickerStream)
		if err != nil {
			return resp, nil, fmt.Errorf("failed parse ticker stream: %w", err)
		}

		return resp, nil, nil
	}

	return resp, nil, fmt.Errorf("unknown message %s", uncompressed.String())
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the tickers that are specified in the config are subscribed to. The only channel that is
// subscribed to is the index tickers channel - which supports spot markets.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	if len(tickers) == 0 {
		return nil, fmt.Errorf("no tickers to subscribe to")
	}

	msgs := make([]handlers.WebsocketEncodedMessage, len(tickers))
	for i, ticker := range tickers {
		msg, err := NewSubscriptionRequest(ticker.GetOffChainTicker())
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs[i] = msg
		h.cache.Add(ticker)
	}

	return msgs, nil
}

// HeartBeatMessages is not used for Huobi.
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
