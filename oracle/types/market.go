package types

import (
	"fmt"

	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
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
				cfg.Name,
				cfg.OffChainTicker,
				cfg.Metadata_JSON,
				DefaultTickerDecimals,
			)
			providerTickers = append(providerTickers, providerTicker)
			seenOffChainTickers[cfg.OffChainTicker] = struct{}{}
		}
	}

	return providerTickers, nil
}

// TickersToProviderTickers is a map of tickers to provider tickers. This should be
// utilized by providers to configure the tickers they will be providing data for.
type TickersToProviderTickers map[mmtypes.Ticker]DefaultProviderTicker

// ToProviderTickers converts the map to a list of provider tickers.
func (tpt *TickersToProviderTickers) ToProviderTickers() []ProviderTicker {
	var providerTickers = make([]ProviderTicker, len(*tpt))

	i := 0
	for _, ticker := range *tpt {
		// If the ticker does not have a set number of decimals, set it to the default.
		if ticker.Decimals == 0 {
			ticker.Decimals = DefaultTickerDecimals
		}

		providerTickers[i] = ticker
		i++
	}

	return providerTickers
}

// MustGetProviderTicker returns the provider ticker for the given x/marketmap ticker.
// This function is mostly used for testing.
func (tpt *TickersToProviderTickers) MustGetProviderTicker(ticker mmtypes.Ticker) ProviderTicker {
	providerTicker, ok := (*tpt)[ticker]
	if !ok {
		panic(fmt.Sprintf("ticker %s not found", ticker))
	}

	// If the ticker does not have a set number of decimals, set it to the default.
	if providerTicker.Decimals == 0 {
		providerTicker.Decimals = DefaultTickerDecimals
	}
	return providerTicker
}
