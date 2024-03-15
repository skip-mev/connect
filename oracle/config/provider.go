package config

import (
	"fmt"
)

// ProviderConfig defines a config for a provider. To add a new provider, add the provider
// config to the oracle configuration.
type ProviderConfig struct {
	// Name identifies which provider this config is for.
	Name string `json:"name"`

	// API is the config for the API based data provider. If the provider does not
	// support API based fetching, this field should be omitted.
	API APIConfig `json:"api"`

	// WebSocket is the config for the websocket based data provider. If the provider
	// does not support websocket based fetching, this field should be omitted.
	WebSocket WebSocketConfig `json:"webSocket"`

	// Type is the type of the provider (i.e. price, market map, other). This is used
	// to determine how to construct the provider.
	Type string `json:"type"`
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

	if len(c.Type) == 0 {
		return fmt.Errorf("type cannot be empty")
	}

	return nil
}
