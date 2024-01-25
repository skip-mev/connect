package okx

import "github.com/skip-mev/slinky/oracle/config"

const (
	// OKX provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://www.okx.com/docs-v5/en/?shell#overview-production-trading-services
	// The two production URLs are defined in ProductionURL and ProductionAWSURL. The
	// DemoURL is used for testing purposes.

	// URL_PROD is the public OKX Websocket URL.
	URL_PROD = "wss://ws.okx.com:8443/ws/v5/public" //nolint

	// URL_PROD_AWS is the public OKX Websocket URL hosted on AWS.
	URL_PROD_AWS = "wss://wsaws.okx.com:8443/ws/v5/public" //nolint

	// DemoURL is the public OKX Websocket URL for test usage.
	URL_DEMO = "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999" //nolint
)

// DefaultWebSocketConfig is the default configuration for the OKX Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                Name,
	Enabled:             true,
	MaxBufferSize:       1000,
	ReconnectionTimeout: config.DefaultReconnectionTimeout,
	WSS:                 URL_PROD,
	ReadBufferSize:      config.DefaultReadBufferSize,
	WriteBufferSize:     config.DefaultWriteBufferSize,
	HandshakeTimeout:    config.DefaultHandshakeTimeout,
	EnableCompression:   config.DefaultEnableCompression,
	ReadTimeout:         config.DefaultReadTimeout,
	WriteTimeout:        config.DefaultWriteTimeout,
}
