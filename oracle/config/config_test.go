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
[[providers]]
name = "erc4626"
token_name_to_metadata = { "ryETH" = { symbol = "0xb5b29320d2Dde5BA5BAFA1EbcD270052070483ec", decimals = 18 } }
[[providers]]
name = "erc4626-share-price-oracle"
token_name_to_metadata = { "ryBTC" = { symbol = "0x0274a704a6D9129F90A62dDC6f6024b33EcDad36", decimals = 18 } }
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
