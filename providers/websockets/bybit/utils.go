package bybit

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// ByBit provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://bybit-exchange.github.io/docs/v5/ws/connect
	// The two production URLs are defined in ProductionURL and TestnetURL.

	// Name is the name of the ByBit provider.
	Name = "bybit_ws"

	// URLProd is the public ByBit Websocket URL.
	URLProd = "wss://stream.bybit.com/v5/public/spot"

	// URLTest is the public testnet ByBit Websocket URL.
	URLTest = "wss://stream-testnet.bybit.com/v5/public/spot"

	// DefaultPingInterval is the default ping interval for the ByBit websocket.
	DefaultPingInterval = 15 * time.Second
)

// DefaultWebSocketConfig is the default configuration for the ByBit Websocket.
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
	PingInterval:                  DefaultPingInterval,
	WriteInterval:                 config.DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
}
