package okx

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// OKX provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://www.okx.com/docs-v5/en/?shell#overview-production-trading-services
	// The two production URLs are defined in ProductionURL and ProductionAWSURL. The
	// DemoURL is used for testing purposes.

	// Name is the name of the OKX provider.
	Name = "okx"

	// URL_PROD is the public OKX Websocket URL.
	URL_PROD = "wss://ws.okx.com:8443/ws/v5/public"

	// URL_PROD_AWS is the public OKX Websocket URL hosted on AWS.
	URL_PROD_AWS = "wss://wsaws.okx.com:8443/ws/v5/public"

	// URL_DEMO is the public OKX Websocket URL for test usage.
	URL_DEMO = "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999"
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
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOM-USD",
		},
		constants.ATOM_USDC: {
			Name:           Name,
			OffChainTicker: "ATOM-USDC",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAX-USD",
		},
		constants.AVAX_USDC: {
			Name:           Name,
			OffChainTicker: "AVAX-USDC",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX-USDT",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTC-USDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC-USDT",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA-USDT",
		},
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDX-USD",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDX-USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ETH-USD",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETH-USDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH-USDT",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOL-USD",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL-USDT",
		},
		constants.USDC_USD: {
			Name:           Name,
			OffChainTicker: "USDC-USD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT-USD",
		},
	}
)
