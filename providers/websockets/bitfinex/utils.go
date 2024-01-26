package bitfinex

import "github.com/skip-mev/slinky/oracle/config"

const (
	// Name is the name of the BitFinex provider.
	Name = "bitfinex"

	// URLProd is the public BitFinex Websocket URL.
	URLProd = "wss://api-pub.bitfinex.com/ws/2"
)

// DefaultWebSocketConfig is the default configuration for the BitFinex Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                Name,
	Enabled:             true,
	MaxBufferSize:       1000,
	ReconnectionTimeout: config.DefaultReconnectionTimeout,
	WSS:                 URLProd,
	ReadBufferSize:      config.DefaultReadBufferSize,
	WriteBufferSize:     config.DefaultWriteBufferSize,
	HandshakeTimeout:    config.DefaultHandshakeTimeout,
	EnableCompression:   config.DefaultEnableCompression,
	ReadTimeout:         config.DefaultReadTimeout,
	WriteTimeout:        config.DefaultWriteTimeout,
}
