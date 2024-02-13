package cryptodotcom

import (
	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
			"ATOM/USD": {
				Ticker:       "ATOMUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
			},
			"ATOM/USDT": {
				Ticker:       "ATOM_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USD": {
				Ticker:       "AVAXUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USD": {
				Ticker:       "BTCUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTC_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIAUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"CELESTIA/USDT": {
				Ticker:       "TIA_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"DYDX/USD": {
				Ticker:       "DYDXUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDX_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH_BTC",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETHUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"OSMOSIS/USD": {
				Ticker:       "OSMO_USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "SOLUSD-PERP",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDT/USD": {
				Ticker:       "USDT_USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
			},
		},
	}
)
