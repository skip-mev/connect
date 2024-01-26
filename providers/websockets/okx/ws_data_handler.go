package okx

import (
	"encoding/json"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the OKX provider.
	Name = "okx"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the OKX websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the OKX web socket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for OKX
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
		cfg:    cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The OKX
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker response message. This is sent when a ticker update is received from the
//     OKX web socket API.
//
// Heartbeat messages are NOT sent by the OKX web socket. The connection is only closed
// iff no data is received within a 30 second interval or if all of the subscriptions
// fail. In the case where no data is received within a 30 second interval, the OKX
// will be restarted after the configured restart interval.
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
		baseMessage BaseMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		h.logger.Debug("unable to unmarshal message into base message", zap.Error(err))
		return resp, nil, err
	}

	eventType := EventType(baseMessage.Event)
	switch {
	case eventType == EventSubscribe || eventType == EventError:
		h.logger.Debug("received subscribe response message")

		var subscribeMessage SubscribeResponseMessage
		if err := json.Unmarshal(message, &subscribeMessage); err != nil {
			h.logger.Error("failed to unmarshal subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %s", err)
		}

		updateMessage, err := h.parseSubscribeResponseMessage(subscribeMessage)
		if err != nil {
			h.logger.Error("failed to parse subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %s", err)
		}

		return resp, updateMessage, nil
	case eventType == EventTickers:
		h.logger.Debug("received ticker response message")

		var tickerMessage IndexTickersResponseMessage
		if err := json.Unmarshal(message, &tickerMessage); err != nil {
			h.logger.Error("failed to unmarshal ticker response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal ticker response message: %s", err)
		}

		resp, err := h.parseTickerResponseMessage(tickerMessage)
		if err != nil {
			h.logger.Error("failed to parse ticker response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse ticker response message: %s", err)
		}

		return resp, nil, nil
	default:
		h.logger.Debug("received unknown message", zap.String("message", string(message)))
		return resp, nil, fmt.Errorf("unknown message type %s", baseMessage.Event)
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebsocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]SubscriptionTopic, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("instrument ID not found for currency pair", zap.String("currency_pair", cp.String()))
			continue
		}

		instruments = append(instruments, SubscriptionTopic{
			Channel:      string(IndexTickersChannel),
			InstrumentID: market.Ticker,
		})
	}

	h.logger.Debug("subscribing to instruments", zap.Any("instruments", instruments))
	return NewSubscribeToTickersRequestMessage(instruments)
}

// HeartBeatMessages is not used for okx.
func (h *WebsocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
