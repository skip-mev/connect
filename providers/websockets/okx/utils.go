package okx

import (
	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

const (
	// OKX provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://www.okx.com/docs-v5/en/?shell#overview-production-trading-services
	// The two production URLs are defined in ProductionURL and ProductionAWSURL. The
	// DemoURL is used for testing purposes.

	// Name is the name of the OKX provider.
	Name = "okx"

	// URL_PROD is the public OKX Websocket URL.
	URL_PROD = "wss://ws.okx.com:8443/ws/v5/public" //nolint

	// URL_PROD_AWS is the public OKX Websocket URL hosted on AWS.
	URL_PROD_AWS = "wss://wsaws.okx.com:8443/ws/v5/public" //nolint

	// URL_DEMO is the public OKX Websocket URL for test usage.
	URL_DEMO = "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999" //nolint
)

var (
	// DefaultWebSocketConfig is the default configuration for the OKX Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                Name,
		Enabled:             true,
		MaxBufferSize:       1000,
		ReconnectionTimeout: config.DefaultReconnectionTimeout,
		WSS:                 URL_PROD,
		ReadBufferSize:      config.DefaultReadBufferSize,
		WriteBufferSize:     config.DefaultWriteBufferSize,
		HandshakeTimeout:    config.DefaultHandshakeTimeout,
		EnableCompression:   config.DefaultEnableCompression,
		ReadTimeout:         config.DefaultReadTimeout,
		WriteTimeout:        config.DefaultWriteTimeout,
		MaxReadErrorCount:   config.DefaultMaxReadErrorCount,
	}

	// DefaultMarketConfig is the default market configuration for OKX.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"ATOM/USD": {
				Ticker:       "ATOM-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
			},
			"ATOM/USDC": {
				Ticker:       "ATOM-USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDC"),
			},
			"ATOM/USDT": {
				Ticker:       "ATOM-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USD": {
				Ticker:       "AVAX-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDC": {
				Ticker:       "AVAX-USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDC"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USD": {
				Ticker:       "BTC-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDC": {
				Ticker:       "BTC-USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTC-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIA-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"CELESTIA/USDT": {
				Ticker:       "TIA-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"DYDX/USD": {
				Ticker:       "DYDX-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDX-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH-BTC",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETH-USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"SOLANA/USD": {
				Ticker:       "SOL-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOL-USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDC/USD": {
				Ticker:       "USDC-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
			},
			"USDC/USDT": {
				Ticker:       "USDC-USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
			},
			"USDT/USD": {
				Ticker:       "USDT-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
			},
		},
	}
)
