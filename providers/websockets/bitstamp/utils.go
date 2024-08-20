package bitstamp

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// Name is the name of the bitstamp provider.
	Name = "bitstamp_ws"

	// WSS is the bitstamp websocket address.
	WSS = "wss://ws.bitstamp.net"

	// DefaultPingInterval is the default ping interval for the bitstamp websocket.
	DefaultPingInterval = 10 * time.Second
)

// DefaultWebSocketConfig returns the default websocket config for bitstamp.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Enabled:                       true,
	Name:                          Name,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: WSS}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	WriteTimeout:                  config.DefaultWriteTimeout,
	ReadTimeout:                   config.DefaultReadTimeout,
	PingInterval:                  DefaultPingInterval,
	WriteInterval:                 config.DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	// Note that BitStamp does not support batch subscriptions. As such each new
	// market will be subscribed to with a new message.
	MaxSubscriptionsPerBatch: config.DefaultMaxSubscriptionsPerBatch,
}
