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
	Name = "okx_ws"

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
	DefaultMarketConfig = types.TickersToProviderTickers{
		constants.APE_USDC: {
			Name:           Name,
			OffChainTicker: "APE-USDC",
		},
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USDC: {
			Name:           Name,
			OffChainTicker: "APT-USDC",
		},
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "APT-USDT",
		},
		constants.ARBITRUM_USDT: {
			Name:           Name,
			OffChainTicker: "ARB-USDT",
		},
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
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCH-USDT",
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
		constants.BLUR_USDT: {
			Name:           Name,
			OffChainTicker: "BLUR-USDT",
		},
		constants.CARDANO_USD: {
			Name:           Name,
			OffChainTicker: "ADA-USD",
		},
		constants.CARDANO_USDC: {
			Name:           Name,
			OffChainTicker: "ADA-USDC",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADA-USDT",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USDT: {
			Name:           Name,
			OffChainTicker: "LINK-USDT",
		},
		constants.COMPOUND_USDT: {
			Name:           Name,
			OffChainTicker: "COMP-USDT",
		},
		constants.CURVE_USDT: {
			Name:           Name,
			OffChainTicker: "CRV-USDT",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "DOGE-USDT",
		},
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDX-USD",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDX-USDT",
		},
		constants.ETC_USDT: {
			Name:           Name,
			OffChainTicker: "ETC-USDT",
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
		constants.FILECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "FIL-USDT",
		},
		constants.LIDO_USDT: {
			Name:           Name,
			OffChainTicker: "LDO-USDT",
		},
		constants.LITECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "LTC-USDT",
		},
		constants.MAKER_USDT: {
			Name:           Name,
			OffChainTicker: "MKR-USDT",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOT-USDT",
		},
		constants.POLYGON_USDT: {
			Name:           Name,
			OffChainTicker: "MATIC-USDT",
		},
		constants.NEAR_USDT: {
			Name:           Name,
			OffChainTicker: "NEAR-USDT",
		},
		constants.OPTIMISM_USDT: {
			Name:           Name,
			OffChainTicker: "OP-USDT",
		},
		constants.PEPE_USDT: {
			Name:           Name,
			OffChainTicker: "PEPE-USDT",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRP-USDT",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIB-USDT",
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
		constants.STELLAR_USDT: {
			Name:           Name,
			OffChainTicker: "XLM-USDT",
		},
		constants.SUI_USDT: {
			Name:           Name,
			OffChainTicker: "SUI-USDT",
		},
		constants.TRON_USDT: {
			Name:           Name,
			OffChainTicker: "TRX-USDT",
		},
		constants.UNISWAP_USDT: {
			Name:           Name,
			OffChainTicker: "UNI-USDT",
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
		constants.WORLD_USDT: {
			Name:           Name,
			OffChainTicker: "WLD-USDT",
		},
	}
)
