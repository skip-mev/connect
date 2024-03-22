package types

import (
	"fmt"
)

// ValidateBasic performs aggregate validation for all fields in the MarketMap. We consider
// the market map to be valid iff:
//
// 1. Each ticker a provider supports is included in the main set of tickers.
// 2. Each ticker is valid.
// 3. Each provider is valid.
// 4. Aggregation function is valid.
func (mm *MarketMap) ValidateBasic() error {
	for _, market := range mm.Markets {
		if err := market.ValidateBasic(); err != nil {
			return err
		}
	}

	return ValidateIndexPriceAggregation(*mm)
}

// String returns the string representation of the market map.
func (mm *MarketMap) String() string {
	return fmt.Sprintf(
		"MarketMap: {Markets %v}",
		mm.Markets,
	)
}

// ValidateBasic performs stateless validation of a Market.
func (m *Market) ValidateBasic() error {
	if err := m.Ticker.ValidateBasic(); err != nil {
		return err
	}

	if uint64(len(m.ProviderConfigs)) < m.Ticker.MinProviderCount {
		return fmt.Errorf("this ticker must have at least %d providers; got %d",
			m.Ticker.MinProviderCount,
			len(m.ProviderConfigs),
		)
	}

	seenProviders := make(map[string]struct{})
	for _, providerConfig := range m.ProviderConfigs {
		if err := providerConfig.ValidateBasic(); err != nil {
			return err
		}

		// check for duplicate providers
		if _, seen := seenProviders[providerConfig.Name]; seen {
			return fmt.Errorf("duplicate provider found: %s", providerConfig.Name)
		}
		seenProviders[providerConfig.Name] = struct{}{}

	}

	return nil
}

// String returns the string representation of the market.
func (m *Market) String() string {
	return fmt.Sprintf(
		"Market: {Ticker %v Providers: %v}", m.Ticker, m.ProviderConfigs,
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
func (mm *MarketMap) ValidateIndexPriceAggregation() error {
	for _, market := range mm.Markets {
		for _, providerConfig := range market.ProviderConfigs {
			if providerConfig.NormalizeByPair != nil {
				if _, found := mm.Markets[providerConfig.NormalizeByPair.String()]; !found {
					return fmt.Errorf("provider index of %s was not found in the marketmap", providerConfig.NormalizeByPair.String())
				}
			}
		}
	}

	return nil
}

// Equal returns true if the MarketMap is equal to the given MarketMap.
func (mm *MarketMap) Equal(other MarketMap) bool {
	if len(mm.Markets) != len(other.Markets) {
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

	if len(m.ProviderConfigs) != len(other.ProviderConfigs) {
		return false
	}

	for i, providerConfig := range m.ProviderConfigs {
		if !providerConfig.Equal(other.ProviderConfigs[i]) {
			return false
		}
	}

	return true
}
