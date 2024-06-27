package okx

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// OKX provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://www.okx.com/docs-v5/en/?shell#overview-production-trading-services
	// The two production URLs are defined in ProductionURL and ProductionAWSURL. The
	// DemoURL is used for testing purposes.

	// Name is the name of the OKX provider.
	Name = "okx_ws"

	// URL_PROD is the public OKX Websocket URL.
	URL_PROD = "wss://ws.okx.com:8443/ws/v5/public"

	// URL_PROD_AWS is the public OKX Websocket URL hosted on AWS.
	URL_PROD_AWS = "wss://wsaws.okx.com:8443/ws/v5/public"

	// URL_DEMO is the public OKX Websocket URL for test usage.
	URL_DEMO = "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999"

	// WriteInterval is the interval at which the OKX Websocket will write to the connection.
	// By default, there can be 3 messages written to the connection every second. Or 480
	// messages every hour.
	//
	// ref: https://www.okx.com/docs-v5/en/#overview-websocket-overview
	WriteInterval = 400 * time.Millisecond

	// MaxSubscriptionsPerConnection is the maximum number of subscriptions that can be
	// assigned to a single connection for the OKX provider. By default the limit is
	// 20 subscriptions per connection. We set the limit to 15 to be safe.
	//
	// ref: https://www.okx.com/docs-v5/en/#overview-websocket-overview
	MaxSubscriptionsPerConnection = 15
)

// DefaultWebSocketConfig is the default configuration for the OKX Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: URL_PROD}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 WriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: MaxSubscriptionsPerConnection,
}
