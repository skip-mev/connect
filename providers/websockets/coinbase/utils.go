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
			"BITCOIN/USD/8": {
				Ticker:       "BTC-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/USD/8": {
				Ticker:       "ETH-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
			},
			"ATOM/USD/8": {
				Ticker:       "ATOM-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD", oracletypes.DefaultDecimals),
			},
			"SOLANA/USD/8": {
				Ticker:       "SOL-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD", oracletypes.DefaultDecimals),
			},
			"CELESTIA/USD/8": {
				Ticker:       "TIA-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD", oracletypes.DefaultDecimals),
			},
			"AVAX/USD/8": {
				Ticker:       "AVAX-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD", oracletypes.DefaultDecimals),
			},
			"DYDX/USD/8": {
				Ticker:       "DYDX-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/BITCOIN/8": {
				Ticker:       "ETH-BTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN", oracletypes.DefaultDecimals),
			},
			"OSMOSIS/USD/8": {
				Ticker:       "OSMO-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD", oracletypes.DefaultDecimals),
			},
		},
	}
)
