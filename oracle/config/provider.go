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

	// Market defines the provider's market configurations. In particular, this defines
	// the mappings between on-chain and off-chain currency pairs.
	Market MarketConfig `mapstructure:"market_config" toml:"market_config"`
}

func (c *ProviderConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	if c.API.Enabled && c.WebSocket.Enabled {
		return fmt.Errorf("provider %s cannot be both API and websocket based", c.Name)
	}

	if !c.API.Enabled && !c.WebSocket.Enabled {
		return fmt.Errorf("provider %s must be either API or websocket based", c.Name)
	}

	if c.API.Enabled {
		if err := c.API.ValidateBasic(); err != nil {
			return fmt.Errorf("api config for %s is not formatted correctly: %w", c.Name, err)
		}

		if c.API.Name != c.Name {
			return fmt.Errorf("received api config for %s but expected %s", c.API.Name, c.Name)
		}
	}

	if c.WebSocket.Enabled {
		if err := c.WebSocket.ValidateBasic(); err != nil {
			return fmt.Errorf("websocket config for %s is not formatted correctly: %w", c.Name, err)
		}

		if c.WebSocket.Name != c.Name {
			return fmt.Errorf("received websocket config for %s but expected %s", c.WebSocket.Name, c.Name)
		}
	}

	if err := c.Market.ValidateBasic(); err != nil {
		return fmt.Errorf("market config for %s is not formatted correctly: %w", c.Name, err)
	}

	if c.Name != c.Market.Name {
		return fmt.Errorf("name must match market config name; %s != %s", c.Name, c.Market.Name)
	}

	return nil
}
