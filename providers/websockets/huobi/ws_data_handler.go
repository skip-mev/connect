package huobi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/klauspost/compress/gzip"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ handlers.WebSocketDataHandler[mmtypes.Ticker, *big.Int] = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Huobi websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the Gate.io API.
	market mmtypes.MarketConfig
	// ws is the config for the Gate.io websocket.
	ws config.WebSocketConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Huobi
// from a given provider configuration.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	marketCfg mmtypes.MarketConfig,
	wsCfg config.WebSocketConfig,
) (handlers.WebSocketDataHandler[mmtypes.Ticker, *big.Int], error) {
	if err := marketCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %w", err)
	}

	if marketCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, marketCfg.Name)
	}

	if wsCfg.Name != Name {
		return nil, fmt.Errorf("expected websocket config name %s, got %s", Name, wsCfg.Name)
	}

	if !wsCfg.Enabled {
		return nil, fmt.Errorf("websocket config for %s is not enabled", Name)
	}

	if err := wsCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid websocket config %w", err)
	}

	return &WebSocketHandler{
		logger: logger,
		market: marketCfg,
		ws:     wsCfg,
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
) (providertypes.GetResponse[mmtypes.Ticker, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp                 providertypes.GetResponse[mmtypes.Ticker, *big.Int]
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
	for i, cp := range tickers {
		market, ok := h.market.TickerConfigs[cp.String()]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", market.Ticker.String())
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
