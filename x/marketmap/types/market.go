package types

import (
	"fmt"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

const (
	// MaxConversionOperations is the maximum number of conversion operations that can be used
	// to convert a set of prices to a common ticker. This implementation only supports a maximum
	// of 2 conversion operations - either a direct conversion or a conversion using the index price.
	// This is specific to the IndexPriceAggregation.
	MaxConversionOperations = 2

	// IndexPrice is the provider name for the index price. This is specific to the IndexPriceAggregation.
	IndexPrice = "index"
)

// ValidateBasic performs aggregate validation for all fields in the MarketMap. We consider
// the market map to be valid iff:
//
// 1. Each ticker a provider supports is included in the main set of tickers.
// 2. Each ticker is valid.
// 3. Each provider is valid.
func (mm *MarketMap) ValidateBasic() error {
	if len(mm.Tickers) < len(mm.Providers) {
		return fmt.Errorf("each ticker a provider includes must have a corresponding ticker in the main set of tickers")
	}

	seenCPs := make(map[string]struct{})
	for tickerStr, ticker := range mm.Tickers {
		if err := ticker.ValidateBasic(); err != nil {
			return err
		}

		if tickerStr != ticker.String() {
			return fmt.Errorf("ticker string %s does not match ticker %s", tickerStr, ticker.String())
		}

		seenCPs[ticker.String()] = struct{}{}
	}

	// check if all providers refer to tickers
	for tickerStr, providers := range mm.Providers {
		// check if the ticker is supported
		if _, ok := mm.Tickers[tickerStr]; !ok {
			return fmt.Errorf("provider %s refers to an unsupported ticker", tickerStr)
		}

		if err := providers.ValidateBasic(); err != nil {
			return fmt.Errorf("ticker %s has invalid providers: %w", tickerStr, err)
		}
	}

	return nil
}

// checkIfProviderSupportsTicker checks if the provider supports the given ticker.
func checkIfProviderSupportsTicker(
	provider string,
	cp slinkytypes.CurrencyPair,
	marketMap MarketMap,
) error {
	providers, ok := marketMap.Providers[cp.String()]
	if !ok {
		return fmt.Errorf("provider %s included a ticker %s that has no providers supporting it", provider, cp.String())
	}

	for _, p := range providers.Providers {
		if p.Name == provider {
			return nil
		}
	}

	return fmt.Errorf("provider %s does not support ticker: %s", provider, cp.String())
}

// Equal returns true iff the MarketMap is equal to the given MarketMap.
func (mm *MarketMap) Equal(other MarketMap) bool {
	if len(mm.Tickers) != len(other.Tickers) {
		return false
	}

	if len(mm.Providers) != len(other.Providers) {
		return false
	}

	if len(mm.Paths) != len(other.Paths) {
		return false
	}

	if mm.AggregationType != other.AggregationType {
		return false
	}

	for ticker, tickerData := range mm.Tickers {
		if !tickerData.Equal(other.Tickers[ticker]) {
			return false
		}
	}

	for ticker, providerData := range mm.Providers {
		if !providerData.Equal(other.Providers[ticker]) {
			return false
		}
	}

	for ticker, pathData := range mm.Paths {
		if !pathData.Equal(other.Paths[ticker]) {
			return false
		}
	}

	return true
}
