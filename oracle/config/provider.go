package config

import (
	"fmt"
)

// ProviderConfig defines a config for a provider. To add a new provider, add the provider
// config to the oracle configuration.
type ProviderConfig struct {
	// Name identifies which provider this config is for.
	Name string `mapstructure:"name" toml:"name"`

	// API is the config for the API based data provider. If the provider does not
	// support API based fetching, this field should be omitted.
	API APIConfig `mapstructure:"api" toml:"api"`

	// WebSocket is the config for the websocket based data provider. If the provider
	// does not support websocket based fetching, this field should be omitted.
	WebSocket WebSocketConfig `mapstructure:"web_socket" toml:"web_socket"`

	// MarketConfig defines the provider's market configurations. In particular, this defines
	// the mappings between on-chain and off-chain currency pairs.
	MarketConfig ProviderMarketConfig `mapstructure:"market_config" toml:"market_config"`
}

func (c *ProviderConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	if c.API.Enabled && c.WebSocket.Enabled {
		return fmt.Errorf("provider cannot be both API and websocket based")
	}

	if err := c.MarketConfig.ValidateBasic(); err != nil {
		return fmt.Errorf("market config is not formatted correctly %w", err)
	}

	if c.Name != c.MarketConfig.Name {
		return fmt.Errorf("name must match market config name")
	}

	if c.API.Enabled {
		return c.API.ValidateBasic()
	}

	if c.WebSocket.Enabled {
		return c.WebSocket.ValidateBasic()
	}

	return fmt.Errorf("provider must be either enable API or websocket based fetching")
}
