package coinbase

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config is the configuration for the Coinbase APIDataHandler.
type Config struct {
	// SymbolMap maps the oracle's equivalent of an asset to the expected coinbase
	// representation of the asset.
	SymbolMap map[string]string `json:"symbolMap" validate:"required"`
}

// NewConfig returns a new config.
func NewConfig(symbolMap map[string]string) Config {
	return Config{
		SymbolMap: symbolMap,
	}
}

// ValidateBasic performs basic validation on the config.
func (c *Config) ValidateBasic() error {
	if len(c.SymbolMap) == 0 {
		return fmt.Errorf("symbol map cannot be empty")
	}

	for k, v := range c.SymbolMap {
		if len(k) == 0 {
			return fmt.Errorf("symbol map key cannot be empty")
		}

		if len(v) == 0 {
			return fmt.Errorf("symbol map value cannot be empty")
		}
	}

	return nil
}

// Format returns the formatted config. This is done in accordance with the
// the Coinbase API's requirements.
func (c *Config) Format() {
	// Capitalize all symbols.
	for k, v := range c.SymbolMap {
		delete(c.SymbolMap, k)
		c.SymbolMap[strings.ToUpper(k)] = strings.ToUpper(v)
	}
}

// ReadCoinbaseConfigFromFile reads a config from a file and returns the config.
func ReadCoinbaseConfigFromFile(path string) (Config, error) {
	var config Config

	// Read in config file.
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read %s: %s", path, err)
	}

	// Unmarshal the config.
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal %s: %s", path, err)
	}

	// Format the config.
	config.Format()

	// Validate the config.
	if err := config.ValidateBasic(); err != nil {
		return config, fmt.Errorf("invalid %s config: %s", Name, err)
	}

	return config, nil
}
