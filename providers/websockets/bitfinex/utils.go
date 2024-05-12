package bitfinex

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the BitFinex provider.
	Name = "bitfinex_ws"

	Type = types.ConfigType

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

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for BitFinex.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.BITCOIN_USD: {
			OffChainTicker: "BTCUSD",
		},
		constants.CELESTIA_USD: {
			OffChainTicker: "TIAUSD",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETHBTC",
		},
		constants.ETHEREUM_USD: {
			OffChainTicker: "ETHUSD",
		},
		constants.SOLANA_USD: {
			OffChainTicker: "SOLUSD",
		},
	}
)
