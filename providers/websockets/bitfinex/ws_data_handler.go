package bitfinex

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the BitFinex websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the BitFinex API.
	market types.ProviderMarketMap
	// ws is the config for the BitFinex websocket.
	ws config.WebSocketConfig
	// channelMap maps a given channel_id to the currency pair its subscription represents.
	channelMap map[int]mmtypes.Ticker
}

// NewWebSocketDataHandler returns a new BitFinex PriceWebSocketDataHandler.
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
		logger:     logger,
		market:     market,
		ws:         ws,
		channelMap: make(map[int]mmtypes.Ticker),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The BitFinex
// provider sends four types of messages:
//
//  1. Subscribed response message. The subscribe response message is used to determine if
//     the subscription was successful.  If successful, the channel ID is saved.
//  2. Error response messages.  These messages provide info about errors from requests
//     sent to the BitFinex websocket API.
//  3. Ticker stream message. This is sent when a ticker update is received from the
//     BitFinex websocket API.
//  4. Heartbeat stream messages.  These are sent every 15 seconds by the BitFinex API.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp              types.PriceResponse
		baseMessage       BaseMessage
		subscribedMessage SubscribedMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		// if it is not a base json struct, we are receiving a stream.
		resp, err := h.handleStream(message)
		return resp, nil, err
	}

	switch Event(baseMessage.Event) {
	case EventSubscribed:
		h.logger.Debug("received subscribed response message")

		if err := json.Unmarshal(message, &subscribedMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
		}
		if err := h.parseSubscribedMessage(subscribedMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %w", err)
		}

		h.logger.Debug(
			"successfully subscribed",
			zap.String("pair", subscribedMessage.Pair),
			zap.Int("channel_id", subscribedMessage.ChannelID),
		)

		return resp, nil, nil

	case EventError:
		h.logger.Debug("received error message")

		var errorMessage ErrorMessage
		if err := json.Unmarshal(message, &errorMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal error message: %w", err)
		}

		updateMessage, err := h.parseErrorMessage(errorMessage)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to parse error message: %w", err)
		}

		return resp, updateMessage, nil
	default:
		return resp, nil, fmt.Errorf("unknown message: %x", message)
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the tickers that are specified in the config are subscribed to. The only channel that is
// subscribed to is the index tickers channel - which supports spot markets.
func (h *WebSocketHandler) CreateMessages(
	tickers []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	if len(tickers) == 0 {
		return nil, nil
	}

	msgs := make([]handlers.WebsocketEncodedMessage, len(tickers))
	for i, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker]
		if !ok {
			return nil, fmt.Errorf("ticker %s not in config", ticker.String())
		}

		msg, err := NewSubscribeMessage(market.OffChainTicker)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs[i] = msg

	}

	return msgs, nil
}

// HeartBeatMessages is not used for BitFinex.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger:     h.logger,
		market:     h.market,
		ws:         h.ws,
		channelMap: make(map[int]mmtypes.Ticker),
	}
}
