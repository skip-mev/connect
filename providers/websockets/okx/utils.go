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

	Type = types.ConfigType

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

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for OKX.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USDC: {
			OffChainTicker: "APE-USDC",
		},
		constants.APE_USDT: {
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USDC: {
			OffChainTicker: "APT-USDC",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "APT-USDT",
		},
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARB-USDT",
		},
		constants.ATOM_USD: {
			OffChainTicker: "ATOM-USD",
		},
		constants.ATOM_USDC: {
			OffChainTicker: "ATOM-USDC",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USD: {
			OffChainTicker: "AVAX-USD",
		},
		constants.AVAX_USDC: {
			OffChainTicker: "AVAX-USDC",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCH-USDT",
		},
		constants.BITCOIN_USD: {
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "BTC-USDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTC-USDT",
		},
		constants.BLUR_USDT: {
			OffChainTicker: "BLUR-USDT",
		},
		constants.CARDANO_USD: {
			OffChainTicker: "ADA-USD",
		},
		constants.CARDANO_USDC: {
			OffChainTicker: "ADA-USDC",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADA-USDT",
		},
		constants.CELESTIA_USD: {
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDT: {
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINK-USDT",
		},
		constants.COMPOUND_USDT: {
			OffChainTicker: "COMP-USDT",
		},
		constants.CURVE_USDT: {
			OffChainTicker: "CRV-USDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGE-USDT",
		},
		constants.DYDX_USD: {
			OffChainTicker: "DYDX-USD",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDX-USDT",
		},
		constants.ETC_USDT: {
			OffChainTicker: "ETC-USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USD: {
			OffChainTicker: "ETH-USD",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETH-USDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETH-USDT",
		},
		constants.FILECOIN_USDT: {
			OffChainTicker: "FIL-USDT",
		},
		constants.LIDO_USDT: {
			OffChainTicker: "LDO-USDT",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTC-USDT",
		},
		constants.MAKER_USDT: {
			OffChainTicker: "MKR-USDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOT-USDT",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATIC-USDT",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "NEAR-USDT",
		},
		constants.OPTIMISM_USDT: {
			OffChainTicker: "OP-USDT",
		},
		constants.PEPE_USDT: {
			OffChainTicker: "PEPE-USDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRP-USDT",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIB-USDT",
		},
		constants.SOLANA_USD: {
			OffChainTicker: "SOL-USD",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOL-USDT",
		},
		constants.STELLAR_USDT: {
			OffChainTicker: "XLM-USDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUI-USDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRX-USDT",
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: "UNI-USDT",
		},
		constants.USDC_USD: {
			OffChainTicker: "USDC-USD",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			OffChainTicker: "USDT-USD",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "WLD-USDT",
		},
	}
)
