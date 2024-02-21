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
func (mm *MarketMap) ValidateBasic() error {
	if len(mm.Tickers) != len(mm.Providers) {
		return fmt.Errorf("each ticker must have a corresponding provider supporting it")
	}

	cps := make(map[types.CurrencyPair]struct{})
	for _, ticker := range mm.Tickers {
		if err := ticker.ValidateBasic(); err != nil {
			return err
		}

		providers, ok := mm.Providers[ticker.String()]
		if !ok {
			return fmt.Errorf("providers for ticker %s not found", ticker.String())
		}

		if err := providers.ValidateBasic(); err != nil {
			return err
		}

		cps[ticker.CurrencyPair] = struct{}{}
	}

	for ticker, paths := range mm.Paths {
		cp, err := types.CurrencyPairFromString(ticker)
		if err != nil {
			return err
		}

		if err := paths.ValidateBasic(cp); err != nil {
			return err
		}

		for _, path := range paths.Paths {
			for _, operation := range path.Operations {
				if _, ok := cps[operation.CurrencyPair]; !ok {
					return fmt.Errorf("currency pair %s not found in market map", operation.CurrencyPair)
				}
			}
		}
	}

	return nil
}
