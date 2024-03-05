package types

import (
	"fmt"

	"github.com/skip-mev/slinky/pkg/types"
)

// ValidateBasic performs aggregate validation for all fields in the MarketMap. We consider
// the market map to be valid iff:
//
// 1. Each ticker has a corresponding provider supporting it.
// 2. Each ticker is valid.
// 3. Each provider is valid.
// 4. Each path is valid.
// 5. Each operation (ticker) in each path is supported by the market map.
// 6. The enabled list is valid.
func (mm *MarketMap) ValidateBasic() error {
	if len(mm.Tickers) != len(mm.Providers) {
		return fmt.Errorf("each ticker must have a corresponding provider list supporting it")
	}

	if len(mm.EnabledTickers.Tickers) > len(mm.Tickers) {
		return fmt.Errorf("enabled tickers cannot be longer than tickers")
	}

	// verify all enabled tickers are in the Tickers Map
	seenEnabledTickers := make(map[string]struct{})
	for _, tickerStr := range mm.EnabledTickers.Tickers {
		if _, found := mm.Tickers[tickerStr]; !found {
			return fmt.Errorf("ticker ID %s in enabled tickers not found in tickers map", tickerStr)
		}

		if _, seen := seenEnabledTickers[tickerStr]; seen {
			return fmt.Errorf("duplicate ticker ID %s found in enabled list", tickerStr)
		}

		seenEnabledTickers[tickerStr] = struct{}{}
	}

	seenCPs := make(map[string]struct{})
	for tickerStr, ticker := range mm.Tickers {
		if err := ticker.ValidateBasic(); err != nil {
			return err
		}

		if tickerStr != ticker.String() {
			return fmt.Errorf("ticker string %s does not match ticker %s", tickerStr, ticker.String())
		}

		providers, ok := mm.Providers[ticker.String()]
		if !ok {
			return fmt.Errorf("providers for ticker %s not found", ticker.String())
		}

		if err := providers.ValidateBasic(); err != nil {
			return err
		}

		seenCPs[ticker.String()] = struct{}{}
	}

	for tickerStr, paths := range mm.Paths {
		cp, err := types.CurrencyPairFromString(tickerStr)
		if err != nil {
			return err
		}

		if err := paths.ValidateBasic(cp); err != nil {
			return err
		}

		for _, path := range paths.Paths {
			for _, operation := range path.Operations {
				if _, ok := seenCPs[operation.CurrencyPair.String()]; !ok {
					return fmt.Errorf("currency pair %s not found in market map", operation.CurrencyPair)
				}
			}
		}
	}

	// check if all providers refer to tickers
	for tickerStr := range mm.Providers {
		if _, ok := seenCPs[tickerStr]; !ok {
			return fmt.Errorf("currency pair %s not found in market map", tickerStr)
		}
	}

	return nil
}

// String returns the string representation of the market map.
func (mm *MarketMap) String() string {
	return fmt.Sprintf(
		"MarketMap: {Tickers: %v, Providers: %v, Paths: %v, EnabledTickers: %v}",
		mm.Tickers,
		mm.Providers,
		mm.Paths,
		mm.EnabledTickers,
	)
}
