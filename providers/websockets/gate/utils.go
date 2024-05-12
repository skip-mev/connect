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

	Type = types.ConfigType

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

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for Gate.io.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARB_USDT",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOM_USDT",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAX_USDT",
		},
		constants.APE_USDT: {
			OffChainTicker: "APE_USDT",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "APT_USDT",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCH_USDT",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTC_USDT",
		},
		constants.BLUR_USDT: {
			OffChainTicker: "BLUR_USDT",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADA_USDT",
		},
		constants.CELESTIA_USDT: {
			OffChainTicker: "TIA_USDT",
		},
		constants.COMPOUND_USDT: {
			OffChainTicker: "COMP_USDT",
		},
		constants.CURVE_USDT: {
			OffChainTicker: "CRV_USDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGE_USDT",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDX_USDT",
		},
		constants.ETC_USDT: {
			OffChainTicker: "ETC_USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETH_BTC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETH_USDT",
		},
		constants.FILECOIN_USDT: {
			OffChainTicker: "FIL_USDT",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "NEAR_USDT",
		},
		constants.OPTIMISM_USDT: {
			OffChainTicker: "OP_USDT",
		},
		constants.PEPE_USDT: {
			OffChainTicker: "PEPE_USDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOT_USDT",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATIC_USDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRP_USDT",
		},
		constants.SEI_USDT: {
			OffChainTicker: "SEI_USDT",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIB_USDT",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOL_USDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOL_USDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUI_USDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRX_USDT",
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: "UNI_USDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDC_USDT",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "WLD_USDT",
		},
	}
)
