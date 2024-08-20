package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
)

func TestProviderConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.ProviderConfig
		expectedErr bool
	}{
		{
			name: "good API config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:          true,
					Timeout:          time.Second,
					Interval:         time.Second,
					ReconnectTimeout: time.Second,
					MaxQueries:       1,
					Name:             "test",
					Atomic:           true,
					Endpoints:        []config.Endpoint{{URL: "http://test.com"}},
				},
				Name: "test",
				Type: "price_provider",
			},
			expectedErr: false,
		},
		{
			name: "good websocket config",
			config: config.ProviderConfig{
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					MaxBufferSize:       1,
					ReconnectionTimeout: time.Second,
					Endpoints: []config.Endpoint{
						{
							URL: "wss://test.com",
						},
					},
					Name:                     "test",
					ReadBufferSize:           config.DefaultReadBufferSize,
					WriteBufferSize:          config.DefaultWriteBufferSize,
					HandshakeTimeout:         config.DefaultHandshakeTimeout,
					EnableCompression:        config.DefaultEnableCompression,
					ReadTimeout:              config.DefaultReadTimeout,
					WriteTimeout:             config.DefaultWriteTimeout,
					MaxSubscriptionsPerBatch: config.DefaultMaxSubscriptionsPerBatch,
				},
				Name: "test",
				Type: "price_provider",
			},
			expectedErr: false,
		},
		{
			name: "no name",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					Endpoints:  []config.Endpoint{{URL: "http://test.com"}},
				},
				Type: "price_provider",
			},
			expectedErr: true,
		},
		{
			name: "no API or websocket config",
			config: config.ProviderConfig{
				Name: "test",
			},
			expectedErr: true,
		},
		{
			name: "both API and websocket config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					Endpoints:  []config.Endpoint{{URL: "http://test.com"}},
				},
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					MaxBufferSize:       1,
					ReconnectionTimeout: time.Second,
					Endpoints: []config.Endpoint{
						{
							URL: "wss://test.com",
						},
					},
					Name:                     "test",
					ReadBufferSize:           config.DefaultReadBufferSize,
					WriteBufferSize:          config.DefaultWriteBufferSize,
					HandshakeTimeout:         config.DefaultHandshakeTimeout,
					EnableCompression:        config.DefaultEnableCompression,
					ReadTimeout:              config.DefaultReadTimeout,
					WriteTimeout:             config.DefaultWriteTimeout,
					MaxSubscriptionsPerBatch: config.DefaultMaxSubscriptionsPerBatch,
				},
				Name: "test",
				Type: "price_provider",
			},
			expectedErr: true,
		},
		{
			name: "bad API config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    2 * time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
				},
				Name: "test",
				Type: "price_provider",
			},
			expectedErr: true,
		},
		{
			name: "bad websocket config",
			config: config.ProviderConfig{
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					ReconnectionTimeout: 2 * time.Second,
				},
				Name: "test",
				Type: "price_provider",
			},
			expectedErr: true,
		},
		{
			name: "mismatch names between provider and ws config",
			config: config.ProviderConfig{
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					MaxBufferSize:       1,
					ReconnectionTimeout: time.Second,
					Endpoints: []config.Endpoint{
						{
							URL: "wss://test.com",
						},
					},
					Name: "test",
				},
				Name: "test2",
				Type: "price_provider",
			},
			expectedErr: true,
		},
		{
			name: "mismatch names between provider and api config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					Endpoints:  []config.Endpoint{{URL: "http://test.com"}},
				},
				Name: "test2",
				Type: "price_provider",
			},
			expectedErr: true,
		},
		{
			name: "no type",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:          true,
					Timeout:          time.Second,
					Interval:         time.Second,
					ReconnectTimeout: time.Second,
					MaxQueries:       1,
					Name:             "test",
					Atomic:           true,
					Endpoints:        []config.Endpoint{{URL: "http://test.com"}},
				},
				Name: "test",
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
