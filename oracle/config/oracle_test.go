package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
)

func TestOracleConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.OracleConfig
		expectedErr bool
	}{
		{
			name: "good config",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]config.ProviderConfig{
					"test": {
						Name: "test",
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
						Type: "price_provider",
					},
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: false,
		},
		{
			name: "bad config w/ no max-price-age",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				Providers: map[string]config.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []config.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    config.DefaultReadBufferSize,
							WriteBufferSize:   config.DefaultWriteBufferSize,
							HandshakeTimeout:  config.DefaultHandshakeTimeout,
							EnableCompression: config.DefaultEnableCompression,
							ReadTimeout:       config.DefaultReadTimeout,
							WriteTimeout:      config.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name:        "bad config with no update interval",
			config:      config.OracleConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with bad metrics",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]config.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []config.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    config.DefaultReadBufferSize,
							WriteBufferSize:   config.DefaultWriteBufferSize,
							HandshakeTimeout:  config.DefaultHandshakeTimeout,
							EnableCompression: config.DefaultEnableCompression,
							ReadTimeout:       config.DefaultReadTimeout,
							WriteTimeout:      config.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Metrics: config.MetricsConfig{
					Enabled: true,
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name: "bad config with missing host",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]config.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []config.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    config.DefaultReadBufferSize,
							WriteBufferSize:   config.DefaultWriteBufferSize,
							HandshakeTimeout:  config.DefaultHandshakeTimeout,
							EnableCompression: config.DefaultEnableCompression,
							ReadTimeout:       config.DefaultReadTimeout,
							WriteTimeout:      config.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name: "bad config with missing port",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]config.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []config.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    config.DefaultReadBufferSize,
							WriteBufferSize:   config.DefaultWriteBufferSize,
							HandshakeTimeout:  config.DefaultHandshakeTimeout,
							EnableCompression: config.DefaultEnableCompression,
							ReadTimeout:       config.DefaultReadTimeout,
							WriteTimeout:      config.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Host: "localhost",
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
