package mexc

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Please refer to the following link for the MEXC websocket documentation:
	// https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-market-streams.

	// Name is the name of the MEXC provider.
	Name = "mexc"

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
	MaxSubscriptionsPerConnection = 30
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
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for the MEXC Websocket.
	DefaultMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USDC": {
				Ticker:         constants.ATOM_USDC,
				OffChainTicker: "ATOMUSDC",
			},
			"ATOM/USDT": {
				Ticker:         constants.ATOM_USDT,
				OffChainTicker: "ATOMUSDT",
			},
			"AVAX/USDC": {
				Ticker:         constants.AVAX_USDC,
				OffChainTicker: "AVAXUSDC",
			},
			"AVAX/USDT": {
				Ticker:         constants.AVAX_USDT,
				OffChainTicker: "AVAXUSDT",
			},
			"BITCOIN/USDC": {
				Ticker:         constants.BITCOIN_USDC,
				OffChainTicker: "BTCUSDC",
			},
			"BITCOIN/USDT": {
				Ticker:         constants.BITCOIN_USDT,
				OffChainTicker: "BTCUSDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         constants.ETHEREUM_BITCOIN,
				OffChainTicker: "ETHBTC",
			},
			"ETHEREUM/USDC": {
				Ticker:         constants.ETHEREUM_USDC,
				OffChainTicker: "ETHUSDC",
			},
			"ETHEREUM/USDT": {
				Ticker:         constants.ETHEREUM_USDT,
				OffChainTicker: "ETHUSDT",
			},
			"SOLANA/USDC": {
				Ticker:         constants.SOLANA_USDC,
				OffChainTicker: "SOLUSDC",
			},
			"SOLANA/USDT": {
				Ticker:         constants.SOLANA_USDT,
				OffChainTicker: "SOLUSDT",
			},
			"USDC/USDT": {
				Ticker:         constants.USDC_USDT,
				OffChainTicker: "USDCUSDT",
			},
		},
	}
)
