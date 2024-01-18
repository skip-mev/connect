package bybit

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type Config struct {
	// SupportedBases maps an oracle base currency to a ByBit base currency.
	SupportedBases map[string]string `json:"supportedBases" validate:"required"`

	// SupportedQuotes maps an oracle quote currency to a ByBit quote currency.
	SupportedQuotes map[string]string `json:"supportedQuotes" validate:"required"`

	// Production is true if the config is for production.
	Production bool `json:"production"`

	// Cache is a cache of currency pair to corresponding bybit pair ID.
	Cache map[oracletypes.CurrencyPair]string

	// ReverseCache is a cache of bybit pair ID to corresponding currency pair.
	ReverseCache map[string]oracletypes.CurrencyPair
}

// ValidateBasic performs basic validation on the config.
func (c *Config) ValidateBasic() error {
	if len(c.SupportedBases) == 0 {
		return fmt.Errorf("must supply at least one supported base currency")
	}

	if len(c.SupportedQuotes) == 0 {
		return fmt.Errorf("must supply at least one supported quote currency")
	}

	for k, v := range c.SupportedBases {
		if len(k) == 0 {
			return fmt.Errorf("supported base currency key cannot be empty")
		}

		if len(v) == 0 {
			return fmt.Errorf("supported base currency value cannot be empty")
		}
	}

	for k, v := range c.SupportedQuotes {
		if len(k) == 0 {
			return fmt.Errorf("supported quote currency key cannot be empty")
		}

		if len(v) == 0 {
			return fmt.Errorf("supported quote currency value cannot be empty")
		}
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

// Format formats the config according to the requirements of the ByBit web socket API.
func (c *Config) Format() error {
	c.Cache = make(map[oracletypes.CurrencyPair]string)
	c.ReverseCache = make(map[string]oracletypes.CurrencyPair)

	for cp, base := range c.SupportedBases {
		delete(c.SupportedBases, cp)
		c.SupportedBases[strings.ToUpper(cp)] = strings.ToUpper(base)
	}

	for cp, quote := range c.SupportedQuotes {
		delete(c.SupportedQuotes, cp)
		c.SupportedQuotes[strings.ToUpper(cp)] = strings.ToUpper(quote)
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
