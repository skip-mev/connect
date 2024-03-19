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

	switch mm.AggregationType {
	case AggregationType_INDEX_PRICE_AGGREGATION:
		return ValidateIndexPriceAggregation(*mm)
	default:
		return nil
	}
}

// String returns the string representation of the market map.
func (mm *MarketMap) String() string {
	return fmt.Sprintf(
		"MarketMap: {Tickers: %v, Providers: %v, Paths: %v, AggregationType: %s}",
		mm.Tickers,
		mm.Providers,
		mm.Paths,
		mm.AggregationType,
	)
}

// ValidateIndexPriceAggregation validates the market map configuration and its expected configuration for
// this aggregator. In particular, this will
//
//  1. Ensure that the market map is valid (ValidateBasic). This ensure's that each of the provider's
//     markets are supported by the market map.
//  2. Ensure that each path has a corresponding ticker.
//  3. Ensure that each path has a valid number of operations.
//  4. Ensure that each operation has a valid ticker and that the provider supports the ticker.
func ValidateIndexPriceAggregation(
	marketMap MarketMap,
) error {
	for ticker, paths := range marketMap.Paths {
		// The ticker must be supported by the market map. Otherwise we do not how to resolve the
		// prices.
		if _, ok := marketMap.Tickers[ticker]; !ok {
			return fmt.Errorf("path includes a ticker that is not supported: %s", ticker)
		}

		for _, path := range paths.Paths {
			operations := path.Operations
			if len(operations) == 0 || len(operations) > MaxConversionOperations {
				return fmt.Errorf(
					"the expected number of operations is between 1 and %d; got %d operations for %s",
					MaxConversionOperations,
					len(operations),
					ticker,
				)
			}

			first := operations[0]
			if _, ok := marketMap.Tickers[first.CurrencyPair.String()]; !ok {
				return fmt.Errorf("operation included a ticker that is not supported: %s", first.CurrencyPair.String())
			}
			if err := checkIfProviderSupportsTicker(first.Provider, first.CurrencyPair, marketMap); err != nil {
				return err
			}

			if len(operations) != 2 {
				continue
			}

			second := operations[1]
			if second.Provider != IndexPrice {
				return fmt.Errorf("expected index price provider for second operation; got %s", second.Provider)
			}
			if _, ok := marketMap.Tickers[second.CurrencyPair.String()]; !ok {
				return fmt.Errorf("index operation included a ticker that is not supported: %s", second.CurrencyPair.String())
			}
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
