package bitfinex

import (
	"encoding/json"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the BitFinex websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// channelMap maps a given channel_id to the currency pair its subscription represents.
	channelMap map[string]config.CurrencyPairMarketConfig

	// config is the config for the BitFinex web socket API.
	cfg config.ProviderConfig
}

// UpdateChannelMap updates the internal map from
func (h *WebsocketDataHandler) UpdateChannelMap(channelID, ticker string) error {
	market, ok := h.cfg.Market.TickerToMarketConfigs[ticker]
	if !ok {
		return fmt.Errorf("unable to find market for currency pair: %s", ticker)
	}

	h.channelMap[channelID] = market
	return nil
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for BitFinex
// from a given provider configuration.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	cfg config.ProviderConfig,
) (handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %s", err)
	}

	if !cfg.WebSocket.Enabled {
		return nil, fmt.Errorf("web socket is not enabled for provider %s", cfg.Name)
	}

	if cfg.Name != Name {
		return nil, fmt.Errorf("invalid provider name %s", cfg.Name)
	}

	return &WebsocketDataHandler{
		cfg:        cfg,
		channelMap: make(map[string]config.CurrencyPairMarketConfig),
		logger:     logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The OKX
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.  If successful, the channel ID is saved
//  2. Ticker stream message. This is sent when a ticker update is received from the
//     BitFinex web socket API.
//  3. Heartbeat messages.  These are sent every 15 seconds by the BitFinex API
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp              providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		baseMessage       BaseMessage
		baseStream        BaseStreamMessage
		subscribedMessage SubscribedMessage
		tickerStream      TickerStream
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		h.logger.Debug("unable to unmarshal message into base message", zap.Error(err))
		return resp, nil, err
	}

	switch Event(baseMessage.Event) {
	case EventSubscribed:
		h.logger.Debug("received subscribed response message")

		if err := json.Unmarshal(message, &subscribedMessage); err != nil {
			h.logger.Error("failed to unmarshal subscription response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %s", err)
		}

		err := h.parseSubscribedMessage(subscribedMessage)
		if err != nil {
			h.logger.Error("failed to parse subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %s", err)
		}

		return resp, nil, nil

	case EventError:
		h.logger.Debug("received error message")

		var errorMessage ErrorMessage
		if err := json.Unmarshal(message, &errorMessage); err != nil {
			h.logger.Error("failed to unmarshal error message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal error message: %s", err)
		}

		updateMessage, err := h.parseErrorMessage(errorMessage)
		if err != nil {
			h.logger.Error("failed to parse error message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse error message: %s", err)
		}

		return resp, updateMessage, nil

		// otherwise, this is a stream message
	case EventNil:
		if err := json.Unmarshal(message, &baseStream); err != nil {
			h.logger.Debug("unable to unmarshal message into stream message", zap.Error(err))
			return resp, nil, err
		}

		// stream messages are defined by their channel ID
		switch ChannelID(baseStream.ChannelID) {
		case ChannelIdHeartbeat:
			// handle incoming heartbeat
			h.logger.Debug("received heartbeat message")
			return resp, nil, nil

		default:
			if err := json.Unmarshal(message, &tickerStream); err != nil {
				h.logger.Error("failed to unmarshal ticker stream message", zap.Error(err))
				return resp, nil, fmt.Errorf("failed to unmarshal ticker stream message: %s", err)
			}

			resp, err := h.parseTickerStream(tickerStream)
			if err != nil {
				h.logger.Error("failed to parse ticker stream message", zap.Error(err))
				return resp, nil, fmt.Errorf("failed to parse ticker stream message: %s", err)
			}

			return resp, nil, nil
		}

	default:
		h.logger.Error("unknown message", zap.Any("message", message))
		return resp, nil, fmt.Errorf("unknown message: %v", message)
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebsocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	var msgs []handlers.WebsocketEncodedMessage

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("instrument ID not found for currency pair", zap.String("currency_pair", cp.String()))
			return nil, fmt.Errorf("currency pair %s not in config", cp.String())
		}

		msg, err := NewSubscribeMessage(market.Ticker)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs = append(msgs, msg)

	}

	h.logger.Debug("subscribing to currency pairs", zap.Any("pairs", cps))
	return msgs, nil
}

// HeartBeatMessages is not used for BitFinex.
func (h *WebsocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
