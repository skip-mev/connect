package config_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/skip-mev/slinky/oracle/config"
)

var (
	goodMetricsConfigContent = `
prometheus_server_address = "testaddress"

[oracle_metrics]
enabled = true

[app_metrics]
enabled = true
validator_cons_address = "cosmosvalcons1weskcvgqsytu2"
`

	badAppMetricsConfigContent = `
prometheus_server_address = "testaddress"

[oracle_metrics]
enabled = true

[app_metrics]
enabled = true
validator_cons_address = "badconsaddress"
`

	badOracleMetricsConfigContent = `
prometheus_server_address = ""

[oracle_metrics]
enabled = true

[app_metrics]
enabled = true
validator_cons_address = "cosmosvalcons1weskcvgqsytu2"
`

	missingFieldMetricsConfigContent = `
prometheus_server_address = ""

[oracle_metrics]
enabled = true
`
)

func TestReadMetricsConfigFromFile(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		expectedErr bool
	}{
		{
			name:        "good config",
			config:      goodMetricsConfigContent,
			expectedErr: false,
		},
		{
			name:        "bad app metrics config",
			config:      badAppMetricsConfigContent,
			expectedErr: true,
		},
		{
			name:        "bad oracle metrics config",
			config:      badOracleMetricsConfigContent,
			expectedErr: true,
		},
		{
			name:        "missing field config",
			config:      missingFieldMetricsConfigContent,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file
			f, err := os.CreateTemp("", "oracle_metrics_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.config)
			assert.NoError(t, err)

			// Read config from file
			_, err = config.ReadMetricsConfigFromFile(f.Name())
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
