package gate

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the Gate.io provider.
	Name = "gate_ws"
	// URL is the base url of for the Gate.io websocket API.
	URL = "wss://api.gateio.ws/ws/v4/"
	// DefaultReconnectionTimeout is the default reconnection timeout for the Gate.io websocket.
	DefaultReconnectionTimeout = 10 * time.Second
)

// DefaultWebSocketConfig is the default configuration for the Gate.io Websocket.
var (
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           DefaultReconnectionTimeout,
		PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
		Endpoints:                     []config.Endpoint{{URL: URL}},
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteInterval:                 config.DefaultWriteInterval,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}
)
