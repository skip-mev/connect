package huobi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Huobi websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the Huobi API.
	market types.ProviderMarketMap
	// ws is the config for the Huobi websocket.
	ws config.WebSocketConfig
}

// NewWebSocketDataHandler returns a new Huobi PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	market types.ProviderMarketMap,
	ws config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
	if err := market.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market config for %s: %w", Name, err)
	}

	if market.Name != Name {
		return nil, fmt.Errorf("expected market config name %s, got %s", Name, market.Name)
	}

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
		market: market,
		ws:     ws,
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
		return resp, updateMessage, err
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
	tickers []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	if len(tickers) == 0 {
		return nil, fmt.Errorf("no tickers to subscribe to")
	}

	msgs := make([]handlers.WebsocketEncodedMessage, len(tickers))
	for i, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", ticker.String())
		}

		msg, err := NewSubscriptionRequest(market.OffChainTicker)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs[i] = msg

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
		market: h.market,
		ws:     h.ws,
	}
}
