package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// The following URLs are used for the Coinbase Websocket feed. These can be found
	// in the Coinbase documentation at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// Name is the name of the Coinbase provider.
	Name = "coinbase"

	// URL is the production Coinbase Websocket URL.
	URL = "wss://ws-feed.exchange.coinbase.com"

	// URL_SANDBOX is the sandbox Coinbase Websocket URL.
	URL_SANDBOX = "wss://ws-feed-public.sandbox.exchange.coinbase.com" //nolint
)

const (
	// The following websocket configuration values were taken from the Coinbase documentation
	// at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// DefaultEnabledCompression is the default enabled compression for the Coinbase Websocket.
	// It is recommended to set this as true (please see the Coinbase documentation for more).
	DefaultEnabledCompression = false

	// DefaultWriteTimeout is the default write timeout for the Coinbase Websocket.
	// As recommended by Coinbase, this is set to 5 seconds.
	DefaultWriteTimeout = 5 * time.Second
)

var (
	// DefaultWebSocketConfig is the default configuration for the Coinbase Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:                       true,
		Name:                          Name,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URL,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             DefaultEnabledCompression,
		WriteTimeout:                  DefaultWriteTimeout,
		ReadTimeout:                   config.DefaultReadTimeout,
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"ATOM/USD": {
				Ticker:       "ATOM-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"ATOM/USDC": {
				Ticker:       "ATOM-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDC"),
			},
			"ATOM/USDT": {
				Ticker:       "ATOM-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USD": {
				Ticker:       "AVAX-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDC": {
				Ticker:       "AVAX-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDC"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USD": {
				Ticker:       "BTC-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDC": {
				Ticker:       "BTC-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTC-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIA-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"CELESTIA/USDC": {
				Ticker:       "TIA-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDC"),
			},
			"CELESTIA/USDT": {
				Ticker:       "TIA-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"DYDX/USD": {
				Ticker:       "DYDX-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"DYDX/USDC": {
				Ticker:       "DYDX-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDC"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDX-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH-BTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETH-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"OSMOSIS/USD": {
				Ticker:       "OSMO-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
			},
			"OSMOSIS/USDC": {
				Ticker:       "OSMO-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDC"),
			},
			"OSMOSIS/USDT": {
				Ticker:       "OSMO-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDT"),
			},
			"SOLANA/USD": {
				Ticker:       "SOL-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOL-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDC/USD": {
				Ticker:       "USDC-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
			},
			"USDC/USDT": {
				Ticker:       "USDC-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
			},
			"USDT/USD": {
				Ticker:       "USDT-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
			},
		},
	}
)
