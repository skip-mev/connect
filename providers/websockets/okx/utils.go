package okx

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
	DefaultMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USD": {
				Ticker:         constants.ATOM_USD,
				OffChainTicker: "ATOM-USD",
			},
			"ATOM/USDC": {
				Ticker:         constants.ATOM_USDC,
				OffChainTicker: "ATOM-USDC",
			},
			"ATOM/USDT": {
				Ticker:         constants.ATOM_USDT,
				OffChainTicker: "ATOM-USDT",
			},
			"AVAX/USD": {
				Ticker:         constants.AVAX_USD,
				OffChainTicker: "AVAX-USD",
			},
			"AVAX/USDC": {
				Ticker:         constants.AVAX_USDC,
				OffChainTicker: "AVAX-USDC",
			},
			"AVAX/USDT": {
				Ticker:         constants.AVAX_USDT,
				OffChainTicker: "AVAX-USDT",
			},
			"BITCOIN/USD": {
				Ticker:         constants.BITCOIN_USD,
				OffChainTicker: "BTC-USD",
			},
			"BITCOIN/USDC": {
				Ticker:         constants.BITCOIN_USDC,
				OffChainTicker: "BTC-USDC",
			},
			"BITCOIN/USDT": {
				Ticker:         constants.BITCOIN_USDT,
				OffChainTicker: "BTC-USDT",
			},
			"CELESTIA/USD": {
				Ticker:         constants.CELESTIA_USD,
				OffChainTicker: "TIA-USD",
			},
			"CELESTIA/USDT": {
				Ticker:         constants.CELESTIA_USDT,
				OffChainTicker: "TIA-USDT",
			},
			"DYDX/USD": {
				Ticker:         constants.DYDX_USD,
				OffChainTicker: "DYDX-USD",
			},
			"DYDX/USDT": {
				Ticker:         constants.DYDX_USDT,
				OffChainTicker: "DYDX-USDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         constants.ETHEREUM_BITCOIN,
				OffChainTicker: "ETH-BTC",
			},
			"ETHEREUM/USD": {
				Ticker:         constants.ETHEREUM_USD,
				OffChainTicker: "ETH-USD",
			},
			"ETHEREUM/USDC": {
				Ticker:         constants.ETHEREUM_USDC,
				OffChainTicker: "ETH-USDC",
			},
			"ETHEREUM/USDT": {
				Ticker:         constants.ETHEREUM_USDT,
				OffChainTicker: "ETH-USDT",
			},
			"SOLANA/USD": {
				Ticker:         constants.SOLANA_USD,
				OffChainTicker: "SOL-USD",
			},
			"SOLANA/USDC": {
				Ticker:         constants.SOLANA_USDC,
				OffChainTicker: "SOL-USDC",
			},
			"SOLANA/USDT": {
				Ticker:         constants.SOLANA_USDT,
				OffChainTicker: "SOL-USDT",
			},
			"USDC/USD": {
				Ticker:         constants.USDC_USD,
				OffChainTicker: "USDC-USD",
			},
			"USDC/USDT": {
				Ticker:         constants.USDC_USDT,
				OffChainTicker: "USDC-USDT",
			},
			"USDT/USD": {
				Ticker:         constants.USDT_USD,
				OffChainTicker: "USDT-USD",
			},
		},
	}
)
