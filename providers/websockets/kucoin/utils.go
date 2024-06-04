package kucoin

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
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
	// per connection for the KuCoin websocket.
	//
	// ref: https://www.kucoin.com/docs/basic-info/request-rate-limit/websocket
	DefaultMaxSubscriptionsPerConnection = 50

	// DefaultWriteInterval is the default write interval for the KuCoin websocket.
	DefaultWriteInterval = 10 * time.Millisecond

	// DefaultReconnectionTimeout is the default reconnection timeout for the KuCoin websocket.
	DefaultReconnectionTimeout = 1 * time.Minute
)

var (
	// DefaultWebSocketConfig defines the default websocket config for Kucoin.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           DefaultReconnectionTimeout,
		PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
		Endpoints:                     []config.Endpoint{{URL: WSS}},
		Name:                          Name,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteInterval:                 DefaultWriteInterval,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: DefaultMaxSubscriptionsPerConnection,
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
