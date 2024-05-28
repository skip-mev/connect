package binance

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

var (
	// Name is the name of the Binance exchange WebSocket provider.
	Name = "binance_ws"
	// WSS is the WSS for the Binance exchange WebSocket API.
	WSS = "wss://stream.binance.com/stream"
	// DefaultMaxSubscriptionsPerConnection is the default maximum number of subscriptions
	// per connection. By default, Binance accepts up to 1024 subscriptions per connection.
	// However, we limit this to 20 to prevent overloading the connection.
	DefaultMaxSubscriptionsPerConnection = 20
)

var (
	// DefaultWebSocketConfig is the default configuration for the Binance exchange WebSocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           WSS,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for Binance.
	// DefaultNonUSMarketConfig is the default market configuration for Binance.
	DefaultNonUSMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTCUSDT",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETHUSDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDCUSDT",
		},
	}
)
