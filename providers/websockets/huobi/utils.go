package huobi

import (
	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
			"ATOM/USDT": {
				Ticker:       "atomusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USDT": {
				Ticker:       "avaxusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USDC": {
				Ticker:       "btcusdc",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"BITCOIN/USDT": {
				Ticker:       "btcusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"CELESTIA/USDT": {
				Ticker:       "tiausdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"DYDX/USDT": {
				Ticker:       "dydxusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ethbtc",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ethusdc",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ethusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"SOLANA/USDT": {
				Ticker:       "solusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDC/USDT": {
				Ticker:       "usdcusdt",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
			},
		},
	}
)
