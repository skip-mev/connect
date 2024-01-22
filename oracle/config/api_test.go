package config_test

import (
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/stretchr/testify/require"
)

func TestAPIConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.APIConfig
		expectedErr bool
	}{
		{
			name: "good config with api enabled",
			config: config.APIConfig{
				Enabled:    true,
				Timeout:    time.Second,
				Interval:   time.Second,
				MaxQueries: 1,
				Name:       "test",
				URL:        "http://test.com",
			},
			expectedErr: false,
		},
		{
			name: "good config with api disabled",
			config: config.APIConfig{
				Enabled: false,
			},
			expectedErr: false,
		},
		{
			name: "bad config with no timeout",
			config: config.APIConfig{
				Enabled:    true,
				Interval:   time.Second,
				MaxQueries: 1,
				Name:       "test",
				URL:        "http://test.com",
			},
			expectedErr: true,
		},
		{
			name: "bad config with no interval",
			config: config.APIConfig{
				Enabled:    true,
				Timeout:    time.Second,
				MaxQueries: 1,
				Name:       "test",
				URL:        "http://test.com",
			},
			expectedErr: true,
		},
		{
			name: "bad config with no max queries",
			config: config.APIConfig{
				Enabled:  true,
				Timeout:  time.Second,
				Interval: time.Second,
				Name:     "test",
				URL:      "http://test.com",
			},
			expectedErr: true,
		},
		{
			name: "bad config with timeout greater than interval",
			config: config.APIConfig{
				Enabled:    true,
				Timeout:    2 * time.Second,
				Interval:   time.Second,
				MaxQueries: 1,
				Name:       "test",
				URL:        "http://test.com",
			},
			expectedErr: true,
		},
		{
			name: "bad config with no name",
			config: config.APIConfig{
				Enabled:    true,
				Timeout:    time.Second,
				Interval:   time.Second,
				MaxQueries: 1,
				URL:        "http://test.com",
			},
			expectedErr: true,
		},
		{
			name: "bad config with no url",
			config: config.APIConfig{
				Enabled:    true,
				Timeout:    time.Second,
				Interval:   time.Second,
				MaxQueries: 1,
				Name:       "test",
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
