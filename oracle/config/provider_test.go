package config_test

import (
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/stretchr/testify/require"
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
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
				},
				Name: "test",
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
				},
				Name: "test",
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
				},
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
				},
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					MaxBufferSize:       1,
					ReconnectionTimeout: time.Second,
				},
				Name: "test",
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
