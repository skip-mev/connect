package okx

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
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
	WriteInterval = 3000 * time.Millisecond

	// MaxSubscriptionsPerConnection is the maximum number of subscriptions that can be
	// assigned to a single connection for the OKX provider.
	//
	// ref: https://www.okx.com/docs-v5/en/#overview-websocket-overview
	MaxSubscriptionsPerConnection = 50

	// MaxSubscriptionsPerBatch is the maximum number of subscriptions that can be
	// assigned to a single batch for the OKX provider. We set the limit to 5 to be safe.
	MaxSubscriptionsPerBatch = 25

	// ReadTimeout is the timeout for reading from the OKX Websocket connection.
	ReadTimeout = 15 * time.Second
)

// DefaultWebSocketConfig is the default configuration for the OKX Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: URL_PROD_AWS}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   ReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 WriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: MaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      MaxSubscriptionsPerBatch,
}
