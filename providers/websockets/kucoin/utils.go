package kucoin

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// Name is the name of the KuCoin provider.
	Name = "kucoin_ws"

	// WSSEndpoint contains the endpoint format for Kucoin websocket API. Specifically
	// this inputs the dynamically generated token from the user and the endpoint.
	WSSEndpoint = "%s?token=%s"

	// WSS is the websocket URL for Kucoin. Note that this may change as the URL is
	// dynamically generated. A token is required to connect to the websocket feed.
	WSS = "wss://ws-api-spot.kucoin.com/"

	// URL is the Kucoin websocket URL. This URL specifically points to the public
	// spot and maring REST API.
	URL = "https://api.kucoin.com"

	// DefaultPingInterval is the default ping interval for the KuCoin websocket.
	DefaultPingInterval = 10 * time.Second

	// DefaultMaxSubscriptionsPerConnection is the default maximum number of subscriptions
	// per connection. By default, KuCoin accepts up to 300 subscriptions per connection.
	// However we limit this to 50 to prevent overloading the connection.
	//
	// ref: https://www.kucoin.com/docs/basic-info/request-rate-limit/websocket
	DefaultMaxSubscriptionsPerConnection = 25

	// DefaultWriteInterval is the default write interval for the KuCoin websocket.
	// Kucoin allows 100 messages to be sent per 10 seconds. We set this to 300ms to
	// prevent overloading the connection.
	//
	// https://www.kucoin.com/docs/basic-info/request-rate-limit/websocket
	DefaultWriteInterval = 300 * time.Millisecond

	// DefaultHandShakeTimeout is the default handshake timeout for the KuCoin websocket.
	// Assuming that we can create 40 subscriptions every 7.5 seconds, we want to space
	// out the subscriptions to prevent overloading the connection. So we set the
	// handshake timeout to 20 seconds.
	DefaultHandShakeTimeout = 20 * time.Second
)

var (
	// DefaultWebSocketConfig defines the default websocket config for Kucoin.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
		Endpoints:                     []config.Endpoint{{URL: WSS}},
		Name:                          Name,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              DefaultHandShakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  DefaultPingInterval,
		WriteInterval:                 DefaultWriteInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: DefaultMaxSubscriptionsPerConnection,
		MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
	}

	// DefaultAPIConfig defines the default API config for KuCoin. This is
	// only utilized on the initial connection to the websocket feed.
	DefaultAPIConfig = config.APIConfig{
		Enabled:    false,
		Timeout:    5 * time.Second, // KuCoin recommends a timeout of 5 seconds.
		Interval:   1 * time.Minute, // This is not used.
		MaxQueries: 1,               // This is not used.
		Endpoints:  []config.Endpoint{{URL: URL}},
		Name:       Name,
	}
)
