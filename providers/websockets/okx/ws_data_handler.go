package okx

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the OKX websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the OKX API.
	market mmtypes.MarketConfig
	// ws is the config for the OKX websocket.
	ws config.WebSocketConfig
}

// NewWebSocketDataHandler returns a new OKX PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	marketCfg mmtypes.MarketConfig,
	wsCfg config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
	if err := marketCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market config for %s: %w", Name, err)
	}

	if marketCfg.Name != Name {
		return nil, fmt.Errorf("expected market config name %s, got %s", Name, marketCfg.Name)
	}

	if wsCfg.Name != Name {
		return nil, fmt.Errorf("expected websocket config name %s, got %s", Name, wsCfg.Name)
	}

	if !wsCfg.Enabled {
		return nil, fmt.Errorf("websocket config for %s is not enabled", Name)

	}

	if err := wsCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid websocket config for %s: %w", Name, err)
	}

	return &WebSocketHandler{
		logger: logger,
		market: marketCfg,
		ws:     wsCfg,
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The OKX
// provider sends two types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker response message. This is sent when a ticker update is received from the
//     OKX websocket API.
//
// Heartbeat messages are NOT sent by the OKX websocket. The connection is only closed
// iff no data is received within a 30-second interval or if all subscriptions
// fail. In the case where no data is received within a 30-second interval, the OKX
// will be restarted after the configured restart interval.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp        types.PriceResponse
		baseMessage BaseMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		return resp, nil, err
	}

	eventType := EventType(baseMessage.Event)
	switch {
	case eventType == EventSubscribe || eventType == EventError:
		h.logger.Debug("received subscribe response message")

		var subscribeMessage SubscribeResponseMessage
		if err := json.Unmarshal(message, &subscribeMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
		}

		updateMessage, err := h.parseSubscribeResponseMessage(subscribeMessage)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %w", err)
		}

		return resp, updateMessage, nil
	case eventType == EventTickers:
		h.logger.Debug("received ticker response message")

		var tickerMessage IndexTickersResponseMessage
		if err := json.Unmarshal(message, &tickerMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker response message: %w", err)
		}

		resp, err := h.parseTickerResponseMessage(tickerMessage)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse ticker response message: %w", err)
		}

		return resp, nil, nil
	default:
		return resp, nil, fmt.Errorf("unknown message type %s", baseMessage.Event)
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebSocketHandler) CreateMessages(
	tickers []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]SubscriptionTopic, 0)
	for _, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker.String()]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", ticker.String())
		}

		instruments = append(instruments, SubscriptionTopic{
			Channel:      string(IndexTickersChannel),
			InstrumentID: market.OffChainTicker,
		})
	}

	return NewSubscribeToTickersRequestMessage(instruments)
}

// HeartBeatMessages is not used for okx.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
