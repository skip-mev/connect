package gate

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the Gate.io provider.
	Name = "gate_ws"
	// URL is the base url of for the Gate.io websocket API.
	URL = "wss://api.gateio.ws/ws/v4/"
)

// DefaultWebSocketConfig is the default configuration for the Gate.io Websocket.
var (
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
	}

	// DefaultProviderConfig is the default market configuration for Gate.io.
	DefaultProviderConfig = types.TickerToProviderConfig{
		constants.ARBITRUM_USDT: {
			Name:           Name,
			OffChainTicker: "ARB_USDT",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOM_USDT",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX_USDT",
		},
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APE_USDT",
		},
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "APT_USDT",
		},
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCH_USDT",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC_USDT",
		},
		constants.BLUR_USDT: {
			Name:           Name,
			OffChainTicker: "BLUR_USDT",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADA_USDT",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA_USDT",
		},
		constants.COMPOUND_USDT: {
			Name:           Name,
			OffChainTicker: "COMP_USDT",
		},
		constants.CURVE_USDT: {
			Name:           Name,
			OffChainTicker: "CRV_USDT",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "DOGE_USDT",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDX_USDT",
		},
		constants.ETC_USDT: {
			Name:           Name,
			OffChainTicker: "ETC_USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH_BTC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH_USDT",
		},
		constants.FILECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "FIL_USDT",
		},
		constants.NEAR_USDT: {
			Name:           Name,
			OffChainTicker: "NEAR_USDT",
		},
		constants.OPTIMISM_USDT: {
			Name:           Name,
			OffChainTicker: "OP_USDT",
		},
		constants.PEPE_USDT: {
			Name:           Name,
			OffChainTicker: "PEPE_USDT",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOT_USDT",
		},
		constants.POLYGON_USDT: {
			Name:           Name,
			OffChainTicker: "MATIC_USDT",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRP_USDT",
		},
		constants.SEI_USDT: {
			Name:           Name,
			OffChainTicker: "SEI_USDT",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIB_USDT",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOL_USDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL_USDT",
		},
		constants.SUI_USDT: {
			Name:           Name,
			OffChainTicker: "SUI_USDT",
		},
		constants.TRON_USDT: {
			Name:           Name,
			OffChainTicker: "TRX_USDT",
		},
		constants.UNISWAP_USDT: {
			Name:           Name,
			OffChainTicker: "UNI_USDT",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC_USDT",
		},
		constants.WORLD_USDT: {
			Name:           Name,
			OffChainTicker: "WLD_USDT",
		},
	}
)
