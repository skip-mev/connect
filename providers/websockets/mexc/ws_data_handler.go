package mexc

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the MEXC websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// market is the config for the MEXC API.
	market mmtypes.MarketConfig
	// ws is the config for the MEXC websocket.
	ws config.WebSocketConfig
}

// NewWebSocketDataHandler returns a new MEXC PriceWebSocketDataHandler.
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

// HandleMessage is used to handle a message received from the data provider. This is called
// when a message is received from the data provider. There are three types of messages that
// can be received from the data provider:
//
// 1. A message that confirms that the client has successfully subscribed to a channel.
// 2. A message that confirms that the client has successfully pinged the server.
// 3. A message that contains the latest price for a ticker.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp      types.PriceResponse
		msg       BaseMessage
		tickerMsg TickerResponseMessage
	)

	if err := json.Unmarshal(message, &msg); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal message %w", err)
	}

	// If the base message is empty, we assume it is a price message.
	if msg.IsEmpty() {
		if err := json.Unmarshal(message, &tickerMsg); err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal ticker message %w", err)
		}

		// Parse the ticker message.
		resp, err := h.parseTickerResponseMessage(tickerMsg)
		return resp, nil, err
	}

	// Otherwise, we assume it is a subscription or pong message.
	switch {
	case strings.HasPrefix(msg.Message, string(MiniTickerChannel)):
		h.logger.Debug("subscribed to ticker channel", zap.String("instruments", msg.Message))
		return resp, nil, nil
	case MethodType(msg.Message) == PongMethod:
		h.logger.Debug("received pong message")
		return resp, nil, nil
	default:
		return resp, nil, fmt.Errorf("invalid message type %s", msg.Message)
	}
}

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given ticker. This is called when the connection to the data provider is
// first established.
func (h *WebSocketHandler) CreateMessages(
	tickers []mmtypes.Ticker,
) ([]handlers.WebsocketEncodedMessage, error) {
	if len(tickers) > MaxSubscriptionsPerConnection {
		return nil, fmt.Errorf("cannot subscribe to more than %d tickers per connection", MaxSubscriptionsPerConnection)

	}

	instruments := make([]string, 0)

	for _, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker.String()]
		if !ok {
			return nil, fmt.Errorf("ticker not found in market configs %s", ticker.String())
		}

		mexcTicker := fmt.Sprintf("%s%s%s", string(MiniTickerChannel), strings.ToUpper(market.OffChainTicker), "@UTC+8")
		instruments = append(instruments, mexcTicker)
	}

	return NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is used by the MEXC handler to send heart beat messages to the data provider.
// This is used to keep the connection alive when no messages are being sent from the data provider.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return NewPingRequestMessage()
}
