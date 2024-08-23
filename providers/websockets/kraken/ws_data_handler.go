package kraken

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Kraken websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Kraken websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Kraken PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	ws config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
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
		logger: logger,
		ws:     ws,
		cache:  types.NewProviderTickers(),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. There are two
// types of messages that are handled by this function:
//  1. Price update messages. This is used to update the price of the given ticker. This
//     is formatted as a JSON array.
//  2. General response messages. This is used to check if the subscription request was successful,
//     heartbeats, and system status updates.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp        types.PriceResponse
		baseMessage BaseMessage
	)

	// If the message is able to be unmarshalled into a base message, then it is a general
	// response message. Otherwise, we check if it is a price update message.
	if err := json.Unmarshal(message, &baseMessage); err == nil {
		updateMessage, err := h.parseBaseMessage(message, Event(baseMessage.Event))
		return resp, updateMessage, err
	}

	// If the response cannot be decoded into a ticker response message, then it is likely
	// an unknown message type.
	tickerResponse, err := DecodeTickerResponseMessage(message)
	if err != nil {
		return resp, nil, fmt.Errorf(
			"failed to decode ticker response message; an unexpected message type was likely received: %w", err,
		)
	}

	// Parse the ticker response message and extract the price.
	resp, err = h.parseTickerMessage(tickerResponse)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to parse ticker message: %w", err)
	}

	return resp, nil, nil
}

// CreateMessages is used to create a message to send to the data provider. This is used to
// subscribe to the given tickers. This is called when the connection to the data provider
// is first established.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, ticker := range tickers {
		instruments = append(instruments, ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	return h.NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is not used for Kraken.
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger: h.logger,
		ws:     h.ws,
		cache:  types.NewProviderTickers(),
	}
}
