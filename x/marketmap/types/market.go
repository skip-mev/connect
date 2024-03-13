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

// String returns the string representation of the market map.
func (mm *MarketMap) String() string {
	return fmt.Sprintf(
		"MarketMap: {Tickers: %v, Providers: %v, Paths: %v}",
		mm.Tickers,
		mm.Providers,
		mm.Paths,
	)
}
