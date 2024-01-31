package cryptodotcom

import (
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// URL is the URL used to connect to the Crypto.com websocket API. This can be found here
	// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#websocket-root-endpoints
	// Note that Crypto.com offers a sandbox and production environment.

	// Name is the name of the Crypto.com provider.
	Name = "crypto_dot_com"

	// URL_PROD is the URL used to connect to the Crypto.com production websocket API.
	URL_PROD = "wss://stream.crypto.com/exchange/v1/market" //nolint

	// URL_SANDBOX is the URL used to connect to the Crypto.com sandbox websocket API. This will
	// return static prices.
	URL_SANDBOX = "wss://uat-stream.3ona.co/exchange/v1/market" //nolint
)

var (
	// DefaultWebSocketConfig is the default configuration for the Crypto.com Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URL_PROD,
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

	// DefaultMarketConfig is the default market configuration for Crypto.com.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTCUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETHUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
			},
			"ATOM/USD": {
				Ticker:       "ATOMUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD", oracletypes.DefaultDecimals),
			},
			"SOLANA/USD": {
				Ticker:       "SOLUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD", oracletypes.DefaultDecimals),
			},
			"CELESTIA/USD": {
				Ticker:       "TIAUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD", oracletypes.DefaultDecimals),
			},
			"AVAX/USD": {
				Ticker:       "AVAXUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD", oracletypes.DefaultDecimals),
			},
			"DYDX/USD": {
				Ticker:       "DYDXUSD-PERP",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH_BTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN", oracletypes.DefaultDecimals),
			},
			"OSMOSIS/USD": {
				Ticker:       "OSMO_USD",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD", oracletypes.DefaultDecimals),
			},
		},
	}
)
