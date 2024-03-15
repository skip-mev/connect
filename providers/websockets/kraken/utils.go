package kraken

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// URL is the websocket URL for Kraken. You can find the documentation here:
	// https://docs.kraken.com/websockets/. Kraken provides an authenticated and
	// unauthenticated websocket. The URLs defined below are all unauthenticated.

	// Name is the name of the Kraken provider.
	Name = "Kraken"

	// URL is the production websocket URL for Kraken.
	URL = "wss://ws.kraken.com"

	// URL_BETA is the demo websocket URL for Kraken.
	URL_BETA = "wss://beta-ws.kraken.com"
)

var (
	// DefaultWebSocketConfig is the default configuration for the Kraken Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           10 * time.Second,
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
		Type:                          types.ConfigType,
	}

	// DefaultMarketConfig is the default market configuration for Kraken.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APE_USD: {
			Name:           Name,
			OffChainTicker: "APE/USD",
		},
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOM/USD",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAX/USD",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX/USDT",
		},
		constants.BCH_USD: {
			Name:           Name,
			OffChainTicker: "BCH/USD",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "XBT/USD",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "XBT/USDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "XBT/USDT",
		},
		constants.BLUR_USD: {
			Name:           Name,
			OffChainTicker: "BLUR/USD",
		},
		constants.CARDANO_USD: {
			Name:           Name,
			OffChainTicker: "ADA/USD",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "TIA/USD",
		},
		constants.CHAINLINK_USD: {
			Name:           Name,
			OffChainTicker: "LINK/USD",
		},
		constants.COMPOUND_USD: {
			Name:           Name,
			OffChainTicker: "COMP/USD",
		},
		constants.CURVE_USD: {
			Name:           Name,
			OffChainTicker: "CRV/USD",
		},
		constants.DOGE_USD: {
			Name:           Name,
			OffChainTicker: "XDG/USD",
		},
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDX/USD",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH/XBT",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ETH/USD",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETH/USDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH/USDT",
		},
		constants.FILECOIN_USD: {
			Name:           Name,
			OffChainTicker: "FIL/USD",
		},
		constants.LIDO_USD: {
			Name:           Name,
			OffChainTicker: "LDO/USD",
		},
		constants.LITECOIN_USD: {
			Name:           Name,
			OffChainTicker: "XLTCZ/USD",
		},
		constants.MAKER_USD: {
			Name:           Name,
			OffChainTicker: "MKR/USD",
		},
		constants.PEPE_USD: {
			Name:           Name,
			OffChainTicker: "PEPE/USD",
		},
		constants.POLKADOT_USD: {
			Name:           Name,
			OffChainTicker: "DOT/USD",
		},
		constants.POLYGON_USD: {
			Name:           Name,
			OffChainTicker: "MATIC/USD",
		},
		constants.RIPPLE_USD: {
			Name:           Name,
			OffChainTicker: "XXRPZ/USD",
		},
		constants.SHIBA_USD: {
			Name:           Name,
			OffChainTicker: "SHIB/USD",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOL/USD",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL/USDT",
		},
		constants.STELLAR_USD: {
			Name:           Name,
			OffChainTicker: "XXLMZ/USD",
		},
		constants.TRON_USD: {
			Name:           Name,
			OffChainTicker: "TRX/USD",
		},
		constants.UNISWAP_USD: {
			Name:           Name,
			OffChainTicker: "UNI/USD",
		},
		constants.USDC_USD: {
			Name:           Name,
			OffChainTicker: "USDC/USD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC/USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT/USD",
		},
	}
)
