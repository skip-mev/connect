package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
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
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: config.DefaultReconnectionTimeout,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
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
				Enabled:             true,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no reconnection timeout",
			config: config.WebSocketConfig{
				Enabled:           true,
				MaxBufferSize:     1,
				Name:              "test",
				WSS:               "wss://test.com",
				ReadBufferSize:    config.DefaultReadBufferSize,
				WriteBufferSize:   config.DefaultWriteBufferSize,
				HandshakeTimeout:  config.DefaultHandshakeTimeout,
				EnableCompression: config.DefaultEnableCompression,
				ReadTimeout:       config.DefaultReadTimeout,
				WriteTimeout:      config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no name",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no wss",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative read buffer size",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      -1,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with negative write buffer size",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     -1,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no handshake timeout",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},

		{
			name: "bad config with no read timeout",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				WriteTimeout:        config.DefaultWriteTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no write timeout",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: time.Second,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad ping interval",
			config: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1,
				ReconnectionTimeout: config.DefaultReconnectionTimeout,
				Name:                "test",
				WSS:                 "wss://test.com",
				ReadBufferSize:      config.DefaultReadBufferSize,
				WriteBufferSize:     config.DefaultWriteBufferSize,
				HandshakeTimeout:    config.DefaultHandshakeTimeout,
				EnableCompression:   config.DefaultEnableCompression,
				ReadTimeout:         config.DefaultReadTimeout,
				WriteTimeout:        config.DefaultWriteTimeout,
				PingInterval:        -1,
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
