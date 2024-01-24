package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// The following URLs are used for the Coinbase Websocket feed. These can be found
	// in the Coinbase documentation at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// URL is the production Coinbase Websocket URL.
	URL = "wss://ws-feed.exchange.coinbase.com"

	// URL_SANDBOX is the sandbox Coinbase Websocket URL.
	URL_SANDBOX = "wss://ws-feed-public.sandbox.exchange.coinbase.com"
)

const (
	// The following web socket configuration values were taken from the Coinbase documentation
	// at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// DefaultEnabledCompression is the default enabled compression for the Coinbase Websocket.
	// It is recommended to set this as true (please see the Coinbase documentation for more).
	DefaultEnabledCompression = false

	// DefaultWriteTimeout is the default write timeout for the Coinbase Websocket.
	// As recommended by Coinbase, this is set to 5 seconds.
	DefaultWriteTimeout = 5 * time.Second
)

var (
	// DefaultWebSocketConfig is the default configuration for the Coinbase Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:             true,
		Name:                Name,
		MaxBufferSize:       config.DefaultMaxBufferSize,
		ReconnectionTimeout: config.DefaultReconnectionTimeout,
		WSS:                 URL,
		ReadBufferSize:      config.DefaultReadBufferSize,
		WriteBufferSize:     config.DefaultWriteBufferSize,
		HandshakeTimeout:    config.DefaultHandshakeTimeout,
		EnableCompression:   DefaultEnabledCompression,
		WriteTimeout:        DefaultWriteTimeout,
		ReadTimeout:         config.DefaultReadTimeout,
	}
)
