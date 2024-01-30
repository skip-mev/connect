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
		Name:                Name,
		Enabled:             true,
		MaxBufferSize:       1000,
		ReconnectionTimeout: config.DefaultReconnectionTimeout,
		WSS:                 URL,
		ReadBufferSize:      config.DefaultReadBufferSize,
		WriteBufferSize:     config.DefaultWriteBufferSize,
		HandshakeTimeout:    config.DefaultHandshakeTimeout,
		EnableCompression:   config.DefaultEnableCompression,
		ReadTimeout:         config.DefaultReadTimeout,
		WriteTimeout:        config.DefaultWriteTimeout,
	}

	// DefaultMarketConfig is the default market configuration for the Huobi Websocket.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "btcusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ethusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ATOM/USD": {
				Ticker:       "atomusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "solusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"CELESTIA/USD": {
				Ticker:       "tiausdt",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"AVAX/USD": {
				Ticker:       "avaxusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"DYDX/USD": {
				Ticker:       "dydxusdt",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ethbtc",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
		},
	}
)
