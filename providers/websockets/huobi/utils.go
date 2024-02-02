package huobi

import (
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Huobi provides the following URLs for its Websocket API. More info can be found in the documentation
	// here: https://huobiapi.github.io/docs/spot/v1/en/#websocket-market-data.

	// Name is the name of the Huobi provider.
	Name = "huobi"

	// URL is the public Huobi Websocket URL.
	URL = "wss://api.huobi.pro/ws"

	// URLAws is the public Huobi Websocket URL hosted on AWS.
	URLAws = "wss://api-aws.huobi.pro/ws"
)

var (
	// DefaultWebSocketConfig is the default configuration for the Huobi Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
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

	// DefaultMarketConfig is the default market configuration for the Huobi Websocket.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD/8": {
				Ticker:       "btcusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/USD/8": {
				Ticker:       "ethusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
			},
			"ATOM/USD/8": {
				Ticker:       "atomusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD", oracletypes.DefaultDecimals),
			},
			"SOLANA/USD/8": {
				Ticker:       "solusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD", oracletypes.DefaultDecimals),
			},
			"CELESTIA/USD/8": {
				Ticker:       "tiausdt",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD", oracletypes.DefaultDecimals),
			},
			"AVAX/USD/8": {
				Ticker:       "avaxusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD", oracletypes.DefaultDecimals),
			},
			"DYDX/USD/8": {
				Ticker:       "dydxusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/BITCOIN/8": {
				Ticker:       "ethbtc",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN", oracletypes.DefaultDecimals),
			},
		},
	}
)
