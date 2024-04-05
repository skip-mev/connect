package types

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadMarketConfigFromFile reads a market map configuration from a file at the given path.
func ReadMarketConfigFromFile(path string) (MarketMap, error) {
	// Initialize the struct to hold the configuration
	var config MarketMap

	// Read the entire file at the given path
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the JSON data into the config struct
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("error unmarshalling config JSON: %w", err)
	}

	if err := config.ValidateBasic(); err != nil {
		return config, fmt.Errorf("error validating config: %w", err)
	}

	return config, nil
}
