package gate

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// Name is the name of the Gate.io provider.
	Name = "gate_ws"
	// URL is the base url of for the Gate.io websocket API.
	URL = "wss://api.gateio.ws/ws/v4/"
)

// DefaultWebSocketConfig is the default configuration for the Gate.io Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 1000,
	ReconnectionTimeout:           10 * time.Second,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: URL}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 config.DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
}
