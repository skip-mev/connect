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
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ handlers.WebSocketDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Huobi websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the Huobi websocket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Huobi
// from a given provider configuration.
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

	return &WebsocketDataHandler{
		cfg:    cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
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
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp                 providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
		pingMessage          PingMessage
		subscriptionResponse SubscriptionResponse
		tickerStream         TickerStream
	)

	reader, err := gzip.NewReader(bytes.NewReader(message))
	if err != nil {
		h.logger.Error("error creating gzip reader", zap.Error(err))
		return resp, nil, err
	}
	defer func(reader *gzip.Reader) {
		closeErr := reader.Close()
		err = fmt.Errorf("error closing reader: %w. other errors: %w", closeErr, err)
	}(reader)

	var uncompressed bytes.Buffer
	_, err = io.Copy(&uncompressed, reader)
	if err != nil {
		h.logger.Error("error reading gzip", zap.Error(err))
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

	h.logger.Error("received unknown message", zap.String("message", uncompressed.String()))
	return resp, nil, fmt.Errorf("unknown message %s", uncompressed.String())
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebsocketDataHandler) CreateMessages(
	cps []slinkytypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	if len(cps) == 0 {
		return nil, nil
	}

	msgs := make([]handlers.WebsocketEncodedMessage, len(cps))

	for i, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("ID not found for currency pair", zap.String("currency_pair", cp.String()))
			return nil, fmt.Errorf("currency pair %s not in config", cp.String())
		}

		msg, err := NewSubscriptionRequest(market.Ticker)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs[i] = msg

	}

	h.logger.Debug("subscribing to currency pairs", zap.Any("pairs", cps))
	return msgs, nil
}

// HeartBeatMessages is not used for Huobi.
func (h *WebsocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
