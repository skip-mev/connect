package gate

import (
	"github.com/skip-mev/slinky/oracle/config"
	"time"
)

const (
	// Name is the name of the Gate.io provider.
	Name = "gate.io"
	// URL is the base url of for the Gate.io websocket API.
	URL = "wss://api.gateio.ws/ws/v4/"
)

// DefaultWebSocketConfig is the default configuration for the Gate.io Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                Name,
	Enabled:             true,
	MaxBufferSize:       1000,
	ReconnectionTimeout: 10 * time.Second,
	WSS:                 URL,
	ReadBufferSize:      config.DefaultReadBufferSize,
	WriteBufferSize:     config.DefaultWriteBufferSize,
	HandshakeTimeout:    config.DefaultHandshakeTimeout,
	EnableCompression:   config.DefaultEnableCompression,
	ReadTimeout:         config.DefaultReadTimeout,
	WriteTimeout:        config.DefaultWriteTimeout,
}
