package coingecko

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// NOTE: To determine the list of supported base currencies, please see
// https://api.coingecko.com/api/v3/coins/list. To see the supported quote
// currencies, please see https://api.coingecko.com/api/v3/simple/supported_vs_currencies.
// Not all base currencies are allowed to be used as quote currencies.

// Config is the config struct for the CoinGecko provider.
type Config struct {
	// APIKey is the API key used to make requests to the CoinGecko API.
	APIKey string `json:"apiKey" validate:"required"`

	// SupportedBases maps an oracle base currency to a CoinGecko base currency.
	SupportedBases map[string]string `json:"supportedBases" validate:"required"`

	// SupportedQuotes maps an oracle quote currency to a CoinGecko quote currency.
	SupportedQuotes map[string]string `json:"supportedQuotes" validate:"required"`
}

// NewConfig returns a new config.
func NewConfig(apiKey string) Config {
	return Config{
		APIKey: apiKey,
	}
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

	return nil
}

// Format returns the formatted config. This is done in accordance with the
// the CoinGecko API's requirements.
func (c *Config) Format() {
	for k, v := range c.SupportedBases {
		delete(c.SupportedBases, k)
		c.SupportedBases[strings.ToUpper(k)] = strings.ToLower(v)
	}

	for k, v := range c.SupportedQuotes {
		delete(c.SupportedQuotes, k)
		c.SupportedQuotes[strings.ToUpper(k)] = strings.ToLower(v)
	}
}

// ReadCoinGeckoConfigFromFile reads a config from a file and returns the config.
func ReadCoinGeckoConfigFromFile(path string) (Config, error) {
	var config Config

	// Read in the config file
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Unmarshal the config.
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// Format the config.
	config.Format()

	// Validate the config.
	if err := config.ValidateBasic(); err != nil {
		return config, err
	}

	return config, nil
}
