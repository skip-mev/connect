package mexc

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Please refer to the following link for the MEXC websocket documentation:
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-market-streams.

	// Name is the name of the MEXC provider.
	Name = "mexc_ws"

	Type = types.ConfigType

	// WSS is the public MEXC Websocket URL.
	WSS = "wss://wbs.mexc.com/ws"

	// DefaultPingInterval is the default ping interval for the MEXC websocket. The documentation
	// specifies that this should be done every 30 seconds, however, the actual threshold should be
	// slightly lower than this to account for network latency.
	DefaultPingInterval = 20 * time.Second

	// MaxSubscriptionsPerConnection is the maximum number of subscriptions that can be made
	// per connection.
	//
	// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-market-streams
	MaxSubscriptionsPerConnection = 20
)

var (
	// DefaultWebSocketConfig is the default configuration for the MEXC Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           WSS,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: MaxSubscriptionsPerConnection,
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for the MEXC Websocket.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USDT: {
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDC: {
			OffChainTicker: "ATOMUSDC",
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
		constants.BLUR_USDT: {
			OffChainTicker: "BLURUSDT",
		},
		constants.CARDANO_USDC: {
			OffChainTicker: "ADAUSDC",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINKUSDT",
		},
		constants.COMPOUND_USDT: {
			OffChainTicker: "COMPUSDT",
		},
		constants.CURVE_USDT: {
			OffChainTicker: "CRVUSDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGEUSDT",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDXUSDT",
		},
		constants.ETC_USDT: {
			OffChainTicker: "ETCUSDT",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETHUSDT",
		},
		constants.FILECOIN_USDT: {
			OffChainTicker: "FILUSDT",
		},
		constants.LIDO_USDT: {
			OffChainTicker: "LDOUSDT",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTCUSDT",
		},
		constants.MAKER_USDT: {
			OffChainTicker: "MKRUSDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOTUSDT",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "NEARUSDT",
		},
		constants.OPTIMISM_USDT: {
			OffChainTicker: "OPUSDT",
		},
		constants.PEPE_USDT: {
			OffChainTicker: "PEPEUSDT",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATICUSDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRPUSDT",
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
		constants.STELLAR_USDT: {
			OffChainTicker: "XLMUSDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRXUSDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDCUSDT",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "WLDUSDT",
		},
	}
)
