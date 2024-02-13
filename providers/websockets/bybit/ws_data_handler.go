package bybit

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ handlers.WebSocketDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*WebsocketDataHandler)(nil)

// WebsocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the ByBit websocket API.
type WebsocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the ByBit websocket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for ByBit.
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

// HandleMessage is used to handle a message received from the data provider. The ByBit
// provider sends three types of messages:
//
//  1. Subscribe response message. The subscribe response message is used to determine if
//     the subscription was successful.
//  2. Ticker update message. This is sent when a ticker update is received from the
//     ByBit websocket API.
//  3. Heartbeat update messages.  This should be sent every 20 seconds to ensure the
//     connection remains open.
func (h *WebsocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp         providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
		baseResponse BaseResponse
		update       TickerUpdateMessage
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseResponse); err != nil {
		h.logger.Error("failed to unmarshal subscribe response message", zap.Error(err))
		return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)

	}

	opType := Operation(baseResponse.Op)
	switch {
	case opType == OperationSubscribe:
		h.logger.Debug("received subscribe response message")

		var subscribeMessage SubscriptionResponse
		if err := json.Unmarshal(message, &subscribeMessage); err != nil {
			h.logger.Error("failed to unmarshal subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to unmarshal subscribe response message: %w", err)
		}

		updateMessage, err := h.parseSubscriptionResponse(subscribeMessage)
		if err != nil {
			h.logger.Error("failed to parse subscribe response message", zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse subscribe response message: %w", err)
		}

		return resp, updateMessage, nil
	case opType == OperationPing:
		h.logger.Debug("received pong response message")

		return resp, nil, nil
	default:
		// if the message is not a base message, then it must be a stream response
		if err := json.Unmarshal(message, &update); err != nil {
			h.logger.Debug("unable to recognize message", zap.Error(err), zap.Binary("message", message))
			return resp, nil, err
		}

		// parse
		resp, err := h.parseTickerUpdate(update)
		if err != nil {
			h.logger.Error("failed to parse ticker update message", zap.Any("base", baseResponse), zap.Any("message", update), zap.Error(err))
			return resp, nil, fmt.Errorf("failed to parse ticker update message: %w", err)
		}

		return resp, nil, nil
	}
}

// CreateMessages is used to create an initial subscription message to send to the data provider.
// Only the currency pairs that are specified in the config are subscribed to. The only channel
// that is subscribed to is the index tickers channel - which supports spot markets.
func (h *WebsocketDataHandler) CreateMessages(
	cps []slinkytypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	pairs := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("pair ID not found for currency pair", zap.String("currency_pair", cp.String()))
			continue
		}

		pairs = append(pairs, string(TickerChannel)+"."+market.Ticker)
	}

	h.logger.Debug("subscribing to pairs", zap.Any("pairs", pairs))
	return NewSubscriptionRequestMessage(pairs)
}

// HeartBeatMessages is used to construct heartbeat messages to be sent to the data provider. Note that
// the handler must maintain the necessary state information to construct the heartbeat messages. This
// can be done on the fly as messages as handled by the handler.
func (h *WebsocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	msg, err := json.Marshal(HeartbeatPing{BaseRequest{
		ReqID: time.Now().String(),
		Op:    string(OperationPing),
	}})
	if err != nil {
		h.logger.Debug("unable to marshal heartbeat ping")
		return nil, fmt.Errorf("unable to marshal heartbeat ping: %w", err)
	}

	return []handlers.WebsocketEncodedMessage{msg}, nil
}
