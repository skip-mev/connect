package binance

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

var (
	// Name is the name of the Binance exchange WebSocket provider.
	Name = "binance_ws"
	// WSS is the WSS for the Binance exchange WebSocket API.
	WSS = "wss://stream.binance.com/stream"
	// DefaultMaxSubscriptionsPerConnection is the default maximum number of subscriptions
	// per connection. By default, Binance accepts up to 1024 subscriptions per connection.
	// However, we limit this to 40 to prevent overloading the connection.
	DefaultMaxSubscriptionsPerConnection = 40
	// DefaultWriteInterval is the default write interval for the Binance exchange WebSocket.
	// Binance allows up to 5 messages to be sent per second. We set this to 300ms to
	// prevent overloading the connection.
	DefaultWriteInterval = 300 * time.Millisecond
	// DefaultHandshakeTimeout is the default handshake timeout for the Binance exchange WebSocket.
	// If we assume that for 20 markets it takes 250ms to write a message, then the handshake
	// timeout should be at least 5 seconds. We add a buffer of 5 seconds to account for network
	// latency.
	DefaultHandshakeTimeout = 20 * time.Second
)

// DefaultWebSocketConfig is the default configuration for the Binance exchange WebSocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	HandshakeTimeout:              DefaultHandshakeTimeout,
	Endpoints:                     []config.Endpoint{{URL: WSS}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: DefaultMaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
}
