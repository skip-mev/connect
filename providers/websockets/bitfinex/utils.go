package bitfinex

import (
	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// Name is the name of the BitFinex provider.
	Name = "bitfinex_ws"

	// URLProd is the public BitFinex Websocket URL.
	URLProd = "wss://api-pub.bitfinex.com/ws/2"

	// DefaultMaxSubscriptionsPerConnection is the default maximum number of subscriptions
	// per connection. By default, BitFinex accepts up to 30 subscriptions per connection.
	// However, we limit this to 20 to prevent overloading the connection.
	DefaultMaxSubscriptionsPerConnection = 20
)

// DefaultWebSocketConfig is the default configuration for the BitFinex Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 1000,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: URLProd}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 config.DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: DefaultMaxSubscriptionsPerConnection,
	// Note that BitFinex does not support batch subscriptions. As such each new
	// market will be subscribed to with a new message.
	MaxSubscriptionsPerBatch: config.DefaultMaxSubscriptionsPerBatch,
}
