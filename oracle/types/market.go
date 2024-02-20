package types

import (
	"fmt"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type (
	// TickerToProviderConfig is a type alias for the map of tickers to their respective
	// provider configurations.
	TickerToProviderConfig = map[mmtypes.Ticker]mmtypes.ProviderConfig

	// ProviderMarketMap is a type alias for the provider market map. This is a map of
	// tickers to their respective markets. Provides simple bidirectional mapping between
	// on-chain and off-chain markets.
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

	offChainMap := make(map[string]mmtypes.Ticker, len(tickerConfigs))
	for ticker, config := range tickerConfigs {
		offChainMap[config.OffChainTicker] = ticker
	}

	return ProviderMarketMap{
		Name:          name,
		TickerConfigs: tickerConfigs,
		OffChainMap:   offChainMap,
	}, nil
}

// ValidateBasic validates the provider market map.
func (pmm *ProviderMarketMap) ValidateBasic() error {
	for ticker, config := range pmm.TickerConfigs {
		if config.Name != pmm.Name {
			return fmt.Errorf("expected provider config name %s, got %s", pmm.Name, config.Name)
		}

		if err := config.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid config for ticker %s: %w", ticker.String(), err)
		}
	}

	if len(pmm.OffChainMap) != len(pmm.TickerConfigs) {
		return fmt.Errorf("off-chain map length does not match ticker configs length")
	}

	return nil
}
