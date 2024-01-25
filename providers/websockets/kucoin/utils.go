package kucoin

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// WSSEndpoint contains the endpoint format for Kucoin websocket API. Specifically
	// this inputs the dynamically generated token from the user and the endpoint.
	WSSEndpoint = "%s?token=%s"
)

// DefaultWebSocketConfig defines the default websocket config for Kucoin.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Enabled:             true,
	MaxBufferSize:       config.DefaultMaxBufferSize,
	ReconnectionTimeout: config.DefaultReconnectionTimeout,
	WSS:                 URL, // We use the REST Url as the websocket URL since it is dynamic.
	Name:                Name,
	ReadBufferSize:      config.DefaultReadBufferSize,
	WriteBufferSize:     config.DefaultWriteBufferSize,
	HandshakeTimeout:    config.DefaultHandshakeTimeout,
	EnableCompression:   config.DefaultEnableCompression,
	ReadTimeout:         config.DefaultReadTimeout,
	WriteTimeout:        config.DefaultWriteTimeout,
	PingInterval:        10 * time.Second,
}
