package coinbase

import (
	"encoding/json"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ handlers.WebSocketDataHandler[mmtypes.Ticker, *big.Int] = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Coinbase websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the Coinbase API.
	market mmtypes.MarketConfig
	// ws is the config for the Coinbase websocket.
	ws config.WebSocketConfig
	// Sequence is the current sequence number for the Coinbase websocket API per currency pair.
	sequence map[mmtypes.Ticker]int64
}

// NewWebSocketDataHandler returns a new Coinbase PriceWebSocketDataHandler.
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
		logger:   logger,
		market:   marketCfg,
		ws:       wsCfg,
		sequence: make(map[mmtypes.Ticker]int64),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Coinbase web
// socket expects the client to send a subscribe message within 5 seconds of the initial connection.
// Otherwise, the connection will be closed. There are two types of messages that can be received
// from the Coinbase websocket API:
//
//  1. SubscriptionsMessage: This is sent by the Coinbase websocket API after a subscribe message
//     is sent. This message contains the list of channels that were successfully subscribed to.
//  2. TickerMessage: This is sent by the Coinbase websocket API when a match happens. This message
//     contains the price of the ticker.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp types.PriceResponse
		msg  BaseMessage
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message %w", err)
	}

	switch MessageType(msg.Type) {
	case SubscriptionsMessage:
		h.logger.Debug("received subscriptions message")

		var subscriptionsMessage SubscribeResponseMessage
		if err := json.Unmarshal(message, &subscriptionsMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal subscriptions message %w", err)
		}

		// Log the channels that were successfully subscribed to.
		for _, channel := range subscriptionsMessage.Channels {
			for _, instrument := range channel.Instruments {
				h.logger.Debug("subscribed to ticker channel", zap.String("instrument", instrument), zap.String("channel", channel.Name))
			}
		}

		return resp, nil, nil
	case TickerMessage:
		h.logger.Debug("received ticker message")

		var tickerMessage TickerResponseMessage
		if err := json.Unmarshal(message, &tickerMessage); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %w", err)
		}

		resp, err := h.parseTickerResponseMessage(tickerMessage)
		return resp, nil, err
	default:
		return resp, nil, fmt.Errorf("invalid message type %s", msg.Type)
	}
}

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given tickers. This is called when the connection to the data provider is
// first established.
func (h *WebSocketHandler) CreateMessages(
	tickers []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker.String()]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", ticker.String())
		}

		instruments = append(instruments, market.OffChainTicker)
	}

	return NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is not used for Coinbase.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
