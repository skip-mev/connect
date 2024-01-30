package kraken

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// URL is the websocket URL for Kraken. You can find the documentation here:
	// https://docs.kraken.com/websockets/. Kraken provides an authenticated and
	// unauthenticated websocket. The URLs defined below are all unauthenticated.

	// Name is the name of the Kraken provider.
	Name = "kraken"

	// URL is the production websocket URL for Kraken.
	URL = "wss://ws.kraken.com"

	// URL_BETA is the demo websocket URL for Kraken.
	URL_BETA = "wss://beta-ws.kraken.com" //nolint
)

var (
	// DefaultWebSocketConfig is the default configuration for the Kraken Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           10 * time.Second,
		WSS:                           URL,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  config.DefaultPingInterval,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for Kraken.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "XBT/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ATOM/USD": {
				Ticker:       "ATOM/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "SOL/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIA/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"AVAX/USD": {
				Ticker:       "AVAX/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"DYDX/USD": {
				Ticker:       "DYDX/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH/XBT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
		},
	}
)

// DecodeTickerResponseMessage decodes a ticker response message .
func DecodeTickerResponseMessage(message []byte) (TickerResponseMessage, error) {
	var rawResponse []json.RawMessage
	if err := json.Unmarshal(message, &rawResponse); err != nil {
		return TickerResponseMessage{}, err
	}

	if len(rawResponse) != ExpectedTickerResponseMessageLength {
		return TickerResponseMessage{}, fmt.Errorf(
			"invalid ticker response message; expected length %d, got %d", ExpectedTickerResponseMessageLength, len(rawResponse),
		)
	}

	var response TickerResponseMessage
	if err := json.Unmarshal(rawResponse[ChannelIDIndex], &response.ChannelID); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[TickerDataIndex], &response.TickerData); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[ChannelNameIndex], &response.ChannelName); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[PairIndex], &response.Pair); err != nil {
		return TickerResponseMessage{}, err
	}

	return response, nil
}
