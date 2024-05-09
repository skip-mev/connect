package types

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadMarketMapFromFile reads a market map configuration from a file at the given path.
func ReadMarketMapFromFile(path string) (MarketMap, error) {
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

// Markets is a type alias for an array of unique markets.
type Markets []Market

// ValidateBasic validates each market and checks that they are mutually exclusive.
func (ms Markets) ValidateBasic() error {
	mm := MarketMap{
		Markets: make(map[string]Market, len(ms)),
	}

	for _, m := range ms {
		if _, found := mm.Markets[m.Ticker.String()]; found {
			return fmt.Errorf("found duplicate market: %s", m.Ticker.String())
		}

		mm.Markets[m.Ticker.String()] = m
	}

	return mm.ValidateBasic()
}

// ReadMarketsFromFile reads a market map configuration from a file at the given path.
func ReadMarketsFromFile(path string) (Markets, error) {
	// Initialize the struct to hold the configuration
	var markets Markets

	// Read the entire file at the given path
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the JSON data into the config struct
	if err := json.Unmarshal(data, &markets); err != nil {
		return nil, fmt.Errorf("error unmarshalling config JSON: %w", err)
	}

	return markets, markets.ValidateBasic()
}

func (ms Markets) ToMarketMap() MarketMap {
	mm := MarketMap{
		Markets: make(map[string]Market, len(ms)),
	}

	for _, m := range ms {
		mm.Markets[m.Ticker.String()] = m
	}

	return mm
}
