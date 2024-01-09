package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ProviderConfig defines a config for a provider. To add a new provider, add the provider
// config to the oracle configuration.
type ProviderConfig struct {
	// Name identifies which provider this config is for.
	Name string `mapstructure:"name" toml:"name"`

	// Path is the path to the json/toml config file for the provider.
	Path string `mapstructure:"path" toml:"path"`

	// Timeout is the amount of time the provider should wait for a response from
	// its API before timing out.
	Timeout time.Duration `mapstructure:"timeout" toml:"timeout"`

	// Interval is the interval at which the provider should update the prices.
	Interval time.Duration `mapstructure:"interval" toml:"interval"`

	// MaxQueries is the maximum number of queries that the provider will make
	// within the interval. If the provider makes more queries than this, it will
	// stop making queries until the next interval.
	MaxQueries int `mapstructure:"max_queries" toml:"max_queries"`
}

func (c *ProviderConfig) ValidateBasic() error {
	if len(c.Name) == 0 || len(c.Path) == 0 {
		return fmt.Errorf("name & path cannot be empty")
	}

	if c.Interval <= 0 || c.Timeout <= 0 {
		return fmt.Errorf("provider interval and timeout must be strictly positive")
	}

	if c.Interval < c.Timeout {
		return fmt.Errorf("provider timeout must be greater than 0 and less than the interval")
	}

	if c.MaxQueries < 1 {
		return fmt.Errorf("provider max queries must be greater than 0")
	}

	return nil
}

func ReadProviderConfigFromFile(path string) (ProviderConfig, error) {
	// Read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return ProviderConfig{}, err
	}

	// Check required fields
	requiredFields := []string{"name", "path", "timeout", "interval", "max_queries"}
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			return ProviderConfig{}, fmt.Errorf("required field %s is missing in config", field)
		}
	}

	// Unmarshal config
	var config ProviderConfig
	if err := viper.Unmarshal(&config); err != nil {
		return ProviderConfig{}, err
	}

	if err := config.ValidateBasic(); err != nil {
		return ProviderConfig{}, err
	}

	return config, nil
}
