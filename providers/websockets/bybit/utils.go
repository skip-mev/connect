package bybit

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// ByBit provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://bybit-exchange.github.io/docs/v5/ws/connect
	// The two production URLs are defined in ProductionURL and TestnetURL.

	// Name is the name of the ByBit provider.
	Name = "bybit_ws"

	Type = types.ConfigType

	// URLProd is the public ByBit Websocket URL.
	URLProd = "wss://stream.bybit.com/v5/public/spot"

	// URLTest is the public testnet ByBit Websocket URL.
	URLTest = "wss://stream-testnet.bybit.com/v5/public/spot"

	// DefaultPingInterval is the default ping interval for the ByBit websocket.
	DefaultPingInterval = 15 * time.Second
)

var (
	// DefaultWebSocketConfig is the default configuration for the ByBit Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URLProd,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for ByBit.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APTOS_USDT: {
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOMUSDT",
		},
		constants.AVAX_USDC: {
			OffChainTicker: "AVAXUSDC",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAXUSDT",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCHUSDT",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "BTCUSDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTCUSDT",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINKUSDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGEUSDT",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDXUSDT",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETHUSDT",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTCUSDT",
		},
		constants.PEPE_USDT: {
			OffChainTicker: "PEPEUSDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOTUSDT",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATICUSDT",
		},
		constants.SEI_USDT: {
			OffChainTicker: "SEIUSDT",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIBUSDT",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOLUSDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOLUSDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRXUSDT",
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: "UNIUSDT",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "WLDUSDT",
		},
		constants.STELLAR_USDT: {
			OffChainTicker: "XLMUSDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRPUSDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDCUSDT",
		},
	}
)
