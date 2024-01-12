package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ProviderConfig defines a config for a provider. To add a new provider, add the provider
// config to the oracle configuration.
type ProviderConfig struct {
	// Name identifies which provider this config is for.
	Name string `mapstructure:"name" toml:"name"`

	// Path is the path to the json/toml config file for the provider.
	Path string `mapstructure:"path" toml:"path"`

	// API is the config for the API based data provider. If the provider does not
	// support API based fetching, this field should be omitted.
	API APIConfig `mapstructure:"api" toml:"api"`

	// WebSocket is the config for the websocket based data provider. If the provider
	// does not support websocket based fetching, this field should be omitted.
	WebSocket WebSocketConfig `mapstructure:"web_socket" toml:"web_socket"`
}

func (c *ProviderConfig) ValidateBasic() error {
	if len(c.Name) == 0 || len(c.Path) == 0 {
		return fmt.Errorf("name & path cannot be empty")
	}

	if c.API.Enabled && c.WebSocket.Enabled {
		return fmt.Errorf("provider cannot be both API and websocket based")
	}

	if c.API.Enabled {
		return c.API.ValidateBasic()
	}

	if c.WebSocket.Enabled {
		return c.WebSocket.ValidateBasic()
	}

	return fmt.Errorf("provider must be either enable API or websocket based fetching")
}

func ReadProviderConfigFromFile(path string) (ProviderConfig, error) {
	// Read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return ProviderConfig{}, err
	}

	// Check required fields
	requiredFields := []string{"name", "path"}
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
