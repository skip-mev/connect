package huobi

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Huobi provides the following URLs for its Websocket API. More info can be found in the documentation
	// here: https://huobiapi.github.io/docs/spot/v1/en/#websocket-market-data.

	// Name is the name of the Huobi provider.
	Name = "huobi"

	// URL is the public Huobi Websocket URL.
	URL = "wss://api.huobi.pro/ws"

	// URLAws is the public Huobi Websocket URL hosted on AWS.
	URLAws = "wss://api-aws.huobi.pro/ws"
)

var (
	// DefaultWebSocketConfig is the default configuration for the Huobi Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
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

	// DefaultMarketConfig is the default market configuration for the Huobi Websocket.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "aptusdt",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "atomusdt",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "avaxusdt",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "btcusdc",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "btcusdt",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "adausdt",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "tiausdt",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "dydxusdt",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ethbtc",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ethusdc",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ethusdt",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "solusdt",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "usdcusdt",
		},
	}
)
