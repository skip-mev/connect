package config_test

import (
	"os"
	"testing"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/stretchr/testify/assert"
)

func TestOracleFromConfig(t *testing.T) {
	// write base oracle config to file
	cfgStr := `
update_interval = "1s"
[[providers]]
name = "coingecko"
[[providers]]
name = "coinbase"
[[currency_pairs]]
base = "ATOM"
quote = "USD"
quote_decimals = 8
`
	// write config to file
	file, err := os.CreateTemp("", "oracle_config")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(cfgStr)
	assert.NoError(t, err)

	// read config from file
	cfg, err := config.ReadConfigFromFile(file.Name())
	assert.NoError(t, err)

	// create oracle from config
	_, err = oracle.NewOracleFromConfig(log.NewNopLogger(), cfg)
	assert.NoError(t, err)
}
