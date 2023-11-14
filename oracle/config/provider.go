package config

import "github.com/spf13/viper"

// ProviderConfig defines a config for a provider. To add a new provider, add the provider
// config to the oracle configuration.
type ProviderConfig struct {
	// Name identifies which provider this config is for.
	Name string `mapstructure:"name" toml:"name"`

	// Path is the path to the json/toml config file for the provider.
	Path string `mapstructure:"path" toml:"path"`
}

// ReadProviderConfigFromFile reads a config from a file and returns the config.
func ReadProviderConfigFromFile(path string) (*ProviderConfig, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// unmarshal config
	var config ProviderConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
