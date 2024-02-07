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
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
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
			"BITCOIN/USDC": {
				Ticker:       "XBT/USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETH/USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH/USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"ATOM/USD": {
				Ticker:       "ATOM/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"ATOM/USDC": {
				Ticker:       "ATOM/USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDC"),
			},
			"SOLANA/USD": {
				Ticker:       "SOL/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL/USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIA/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"AVAX/USD": {
				Ticker:       "AVAX/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX/USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"DYDX/USD": {
				Ticker:       "DYDX/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH/XBT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"USDT/USD": {
				Ticker:       "USDT/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
			},
			"USDC/USD": {
				Ticker:       "USDC/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
			},
			"USDC/USDT": {
				Ticker:       "USDC/USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
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
