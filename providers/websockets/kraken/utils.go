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
			"ATOM/USD": {
				Ticker:       "ATOM/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
			},
			"AVAX/USD": {
				Ticker:       "AVAX/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX/USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USD": {
				Ticker:       "XBT/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDC": {
				Ticker:       "XBT/USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"BITCOIN/USDT": {
				Ticker:       "XBT/USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIA/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"DYDX/USD": {
				Ticker:       "DYDX/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH/XBT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETH/USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH/USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"SOLANA/USD": {
				Ticker:       "SOL/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL/USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDC/USD": {
				Ticker:       "USDC/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
			},
			"USDC/USDT": {
				Ticker:       "USDC/USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
			},
			"USDT/USD": {
				Ticker:       "USDT/USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
