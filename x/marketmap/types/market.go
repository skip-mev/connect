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
	IndexPrice = "index_price"
)

// ValidateBasic performs aggregate validation for all fields in the MarketMap. We consider
// the market map to be valid iff:
//
// 1. Each ticker a provider supports is included in the main set of tickers.
// 2. Each ticker is valid.
// 3. Each provider is valid.
func (mm *MarketMap) ValidateBasic() error {
	seenCPs := make(map[string]struct{})
	for tickerStr, market := range mm.Markets {
		if err := market.Ticker.ValidateBasic(); err != nil {
			return err
		}

		if tickerStr != market.Ticker.String() {
			return fmt.Errorf("ticker string %s does not match ticker %s", tickerStr, market.Ticker.String())
		}

		seenCPs[market.Ticker.String()] = struct{}{}

		if err := market.Providers.ValidateBasic(); err != nil {
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
		"MarketMap: {Markets %v AggregationType: %s}",
		mm.Markets,
		mm.AggregationType,
	)
}

// ValidateBasic performs stateless validation of a Market.
func (m *Market) ValidateBasic() error {
	if err := m.Ticker.ValidateBasic(); err != nil {
		return err
	}

	for _, path := range m.Paths.Paths {
		if err := path.ValidateBasic(); err != nil {
			return err
		}
	}

	if uint64(len(m.Providers.Providers)) < m.Ticker.MinProviderCount {
		return fmt.Errorf("this ticker must have at least %d providers; got %d",
			m.Ticker.MinProviderCount,
			len(m.Providers.Providers),
		)
	}

	seenProviders := make(map[string]struct{})
	for _, provider := range m.Providers.Providers {
		if err := provider.ValidateBasic(); err != nil {
			return err
		}

		// check for duplicate providers
		if _, seen := seenProviders[provider.Name]; seen {
			return fmt.Errorf("duplicate provider found: %s", provider.Name)
		}
		seenProviders[provider.Name] = struct{}{}

	}

	return nil
}

// String returns the string representation of the market.
func (m *Market) String() string {
	return fmt.Sprintf(
		"Market: {Ticker %v Paths: %v Providers: %v}", m.Ticker, m.Paths, m.Providers,
	)
}

// ValidateIndexPriceAggregation validates the market map configuration and its expected configuration for
// this aggregator. In particular, this will
//
//  1. Ensure that the market map is valid (ValidateBasic). This ensures that each of the provider's
//     markets are supported by the market map.
//  2. Ensure that each path has a corresponding ticker.
//  3. Ensure that each path has a valid number of operations.
//  4. Ensure that each operation has a valid ticker and that the provider supports the ticker.
func ValidateIndexPriceAggregation(
	marketMap MarketMap,
) error {
	for tickerStr, market := range marketMap.Markets {
		// The ticker must be supported by the market map. Otherwise, we do not how to resolve the
		// prices.
		if _, ok := marketMap.Markets[tickerStr]; !ok {
			return fmt.Errorf("path includes a ticker that is not supported: %s", tickerStr)
		}

		for _, path := range market.Paths.Paths {
			operations := path.Operations
			if len(operations) == 0 || len(operations) > MaxConversionOperations {
				return fmt.Errorf(
					"the expected number of operations is between 1 and %d; got %d operations for %s",
					MaxConversionOperations,
					len(operations),
					tickerStr,
				)
			}

			first := operations[0]
			if _, ok := marketMap.Markets[first.CurrencyPair.String()]; !ok {
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
			if _, ok := marketMap.Markets[second.CurrencyPair.String()]; !ok {
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
	market, ok := marketMap.Markets[cp.String()]
	if !ok {
		return fmt.Errorf("provider %s included a ticker %s that has no providers supporting it", provider, cp.String())
	}

	for _, p := range market.Providers.Providers {
		if p.Name == provider {
			return nil
		}
	}

	return fmt.Errorf("provider %s does not support ticker: %s", provider, cp.String())
}

// Equal returns true if the MarketMap is equal to the given MarketMap.
func (mm *MarketMap) Equal(other MarketMap) bool {
	if len(mm.Markets) != len(other.Markets) {
		return false
	}

	if mm.AggregationType != other.AggregationType {
		return false
	}

	for tickerStr, market := range mm.Markets {
		otherMarket, found := other.Markets[tickerStr]
		if !found {
			return false
		}

		if !market.Equal(otherMarket) {
			return false
		}
	}

	return true
}

// Equal returns true if the Market is equal to the given Market.
func (m *Market) Equal(other Market) bool {
	if !m.Ticker.Equal(other.Ticker) {
		return false
	}

	if !m.Providers.Equal(other.Providers) {
		return false
	}

	if !m.Paths.Equal(other.Paths) {
		return false

	}

	return true
}
