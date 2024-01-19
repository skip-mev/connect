package okx

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type Config struct {
	// Markets is the list of markets to subscribe to. The key is the currency pair and the value
	// is the instrument ID. The instrument ID must correspond to the spot market. For example,
	// the instrument ID for the BITCOIN/USDT market is BTC-USDT.
	Markets map[string]string `json:"markets"`

	// Production is true if the config is for production.
	Production bool `json:"production"`

	// Cache is a cache of currency pair to corresponding instrument ID.
	Cache map[oracletypes.CurrencyPair]string

	// ReverseCache is a cache of instrument ID to corresponding currency pair.
	ReverseCache map[string]oracletypes.CurrencyPair
}

// ValidateBasic performs basic validation on the config.
func (c *Config) ValidateBasic() error {
	if len(c.Markets) == 0 {
		return fmt.Errorf("no markets specified")
	}

	if len(c.Markets) != len(c.Cache) || len(c.Markets) != len(c.ReverseCache) {
		return fmt.Errorf("cache does not match markets size")
	}

	seenMarkets := make(map[string]bool)
	for cp, market := range c.Cache {
		if err := cp.ValidateBasic(); err != nil {
			return fmt.Errorf("cache contains invalid currency pair %s: %s", cp.String(), err)
		}

		if len(market) == 0 {
			return fmt.Errorf("cache contains empty instrument ID")
		}

		if _, ok := c.ReverseCache[market]; !ok {
			return fmt.Errorf("failed to find currency pair for market %s in reverse cache", market)
		}

		if _, ok := seenMarkets[market]; ok {
			return fmt.Errorf("duplicate market %s", market)
		}

		seenMarkets[market] = true
	}

	return nil
}

// Format formats the config according to the requirements of the OKX web socket API.
func (c *Config) Format() error {
	c.Cache = make(map[oracletypes.CurrencyPair]string)
	c.ReverseCache = make(map[string]oracletypes.CurrencyPair)

	for cp, instID := range c.Markets {
		delete(c.Markets, cp)

		// OKX expects all of the instrument IDs to be uppercase.
		key := strings.ToUpper(cp)
		value := strings.ToUpper(instID)
		oracleCP, err := oracletypes.CurrencyPairFromString(key)
		if err != nil {
			return fmt.Errorf("invalid currency pair %s: %s", key, err)
		}

		// Update the config.
		c.Markets[key] = value
		c.Cache[oracleCP] = value
		c.ReverseCache[value] = oracleCP
	}

	return nil
}

// ReadConfigFromFile reads the config from the given file path.
func ReadConfigFromFile(path string) (Config, error) {
	var config Config

	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	// Read in the config file.
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Unmarshal the config.
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// Format the config.
	if err := config.Format(); err != nil {
		return config, err
	}

	// Validate the config.
	if err := config.ValidateBasic(); err != nil {
		return config, err
	}

	return config, nil
}
