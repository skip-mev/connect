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
	Name = "kraken_ws"

	Type = types.ConfigType

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
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for Kraken.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USD: {
			OffChainTicker: "APE/USD",
		},
		constants.ATOM_USD: {
			OffChainTicker: "ATOM/USD",
		},
		constants.AVAX_USD: {
			OffChainTicker: "AVAX/USD",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAX/USDT",
		},
		constants.BCH_USD: {
			OffChainTicker: "BCH/USD",
		},
		constants.BITCOIN_USD: {
			OffChainTicker: "XBT/USD",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "XBT/USDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "XBT/USDT",
		},
		constants.BLUR_USD: {
			OffChainTicker: "BLUR/USD",
		},
		constants.CARDANO_USD: {
			OffChainTicker: "ADA/USD",
		},
		constants.CELESTIA_USD: {
			OffChainTicker: "TIA/USD",
		},
		constants.CHAINLINK_USD: {
			OffChainTicker: "LINK/USD",
		},
		constants.COMPOUND_USD: {
			OffChainTicker: "COMP/USD",
		},
		constants.CURVE_USD: {
			OffChainTicker: "CRV/USD",
		},
		constants.DOGE_USD: {
			OffChainTicker: "XDG/USD",
		},
		constants.DYDX_USD: {
			OffChainTicker: "DYDX/USD",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETH/XBT",
		},
		constants.ETHEREUM_USD: {
			OffChainTicker: "ETH/USD",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETH/USDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETH/USDT",
		},
		constants.FILECOIN_USD: {
			OffChainTicker: "FIL/USD",
		},
		constants.LIDO_USD: {
			OffChainTicker: "LDO/USD",
		},
		constants.LITECOIN_USD: {
			OffChainTicker: "XLTCZ/USD",
		},
		constants.MAKER_USD: {
			OffChainTicker: "MKR/USD",
		},
		constants.PEPE_USD: {
			OffChainTicker: "PEPE/USD",
		},
		constants.POLKADOT_USD: {
			OffChainTicker: "DOT/USD",
		},
		constants.POLYGON_USD: {
			OffChainTicker: "MATIC/USD",
		},
		constants.RIPPLE_USD: {
			OffChainTicker: "XXRPZ/USD",
		},
		constants.SHIBA_USD: {
			OffChainTicker: "SHIB/USD",
		},
		constants.SOLANA_USD: {
			OffChainTicker: "SOL/USD",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOL/USDT",
		},
		constants.STELLAR_USD: {
			OffChainTicker: "XXLMZ/USD",
		},
		constants.TRON_USD: {
			OffChainTicker: "TRX/USD",
		},
		constants.UNISWAP_USD: {
			OffChainTicker: "UNI/USD",
		},
		constants.USDC_USD: {
			OffChainTicker: "USDC/USD",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDC/USDT",
		},
		constants.USDT_USD: {
			OffChainTicker: "USDT/USD",
		},
	}
)
