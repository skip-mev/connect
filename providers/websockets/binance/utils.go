package binance

import (
	"github.com/skip-mev/slinky/oracle/config"
)

var (
	// Name is the name of the Binance exchange WebSocket provider.
	Name = "binance_ws"
	// WSS is the WSS for the Binance exchange WebSocket API.
	WSS = "wss://stream.binance.com/stream"
	// DefaultMaxSubscriptionsPerConnection is the default maximum number of subscriptions
	// per connection. By default, Binance accepts up to 1024 subscriptions per connection.
	// However, we limit this to 20 to prevent overloading the connection.
	DefaultMaxSubscriptionsPerConnection = 100
)

var (
	// DefaultWebSocketConfig is the default configuration for the Binance exchange WebSocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		Endpoints:                     []config.Endpoint{{URL: WSS}},
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
)
