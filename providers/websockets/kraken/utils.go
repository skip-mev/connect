package kraken

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// URL is the websocket URL for Kraken. You can find the documentation here:
	// https://docs.kraken.com/websockets/. Kraken provides an authenticated and
	// unauthenticated websocket. The URLs defined below are all unauthenticated.

	// Name is the name of the Kraken provider.
	Name = "kraken_ws"

	// URL is the production websocket URL for Kraken.
	URL = "wss://ws.kraken.com"

	// URL_BETA is the demo websocket URL for Kraken.
	URL_BETA = "wss://beta-ws.kraken.com"
)

// DefaultWebSocketConfig is the default configuration for the Kraken Websocket.
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
