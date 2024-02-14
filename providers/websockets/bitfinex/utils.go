package bitfinex

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Name is the name of the BitFinex provider.
	Name = "bitfinex"

	// URLProd is the public BitFinex Websocket URL.
	URLProd = "wss://api-pub.bitfinex.com/ws/2"

	MaxSubscriptionsPerConnection = 30
)

var (
	// DefaultWebSocketConfig is the default configuration for the BitFinex Websocket.
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
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for BitFinex.
	DefaultMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"BITCOIN/USD": {
				Ticker:         constants.BITCOIN_USD,
				OffChainTicker: "BTCUSD",
			},
			"CELESTIA/USD": {
				Ticker:         constants.CELESTIA_USD,
				OffChainTicker: "TIAUSD",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         constants.ETHEREUM_BITCOIN,
				OffChainTicker: "ETHBTC",
			},
			"ETHEREUM/USD": {
				Ticker:         constants.ETHEREUM_USD,
				OffChainTicker: "ETHUSD",
			},
			"SOLANA/USD": {
				Ticker:         constants.SOLANA_USD,
				OffChainTicker: "SOLUSD",
			},
		},
	}
)
