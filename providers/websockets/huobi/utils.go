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
	Name = "huobi_ws"

	Type = types.ConfigType

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

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for the Huobi Websocket.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.ARBITRUM_USDT: {
			OffChainTicker: "arbusdt",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "atomusdt",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "avaxusdt",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "aptusdt",
		},
		constants.BCH_USDT: {
			OffChainTicker: "bchusdt",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "btcusdc",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "btcusdt",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "adausdt",
		},
		constants.CELESTIA_USDT: {
			OffChainTicker: "tiausdt",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "dogeusdt",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "dydxusdt",
		},
		constants.ETC_USDT: {
			OffChainTicker: "etcusdt",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ethbtc",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ethusdc",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ethusdt",
		},
		constants.FILECOIN_USDT: {
			OffChainTicker: "filusdt",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "ltcusdt",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "nearusdt",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "maticusdt",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "xrpusdt",
		},
		constants.SEI_USDT: {
			OffChainTicker: "seiusdt",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "solusdt",
		},
		constants.SUI_USDT: {
			OffChainTicker: "suiusdt",
		},
		constants.TRON_USDT: {
			OffChainTicker: "trxusdt",
		},
		constants.USDC_USDT: {
			OffChainTicker: "usdcusdt",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "wldusdt",
		},
	}
)
