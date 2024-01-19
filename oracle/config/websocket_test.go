package config_test

import (
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/stretchr/testify/require"
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
				ReconnectionTimeout: time.Second,
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
			},
			expectedErr: true,
		},
		{
			name: "bad config with no reconnection timeout",
			config: config.WebSocketConfig{
				Enabled:       true,
				MaxBufferSize: 1,
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
