package cryptodotcom

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// URL is the URL used to connect to the Crypto.com websocket API. This can be found here
	// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#websocket-root-endpoints
	// Note that Crypto.com offers a sandbox and production environment.

	// Name is the name of the Crypto.com provider.
	Name = "crypto_dot_com_ws"

	// URL_PROD is the URL used to connect to the Crypto.com production websocket API.
	URL_PROD = "wss://stream.crypto.com/exchange/v1/market"

	// URL_SANDBOX is the URL used to connect to the Crypto.com sandbox websocket API. This will
	// return static prices.
	URL_SANDBOX = "wss://uat-stream.3ona.co/exchange/v1/market"

	// DefaultMaxSubscriptionsPerConnection is the default maximum number of subscriptions per connection.
	// Crypto.com has a limit of 400 but we set it to 200 to be safe.
	//
	// ref: https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html#introduction-2
	DefaultMaxSubscriptionsPerConnection = 200

	// DefaultPostConnectionTimeout is the default timeout for post connection. This is the recommended behaviour
	// from the Crypto.com documentation.
	//
	// ref: https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html#introduction-2
	DefaultPostConnectionTimeout = 1 * time.Second
)

// DefaultWebSocketConfig is the default configuration for the Crypto.com Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: URL_PROD}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 config.DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: DefaultMaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
}
