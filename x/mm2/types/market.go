package types

import "fmt"

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

		if err := market.ProviderConfigs.ValidateBasic(); err != nil {
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
