package config

import (
	"encoding/json"
	"fmt"
	"os"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func ReadMarketConfigFromFile(path string) (mmtypes.AggregateMarketConfig, error) {
	// Initialize the struct to hold the configuration
	var config mmtypes.AggregateMarketConfig

	// Read the entire file at the given path
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the JSON data into the config struct
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("error unmarshalling config JSON: %w", err)
	}

	// Optionally, validate the config if needed
	if err := config.ValidateBasic(); err != nil {
		return config, fmt.Errorf("error validating config: %w", err)
	}

	// Return the populated config struct
	return config, nil
}
