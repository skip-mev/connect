package config_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/skip-mev/slinky/oracle/config"
)

var (
	goodInProcessOracleConfigContent = `
enabled = true
in_process = true
remote_address = ""
client_timeout = "5s"
update_interval = "10s"
production = true

[[providers]]
name = "testname"
path = "testpath"
api.enabled = true
api.timeout = "5s"
api.interval = "10s"
api.max_queries = 10

[[currency_pairs]]
base = "TESTBASE"
quote = "TESTQUOTE"
`

	badOutOfProcessOracleConfigContent = `
enabled = true
in_process = false
remote_address = ""
client_timeout = "5s"
update_interval = "10s"
production = true

[[providers]]
name = "testname"
path = "testpath"
api.enabled = true
api.timeout = "5s"
api.interval = "10s"
api.max_queries = 10

[[currency_pairs]]
base = "TESTBASE"
quote = "TESTQUOTE"
`

	goodOutOfProcessOracleConfigContent = `
enabled = true
in_process = false
remote_address = "testaddress"
client_timeout = "5s"
update_interval = "10s"
production = true

[[providers]]
name = "testname"
path = "testpath"
api.enabled = true
api.timeout = "5s"
api.interval = "10s"
api.max_queries = 10

[[currency_pairs]]
base = "TESTBASE"
quote = "TESTQUOTE"
`

	badProviderOracleConfigContent = `
enabled = true
in_process = false
remote_address = "testaddress"
client_timeout = "5s"
update_interval = "10s"
production = true

[[providers]]
name = ""
path = ""

[[currency_pairs]]
base = "TESTBASE"
quote = "TESTQUOTE"
`

	badCurrencyPairOracleConfigContent = `
enabled = true
in_process = false
remote_address = "testaddress"
client_timeout = "5s"
update_interval = "10s"
production = true

[[providers]]
name = "testname"
path = "testpath"
api.enabled = true
api.timeout = "5s"
api.interval = "10s"
api.max_queries = 10

[[currency_pairs]]
base = ""
quote = ""
`

	badTimeoutsOracleConfigContent = `
enabled = true
in_process = false
remote_address = "testaddress"
client_timeout = "10"
update_interval = "60"
production = true

[[providers]]
name = "testname"
path = "testpath"
api.enabled = true
api.timeout = "5s"
api.interval = "10s"
api.max_queries = 10

[[currency_pairs]]
base = "TESTBASE"
quote = "TESTQUOTE"
`

	missingFieldOracleConfigContent = `
enabled = true
in_process = false
remote_address = "testaddress"
client_timeout = "5s"
`
)

func TestOracleConfigFromFile(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		expectedErr bool
	}{
		{
			name:        "good in process config",
			config:      goodInProcessOracleConfigContent,
			expectedErr: false,
		},
		{
			name:        "bad out of process config",
			config:      badOutOfProcessOracleConfigContent,
			expectedErr: true,
		},
		{
			name:        "good out of process config",
			config:      goodOutOfProcessOracleConfigContent,
			expectedErr: false,
		},
		{
			name:        "bad provider config",
			config:      badProviderOracleConfigContent,
			expectedErr: true,
		},
		{
			name:        "bad currency pair config",
			config:      badCurrencyPairOracleConfigContent,
			expectedErr: true,
		},
		{
			name:        "bad timeouts config",
			config:      badTimeoutsOracleConfigContent,
			expectedErr: true,
		},
		{
			name:        "missing field config",
			config:      missingFieldOracleConfigContent,
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
			_, err = config.ReadOracleConfigFromFile(f.Name())
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
