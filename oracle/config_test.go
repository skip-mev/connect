package oracle_test

import (
	"os"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/skip-mev/slinky/oracle"
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
	cfg, err := oracle.ReadConfigFromFile(file.Name())
	assert.NoError(t, err)

	// create oracle from config
	tmLogger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	_, err = oracle.NewOracleFromConfig(tmLogger, cfg)
	assert.NoError(t, err)
}
