package config_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/skip-mev/slinky/oracle/config"
)

var (
	goodConfig = `
oracle_path = "testpath"
metrics_path = "testpath"
`

	badConfig = `
oracle_path = ""
metrics_path = ""
`

	missingFieldConfig = `
oracle_path = ""
`
)

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
			config:      badConfig,
			expectedErr: true,
		},
		{
			name:        "missing field config",
			config:      missingFieldConfig,
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
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
