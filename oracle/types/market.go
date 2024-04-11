package types

import (
	"fmt"

	pkgtypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// ProviderTickersFromMarketMap returns the set of provider tickers a given provider should
// be providing data for based on the market map.
func ProviderTickersFromMarketMap(
	name string,
	marketMap mmtypes.MarketMap,
) ([]ProviderTicker, error) {
	if err := marketMap.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market map: %w", err)
	}

	var (
		// Track all of the tickers that the provider will be providing data for.
		providerTickers = make([]ProviderTicker, 0)
		// Maintain a set of off-chain tickers that have been seen to avoid duplicates.
		// Notably, the side-car provider enforces a uniquness constraint for off-chain tickers.
		seenOffChainTickers = make(map[string]struct{})
	)

	// Iterate through every single market and its provider configurations to find the
	// provider configurations that match the provider name.
	for _, market := range marketMap.Markets {
		for _, cfg := range market.ProviderConfigs {
			if cfg.Name != name {
				continue
			}
			if _, ok := seenOffChainTickers[cfg.OffChainTicker]; ok {
				continue
			}

			providerTicker := NewProviderTicker(
				cfg.OffChainTicker,
				cfg.Metadata_JSON,
			)
			providerTickers = append(providerTickers, providerTicker)
			seenOffChainTickers[cfg.OffChainTicker] = struct{}{}
		}
	}

	return providerTickers, nil
}

// CurrencyPairsToProviderTickers is a map of tickers to provider tickers. This should be
// utilized by providers to configure the tickers they will be providing data for.
type CurrencyPairsToProviderTickers map[pkgtypes.CurrencyPair]DefaultProviderTicker

// ToProviderTickers converts the map to a list of provider tickers.
func (tpt CurrencyPairsToProviderTickers) ToProviderTickers() []ProviderTicker {
	providerTickers := make([]ProviderTicker, len(tpt))

	i := 0
	for _, ticker := range tpt {
		providerTickers[i] = ticker
		i++
	}

	return providerTickers
}

// MustGetProviderTicker returns the provider ticker for the given currency pair.
// This function is mostly used for testing.
func (tpt CurrencyPairsToProviderTickers) MustGetProviderTicker(cp pkgtypes.CurrencyPair) ProviderTicker {
	providerTicker, ok := tpt[cp]
	if !ok {
		panic(fmt.Sprintf("currency pair %s not found", cp.String()))
	}
	return providerTicker
}

// MustGetProviderConfig returns the provider config for the given currency pair.
// This function is mostly used for testing.
func (tpt CurrencyPairsToProviderTickers) MustGetProviderConfig(name string, cp pkgtypes.CurrencyPair) mmtypes.ProviderConfig {
	providerTicker := tpt.MustGetProviderTicker(cp)
	return mmtypes.ProviderConfig{
		Name:           name,
		OffChainTicker: providerTicker.GetOffChainTicker(),
		Metadata_JSON:  providerTicker.GetJSON(),
	}
}
