package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
)

func TestWebSocketConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.WebSocketConfig
		expectedErr bool
	}{
		{
			name: "good config with websocket enabled",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   config.DefaultReconnectionTimeout,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: false,
		},
		{
			name: "good config with websocket disabled",
			config: config.WebSocketConfig{
				Enabled: false,
			},
			expectedErr: false,
		},
		{
			name: "bad config with no max buffer size",
			config: config.WebSocketConfig{
				Enabled:               true,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no reconnection timeout",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no name",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no wss",
			config: config.WebSocketConfig{
				Enabled:                       true,
				MaxBufferSize:                 1,
				ReconnectionTimeout:           time.Second,
				PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
				Name:                          "test",
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative read buffer size",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                -1,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative write buffer size",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               -1,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no handshake timeout",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},

		{
			name: "bad config with no read timeout",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no write timeout",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   time.Second,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad ping interval",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   config.DefaultReconnectionTimeout,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  -1,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad max subscriptions per connection",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   config.DefaultReconnectionTimeout,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: -1,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative max read error count",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   config.DefaultReconnectionTimeout,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             -1,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative write interval",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   config.DefaultReconnectionTimeout,
				PostConnectionTimeout: config.DefaultPostConnectionTimeout,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 -1,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative post connection timeout",
			config: config.WebSocketConfig{
				Enabled:               true,
				MaxBufferSize:         1,
				ReconnectionTimeout:   config.DefaultReconnectionTimeout,
				PostConnectionTimeout: -1,
				Name:                  "test",
				Endpoints: []config.Endpoint{
					{
						URL: "wss://test.com",
					},
				},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				PingInterval:                  config.DefaultPingInterval,
				WriteInterval:                 config.DefaultWriteInterval,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
			},
			expectedErr: true,
		},
		{
			name: "bad config with less than 1 max subscriptions per batch",
			config: config.WebSocketConfig{
				Enabled:                       true,
				MaxBufferSize:                 1,
				ReconnectionTimeout:           config.DefaultReconnectionTimeout,
				PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
				Name:                          "test",
				Endpoints:                     []config.Endpoint{{URL: "wss://test.com"}},
				ReadBufferSize:                config.DefaultReadBufferSize,
				WriteBufferSize:               config.DefaultWriteBufferSize,
				HandshakeTimeout:              config.DefaultHandshakeTimeout,
				EnableCompression:             config.DefaultEnableCompression,
				ReadTimeout:                   config.DefaultReadTimeout,
				WriteTimeout:                  config.DefaultWriteTimeout,
				MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
				MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
				MaxSubscriptionsPerBatch:      0,
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
