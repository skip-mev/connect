package huobi

import "github.com/skip-mev/slinky/oracle/config"

const (
	// Huobi provides the following URLS fro its Websocket API. More info can be found in the documentation
	// here: https://huobiapi.github.io/docs/spot/v1/en/#websocket-market-data.

	// URL is the public Huobi Websocket URL.
	URL = "wss://api.huobi.pro/ws"

	// URL_AWS is the public Huobi Websocket URL hosted on AWS.
	URL_AWS = "wss://api-aws.huobi.pro/ws"

	Name = "huobi"
)

// DefaultWebSocketConfig is the default configuration for the Huobi Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                Name,
	Enabled:             true,
	MaxBufferSize:       1000,
	ReconnectionTimeout: config.DefaultReconnectionTimeout,
	WSS:                 URL,
	ReadBufferSize:      config.DefaultReadBufferSize,
	WriteBufferSize:     config.DefaultWriteBufferSize,
	HandshakeTimeout:    config.DefaultHandshakeTimeout,
	EnableCompression:   config.DefaultEnableCompression,
	ReadTimeout:         config.DefaultReadTimeout,
	WriteTimeout:        config.DefaultWriteTimeout,
}
