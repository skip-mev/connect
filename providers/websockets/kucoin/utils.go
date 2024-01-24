package kucoin

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// WSSEndpoint contains the endpoint format for Kucoin websocket API. Specifically
	// this inputs the dynamically generated token from the user and the endpoint.
	WSSEndpoint = "%s?token=%s"

	// WSS is the websocket URL for Kucoin. Note that this may change as the URL is
	WSS = "wss://ws-api-spot.kucoin.com/"

	// URL is the Kucoin websocket URL. This URL specifically points to the public
	// spot and maring REST API.
	URL = "https://api.kucoin.com"
)

// DefaultWebSocketConfig defines the default websocket config for Kucoin.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Enabled:             true,
	MaxBufferSize:       config.DefaultMaxBufferSize,
	ReconnectionTimeout: config.DefaultReconnectionTimeout,
	WSS:                 WSS, // Note that this may change as the URL is dynamically generated.
	Name:                Name,
	ReadBufferSize:      config.DefaultReadBufferSize,
	WriteBufferSize:     config.DefaultWriteBufferSize,
	HandshakeTimeout:    config.DefaultHandshakeTimeout,
	EnableCompression:   config.DefaultEnableCompression,
	ReadTimeout:         config.DefaultReadTimeout,
	WriteTimeout:        config.DefaultWriteTimeout,
	PingInterval:        10 * time.Second,
}

// DefaultAPIConfig defines the default API config for Kucoin. This is
// only utilized on the initial connection to the websocket feed.
var DefaultAPIConfig = config.APIConfig{
	Enabled:    false,
	Timeout:    5 * time.Second, // Kucoin recommends a timeout of 5 seconds.
	Interval:   1 * time.Minute, // This is not used.
	MaxQueries: 1,               // This is not used.
	URL:        URL,
	Name:       Name,
}
