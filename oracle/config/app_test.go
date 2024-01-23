package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
)

var (
	goodConfig = `
enabled = true
oracle_address = "localhost:8080"
client_timeout = "10s"
`

	missingAddressConfig = `
enabled = true
client_timeout = "10s"
`

	missingTimeoutConfig = `
enabled = true
oracle_address = "localhost:8080"
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
			},
			expectedErr: false,
		},
		{
			name: "good config with metrics",
			config: config.AppConfig{
				Enabled:                 true,
				OracleAddress:           "localhost:8080",
				ClientTimeout:           time.Second,
				MetricsEnabled:          true,
				PrometheusServerAddress: "localhost:9090",
			},
			expectedErr: false,
		},
		{
			name: "good config with metrics and validator consensus address",
			config: config.AppConfig{
				Enabled:                 true,
				OracleAddress:           "localhost:8080",
				ClientTimeout:           time.Second,
				MetricsEnabled:          true,
				PrometheusServerAddress: "localhost:9090",
				ValidatorConsAddress:    "cosmosvalcons1d3hkxctvdphhxap6xgmrvdfhhg8024",
			},
			expectedErr: false,
		},
		{
			name: "bad config with no oracle address",
			config: config.AppConfig{
				Enabled:       true,
				ClientTimeout: time.Second,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no client timeout",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
			},
			expectedErr: true,
		},
		{
			name: "bad config with no prometheus server address",
			config: config.AppConfig{
				Enabled:                 true,
				OracleAddress:           "localhost:8080",
				ClientTimeout:           time.Second,
				MetricsEnabled:          true,
				PrometheusServerAddress: "",
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad validator consensus address",
			config: config.AppConfig{
				Enabled:                 true,
				OracleAddress:           "localhost:8080",
				ClientTimeout:           time.Second,
				MetricsEnabled:          true,
				PrometheusServerAddress: "localhost:9090",
				ValidatorConsAddress:    "absolutely 0 rizz validator address",
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
			name:        "missing field config",
			config:      missingTimeoutConfig,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file
			f, err := os.CreateTemp("", "oracle_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.config)
			assert.NoError(t, err)

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
