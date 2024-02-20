package types

import "fmt"

// ValidateBasic performs aggregate validation for all fields in the MarketMap.
func (mm *MarketMap) ValidateBasic() error {
	if len(mm.Paths) != len(mm.Tickers) || len(mm.Paths) != len(mm.Providers) {
		return fmt.Errorf("length of Tickers, Paths and Providers must be equal")
	}

	for _, ticker := range mm.Tickers {
		if err := ticker.ValidateBasic(); err != nil {
			return err
		}

		paths, ok := mm.Paths[ticker.String()]
		if !ok {
			return fmt.Errorf("paths for ticker %s not found", ticker.String())
		}

		if err := paths.ValidateBasic(ticker.CurrencyPair); err != nil {
			return err
		}

		providers, ok := mm.Providers[ticker.String()]
		if !ok {
			return fmt.Errorf("providers for ticker %s not found", ticker.String())
		}

		if err := providers.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
