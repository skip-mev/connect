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
}

func ReadProviderConfigFromFile(path string) (*ProviderConfig, error) {
	// Read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Check required fields
	requiredFields := []string{"name", "path", "timeout", "interval"}
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			return nil, fmt.Errorf("required field %s is missing in config", field)
		}
	}

	// Unmarshal config
	var config ProviderConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
