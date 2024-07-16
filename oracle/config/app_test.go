package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
)

var (
	goodConfig = `
enabled = true
oracle_address = "localhost:8080"
client_timeout = "10s"
interval = "10s"
max_age = "20s"
`

	missingAddressConfig = `
enabled = true
client_timeout = "10s"
interval = "10s"
max_age = "10s"
`

	missingTimeoutConfig = `
enabled = true
oracle_address = "localhost:8080"
interval = "10s"
max_age = "10s"
`

	missingMaxAgeConfig = `
enabled = true
oracle_address = "localhost:8080"
client_timeout = "10s"
interval = "10s"
`

	missingIntervalConfig = `
enabled = true
oracle_address = "localhost:8080"
client_timeout = "10s"
max_age = "10s"
`

	invalidMaxAgeConfig = `
enabled = true
oracle_address = "localhost:8080"
client_timeout = "10s"
interval = "10s"
max_age = "lel"
`

	invalidIntervalConfig = `
enabled = true
oracle_address = "localhost:8080"
client_timeout = "10s"
interval = "lel"
max_age = "10s"
`
)

func TestValidateBasic(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.AppConfig
		expectedErr bool
	}{
		{
			name:        "good config with a disabled oracle",
			config:      config.AppConfig{},
			expectedErr: false,
		},
		{
			name: "good config with no metrics",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				MaxAge:        time.Second * 2,
			},
			expectedErr: false,
		},
		{
			name: "good config with metrics",
			config: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8080",
				ClientTimeout:  time.Second,
				MetricsEnabled: true,
				Interval:       time.Second,
				MaxAge:         time.Second * 2,
			},
			expectedErr: false,
		},
		{
			name: "bad config with no oracle address",
			config: config.AppConfig{
				Enabled:       true,
				ClientTimeout: time.Second,
				Interval:      time.Second,
				MaxAge:        time.Second * 2,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no client timeout",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				Interval:      time.Second,
				MaxAge:        time.Second * 2,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no max age",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no interval",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				MaxAge:        time.Second * 2,
			},
			expectedErr: true,
		},
		{
			name: "bad config with max age being less than interval",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				MaxAge:        time.Millisecond,
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

func TestReadConfigFromFile(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		expectedErr bool
	}{
		{
			name:        "good config",
			config:      goodConfig,
			expectedErr: false,
		},
		{
			name:        "bad config",
			config:      missingAddressConfig,
			expectedErr: true,
		},
		{
			name:        "missing timeout field config",
			config:      missingTimeoutConfig,
			expectedErr: true,
		},
		{
			name:        "missing max age field config",
			config:      missingMaxAgeConfig,
			expectedErr: true,
		},
		{
			name:        "missing interval field config",
			config:      missingIntervalConfig,
			expectedErr: true,
		},
		{
			name:        "invalid max age config",
			config:      invalidMaxAgeConfig,
			expectedErr: true,
		},
		{
			name:        "invalid interval config",
			config:      invalidIntervalConfig,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file
			f, err := os.CreateTemp("", "oracle_config")
			require.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.config)
			require.NoError(t, err)

			// Read config from file
			_, err = config.ReadConfigFromFile(f.Name())
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
