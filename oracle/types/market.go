package types

import (
	"fmt"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
		OffChainMap map[string]mmtypes.Ticker
	}
)

// NewProviderMarketMap returns a new provider market map.
func NewProviderMarketMap(name string, tickerConfigs TickerToProviderConfig) (ProviderMarketMap, error) {
	if len(tickerConfigs) == 0 {
		return ProviderMarketMap{}, fmt.Errorf("ticker configs cannot be empty")
	}

	if len(name) == 0 {
		return ProviderMarketMap{}, fmt.Errorf("provider name cannot be empty")
	}

	offChainMap := make(map[string]mmtypes.Ticker)
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

		offChainMap[config.OffChainTicker] = ticker
	}

	return ProviderMarketMap{
		Name:          name,
		TickerConfigs: tickerConfigs,
		OffChainMap:   offChainMap,
	}, nil
}

// ValidateBasic performs basic validation on the provider market map.
func (pmm ProviderMarketMap) ValidateBasic() error {
	if len(pmm.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	if len(pmm.TickerConfigs) == 0 {
		return fmt.Errorf("ticker configs cannot be empty")
	}

	if len(pmm.OffChainMap) != len(pmm.TickerConfigs) {
		return fmt.Errorf("off-chain map length mismatch")
	}

	for ticker, config := range pmm.TickerConfigs {
		if err := ticker.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid ticker %s: %w", ticker, err)
		}

		if err := config.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid provider config for %s: %w", ticker, err)

		}

		t, ok := pmm.OffChainMap[config.OffChainTicker]
		if !ok {
			return fmt.Errorf("off-chain ticker %s not found in off-chain map", config.OffChainTicker)
		}

		if t != ticker {
			return fmt.Errorf("off-chain ticker %s does not match on-chain ticker %s", config.OffChainTicker, ticker)
		}
	}

	return nil
}
