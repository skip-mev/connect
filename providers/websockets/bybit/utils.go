package bybit

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// ByBit provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://bybit-exchange.github.io/docs/v5/ws/connect
	// The two production URLs are defined in ProductionURL and TestnetURL. The

	// URLProd is the public ByBit Websocket URL.
	URLProd = "wss://stream.bybit.com/v5/public/spot"

	// URLTest is the public testnet ByBit Websocket URL.
	URLTest = "wss://stream-testnet.bybit.com/v5/public/spot"
)

// DefaultWebSocketConfig is the default configuration for the ByBit Websocket.
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
	PingInterval:        15 * time.Second,
}
