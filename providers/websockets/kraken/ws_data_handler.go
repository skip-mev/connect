package kraken

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
	// Name is the name of the Kraken provider.
	Name = "kraken"
)

var _ handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] = (*WebSocketDataHandler)(nil)

// WebSocketDataHandler implements the WebSocketDataHandler interface. This is used to
// handle messages received from the Kraken websocket API.
type WebSocketDataHandler struct {
	logger *zap.Logger

	// config is the config for the Kraken web socket API.
	cfg config.ProviderConfig
}

// NewWebSocketDataHandler returns a new WebSocketDataHandler implementation for Kraken.
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

	return &WebSocketDataHandler{
		cfg:    cfg,
		logger: logger.With(zap.String("web_socket_data_handler", Name)),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. There are two
// types of messages that are handled by this function:
//  1. Price update messages. This is used to update the price of the given currency pair. This
//     is formated as a JSON array.
//  2. General response messages. This is used to check if the subscription request was successful,
//     heartbeats, and system status updates.
func (h *WebSocketDataHandler) HandleMessage(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], []handlers.WebsocketEncodedMessage, error) {
	var (
		resp        providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
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
			"failed to decode ticker response message; an unexpected message type was likely received: %s", err,
		)
	}

	// Parse the ticker response message and extract the price.
	resp, err = h.parseTickerMessage(tickerResponse)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to parse ticker message: %s", err)
	}

	return resp, nil, nil
}

// CreateMessage is used to create a message to send to the data provider. This is used to
// subscribe to the given currency pairs. This is called when the connection to the data
// provider is first established.
func (h *WebSocketDataHandler) CreateMessages(
	cps []oracletypes.CurrencyPair,
) ([]handlers.WebsocketEncodedMessage, error) {
	instruments := make([]string, 0)

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			h.logger.Debug("no instrument found for currency pair", zap.String("currency_pair", cp.String()))
			continue
		}

		instruments = append(instruments, market.Ticker)
	}

	return NewSubscribeRequestMessage(instruments)
}

// HeartBeatMessages is not used for Kraken.
func (h *WebSocketDataHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}
