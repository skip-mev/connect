package kraken

import (
	"fmt"

	"github.com/spf13/viper"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// Config is the configuration for the Crypto.com provider. To access the available
// markets, please check the following link:
// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#reference-and-market-data-api
type Config struct {
	// Markets is a map of currency pair to perpetual market ID. For an example of the
	// how to configure the markets, please check the readme.
	Markets map[string]string `json:"markets" validate:"required"`

	// Production is true if the provider is running in production mode. Note that if the
	// production setting is set to false, all prices returned by any subscribed markets
	// will be static.
	Production bool `json:"production" validate:"required"`

	// Cache is a map of currency pair to market ID. This is used to cache the market ID
	// for each currency pair.
	Cache map[oracletypes.CurrencyPair]string

	// ReverseCache is a map of perpetual market ID to currency pair.
	ReverseCache map[string]oracletypes.CurrencyPair
}

// ValidateBasic performs basic validation on the config.
func (c *Config) ValidateBasic() error {
	if len(c.Markets) == 0 {
		return fmt.Errorf("no markets specified")
	}

	if len(c.Markets) != len(c.Cache) || len(c.Cache) != len(c.ReverseCache) {
		return fmt.Errorf("reverse markets map is not the same size as the markets map")
	}

	seenMarkets := make(map[string]bool)
	for cp, market := range c.Cache {
		if err := cp.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid currency pair %s", cp)
		}

		if len(market) == 0 {
			return fmt.Errorf("market value cannot be empty")
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

// Format formats the configuration as required by the provider. This will
// also create a reverse map of the markets for convenience.
func (c *Config) Format() error {
	cache := make(map[oracletypes.CurrencyPair]string)
	reverseCache := make(map[string]oracletypes.CurrencyPair)

	for cp, market := range c.Markets {
		cp, err := oracletypes.CurrencyPairFromString(cp)
		if err != nil {
			return err
		}

		cache[cp] = market
		reverseCache[market] = cp
	}

	c.Cache = cache
	c.ReverseCache = reverseCache

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
