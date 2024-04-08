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

	// DefaultMarketConfig is the default market configuration for the MEXC Websocket.
	DefaultMarketConfig = types.TickersToProviderTickers{
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APEUSDT",
		},
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "APTUSDT",
		},
		constants.ARBITRUM_USDT: {
			Name:           Name,
			OffChainTicker: "ARBUSDT",
		},
		constants.ATOM_USDC: {
			Name:           Name,
			OffChainTicker: "ATOMUSDC",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOMUSDT",
		},
		constants.AVAX_USDC: {
			Name:           Name,
			OffChainTicker: "AVAXUSDC",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAXUSDT",
		},
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCHUSDT",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTCUSDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTCUSDT",
		},
		constants.BLUR_USDT: {
			Name:           Name,
			OffChainTicker: "BLURUSDT",
		},
		constants.CARDANO_USDC: {
			Name:           Name,
			OffChainTicker: "ADAUSDC",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADAUSDT",
		},
		constants.CHAINLINK_USDT: {
			Name:           Name,
			OffChainTicker: "LINKUSDT",
		},
		constants.COMPOUND_USDT: {
			Name:           Name,
			OffChainTicker: "COMPUSDT",
		},
		constants.CURVE_USDT: {
			Name:           Name,
			OffChainTicker: "CRVUSDT",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "DOGEUSDT",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDXUSDT",
		},
		constants.ETC_USDT: {
			Name:           Name,
			OffChainTicker: "ETCUSDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETHUSDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETHUSDT",
		},
		constants.FILECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "FILUSDT",
		},
		constants.LIDO_USDT: {
			Name:           Name,
			OffChainTicker: "LDOUSDT",
		},
		constants.LITECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "LTCUSDT",
		},
		constants.MAKER_USDT: {
			Name:           Name,
			OffChainTicker: "MKRUSDT",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOTUSDT",
		},
		constants.NEAR_USDT: {
			Name:           Name,
			OffChainTicker: "NEARUSDT",
		},
		constants.OPTIMISM_USDT: {
			Name:           Name,
			OffChainTicker: "OPUSDT",
		},
		constants.PEPE_USDT: {
			Name:           Name,
			OffChainTicker: "PEPEUSDT",
		},
		constants.POLYGON_USDT: {
			Name:           Name,
			OffChainTicker: "MATICUSDT",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRPUSDT",
		},
		constants.SEI_USDT: {
			Name:           Name,
			OffChainTicker: "SEIUSDT",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIBUSDT",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOLUSDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOLUSDT",
		},
		constants.STELLAR_USDT: {
			Name:           Name,
			OffChainTicker: "XLMUSDT",
		},
		constants.SUI_USDT: {
			Name:           Name,
			OffChainTicker: "SUIUSDT",
		},
		constants.TRON_USDT: {
			Name:           Name,
			OffChainTicker: "TRXUSDT",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDCUSDT",
		},
		constants.WORLD_USDT: {
			Name:           Name,
			OffChainTicker: "WLDUSDT",
		},
	}
)
