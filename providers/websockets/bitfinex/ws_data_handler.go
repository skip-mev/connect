package bitfinex

import (
	"encoding/json"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

var _ handlers.WebSocketDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the BitFinex websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// channelMap maps a given channel_id to the currency pair its subscription represents.
	channelMap map[int]config.CurrencyPairMarketConfig

	// config is the config for the BitFinex websocket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for BitFinex
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
		cfg:        cfg,
		channelMap: make(map[int]config.CurrencyPairMarketConfig),
		logger:     logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The BitFinex
// provider sends four types of messages:
//
//  1. Subscribed response message. The subscribe response message is used to determine if
//     the subscription was successful.  If successful, the channel ID is saved
//  2. Error response messages.  These messages provide info about errors from requests
//     sent to the BitFinex websocket API
//  3. Ticker stream message. This is sent when a ticker update is received from the
//     BitFinex websocket API.
//  4. Heartbeat stream messages.  These are sent every 15 seconds by the BitFinex API
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp              providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
		baseMessage       BaseMessage
		subscribedMessage SubscribedMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		// if it is not a base json struct, we are receiving a stream
		resp, err := h.handleStream(message)
		return resp, nil, err
	}

	switch Event(baseMessage.Event) {
	case EventSubscribed:
		h.logger.Debug("received subscribed response message")

		if err := json.Unmarshal(message, &subscribedMessage); err != nil {
			h.logger.Error("failed to unmarshal subscription response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
		}

		err := h.parseSubscribedMessage(subscribedMessage)
		if err != nil {
			h.logger.Error("failed to parse subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %w", err)
		}

		h.logger.Debug("successfully subscribed", zap.String("pair", subscribedMessage.Pair), zap.Int("channel_id", subscribedMessage.ChannelID))

		return resp, nil, nil

	case EventError:
		h.logger.Debug("received error message")

		var errorMessage ErrorMessage
		if err := json.Unmarshal(message, &errorMessage); err != nil {
			h.logger.Error("failed to unmarshal error message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal error message: %w", err)
		}

		updateMessage, err := h.parseErrorMessage(errorMessage)
		if err != nil {
			h.logger.Error("failed to parse error message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse error message: %w", err)
		}

		return resp, updateMessage, nil
	default:
		h.logger.Error("unknown message", zap.Binary("message", message))
		return resp, nil, fmt.Errorf("unknown message: %x", message)
	}
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
			h.logger.Debug("instrument ID not found for currency pair", zap.String("currency_pair", cp.String()))
			return nil, fmt.Errorf("currency pair %s not in config", cp.String())
		}

		msg, err := NewSubscribeMessage(market.Ticker)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs[i] = msg

	}

	h.logger.Debug("subscribing to currency pairs", zap.Any("pairs", cps))
	return msgs, nil
}

// HeartBeatMessages is not used for BitFinex.
func (h *WebsocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// UpdateChannelMap updates the internal map for the given channelID and ticker.
func (h *WebsocketDataHandler) UpdateChannelMap(channelID int, ticker string) error {
	market, ok := h.cfg.Market.TickerToMarketConfigs[ticker]
	if !ok {
		return fmt.Errorf("unable to find market for currency pair: %s", ticker)
	}

	h.channelMap[channelID] = market
	return nil
}
