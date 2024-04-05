package types

import (
	"fmt"

	"github.com/skip-mev/slinky/pkg/json"
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

// SimpleProviderTicker is a simple representation of a provider ticker. Provider's
// can utilize this to easily configure custom tickers.
type SimpleProviderTicker struct {
	Name           string
	OffChainTicker string
	JSON           string
}

// ValidateBasic validates the simple provider ticker.
func (spt *SimpleProviderTicker) ValidateBasic() error {
	if spt.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if spt.OffChainTicker == "" {
		return fmt.Errorf("off-chain ticker cannot be empty")
	}
	return json.IsValid([]byte(spt.JSON))
}

// TickersToProviderTickers is a map of tickers to provider tickers. This should be
// utilized by providers to configure the tickers they will be providing data for.
type TickersToProviderTickers map[mmtypes.Ticker]SimpleProviderTicker

// MustToProviderTickers converts the map to a list of provider tickers.
func (tpt *TickersToProviderTickers) MustToProviderTickers() []ProviderTicker {
	var providerTickers []ProviderTicker

	for _, simpleProviderTicker := range *tpt {
		if err := simpleProviderTicker.ValidateBasic(); err != nil {
			panic(err)
		}

		providerTicker := NewProviderTicker(
			simpleProviderTicker.Name,
			simpleProviderTicker.OffChainTicker,
			simpleProviderTicker.JSON,
			DefaultTickerDecimals,
		)
		providerTickers = append(providerTickers, providerTicker)
	}

	return providerTickers
}

// MustToProviderTickersMap converts the ticker -> provider ticker map to a map of
// ticker.String() -> provider ticker. This is mostly used for testing purposes.
func (tpt *TickersToProviderTickers) MustToProviderTickersMap() map[string]ProviderTicker {
	providerTickers := make(map[string]ProviderTicker)

	for ticker, simpleProviderTicker := range *tpt {
		if err := simpleProviderTicker.ValidateBasic(); err != nil {
			panic(err)
		}

		providerTicker := NewProviderTicker(
			simpleProviderTicker.Name,
			simpleProviderTicker.OffChainTicker,
			simpleProviderTicker.JSON,
			DefaultTickerDecimals,
		)
		providerTickers[ticker.String()] = providerTicker
	}

	return providerTickers
}

// GetProviderTicker returns the provider ticker for the given ticker. This function
// is mostly used for testing purposes.
func (tpt *TickersToProviderTickers) MustGetProviderTicker(ticker mmtypes.Ticker) ProviderTicker {
	simpleProviderTicker, ok := (*tpt)[ticker]
	if !ok {
		panic(fmt.Sprintf("ticker %s not found", ticker))
	}

	if err := simpleProviderTicker.ValidateBasic(); err != nil {
		panic(err)
	}

	providerTicker := NewProviderTicker(
		simpleProviderTicker.Name,
		simpleProviderTicker.OffChainTicker,
		simpleProviderTicker.JSON,
		DefaultTickerDecimals,
	)
	return providerTicker
}
