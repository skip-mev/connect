package types

import (
	"encoding/json"
	"fmt"
	"os"

	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

type (
	// TickerToProviderConfig is a type alias for the map of tickers to their respective
	// provider configurations.
	TickerToProviderConfig = map[mmtypes.Ticker]mmtypes.ProviderConfig

	// ProviderMarketMap provides a simple adapter of what the market map module expects
	// from a provider and how the provider configures its markets. It provides a bi-directional
	// mapping between on-chain and off-chain markets. Every ProviderMarketMap is expected to be
	// constructed from a given market map module configuration.
	ProviderMarketMap struct {
		// Name is the name of the provider.
		Name string

		// TickerConfigs is a map of tickers to their respective off-chain markets as
		// configured by the market map module.
		TickerConfigs TickerToProviderConfig

		// OffChainMap is a map of tickers to their respective on-chain markets.
		OffChainMap map[string][]mmtypes.Ticker
	}
)

// ProviderMarketMapFromMarketMap returns a provider market map from a market map provided by the
// market map module.
func ProviderMarketMapFromMarketMap(name string, marketMap mmtypes.MarketMap) (ProviderMarketMap, error) {
	if err := marketMap.ValidateBasic(); err != nil {
		return ProviderMarketMap{}, fmt.Errorf("invalid market map: %w", err)
	}

	tickers := make(TickerToProviderConfig)

	// Iterate over the providers and their respective tickers.
	for _, market := range marketMap.Markets {
		for _, provider := range market.ProviderConfigs {
			if provider.Name != name {
				continue
			}

			tickers[market.Ticker] = provider
			break
		}
	}

	return NewProviderMarketMap(name, tickers)
}

// NewProviderMarketMap returns a new provider market map.
func NewProviderMarketMap(name string, tickerConfigs TickerToProviderConfig) (ProviderMarketMap, error) {
	if len(name) == 0 {
		return ProviderMarketMap{}, fmt.Errorf("provider name cannot be empty")
	}

	if len(tickerConfigs) == 0 {
		return ProviderMarketMap{
			Name:          name,
			TickerConfigs: make(map[mmtypes.Ticker]mmtypes.ProviderConfig),
			OffChainMap:   make(map[string][]mmtypes.Ticker),
		}, nil
	}

	offChainMap := make(map[string][]mmtypes.Ticker, len(tickerConfigs))
	for ticker, config := range tickerConfigs {
		if err := ticker.ValidateBasic(); err != nil {
			return ProviderMarketMap{}, fmt.Errorf("invalid ticker %s: %w", ticker, err)
		}

		if err := config.ValidateBasic(); err != nil {
			return ProviderMarketMap{}, fmt.Errorf("invalid provider config for %s: %w", ticker, err)
		}

		if config.Name != name {
			return ProviderMarketMap{}, fmt.Errorf("expected provider config name %s, got %s", name, config.Name)
		}

		if offChainMap[config.OffChainTicker] == nil {
			offChainMap[config.OffChainTicker] = []mmtypes.Ticker{ticker}
			continue
		}

		offChainMap[config.OffChainTicker] = append(offChainMap[config.OffChainTicker], ticker)
	}

	pmm := ProviderMarketMap{
		Name:          name,
		TickerConfigs: tickerConfigs,
		OffChainMap:   offChainMap,
	}

	return pmm, pmm.ValidateBasic()
}

// GetTickers returns the tickers from the provider market map.
func (pmm *ProviderMarketMap) GetTickers() []mmtypes.Ticker {
	tickers := make([]mmtypes.Ticker, 0, len(pmm.TickerConfigs))
	for ticker := range pmm.TickerConfigs {
		tickers = append(tickers, ticker)
	}
	return tickers
}

// ValidateBasic performs basic validation on the provider market map.
func (pmm *ProviderMarketMap) ValidateBasic() error {
	if len(pmm.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}
	if len(pmm.OffChainMap) > len(pmm.TickerConfigs) {
		return fmt.Errorf("off-chain map length invalid %d>%d, %v %v", len(pmm.OffChainMap), len(pmm.TickerConfigs), pmm.OffChainMap, pmm.TickerConfigs)
	}

	for ticker, config := range pmm.TickerConfigs {
		if err := ticker.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid ticker %s: %w", ticker, err)
		}

		if err := config.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid provider config for %s: %w", ticker, err)
		}

		tickers, ok := pmm.OffChainMap[config.OffChainTicker]
		if !ok {
			return fmt.Errorf("off-chain ticker %s not found in off-chain map", config.OffChainTicker)
		}

		found := false
		for _, t := range tickers {
			if t == ticker {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("ticker %s not found in off-chain map", config.OffChainTicker)
		}
	}

	return nil
}

// ReadMarketConfigFromFile reads a market map configuration from a file at the given path.
func ReadMarketConfigFromFile(path string) (mmtypes.MarketMap, error) {
	// Initialize the struct to hold the configuration
	var config mmtypes.MarketMap

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
