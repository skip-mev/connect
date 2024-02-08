package types

import (
	"fmt"
	"strings"
)

const (
	// DefaultMaxDecimals is the maximum number of decimals allowed for a ticker.
	DefaultMaxDecimals = 36
	//DefaultMinProviderCount is the minimum number of providers required for a ticker.
	DefaultMinProviderCount = 1
)

// NewTicker returns a new Ticker instance.
func NewTicker(id uint64, base, quote string, decimals, minProviderCount uint64) (Ticker, error) {
	t := Ticker{
		Id:               id,
		Base:             base,
		Quote:            quote,
		Decimals:         decimals,
		MinProviderCount: minProviderCount,
	}

	if err := t.ValidateBasic(); err != nil {
		return Ticker{}, err
	}

	return t, nil
}

// NewTickerConfig returns a new TickerConfig instance.
func NewTickerConfig(ticker Ticker, offChainTicker string) (TickerConfig, error) {
	config := TickerConfig{
		Ticker:         ticker,
		OffChainTicker: offChainTicker,
	}

	if err := config.ValidateBasic(); err != nil {
		return TickerConfig{}, err
	}

	return config, nil
}

// ValidateBasic performs basic validation on the Ticker.
func (t *Ticker) ValidateBasic() error {
	if len(t.Base) == 0 {
		return fmt.Errorf("base cannot be empty")
	}

	if len(t.Quote) == 0 {
		return fmt.Errorf("quote cannot be empty")
	}

	base := strings.ToUpper(t.Base)
	if base != t.Base {
		return fmt.Errorf("base must be upper case; got %s", t.Base)
	}

	quote := strings.ToUpper(t.Quote)
	if quote != t.Quote {
		return fmt.Errorf("quote must be upper case; got %s", t.Quote)
	}

	if t.Decimals > DefaultMaxDecimals || t.Decimals == 0 {
		return fmt.Errorf("decimals must be between 1 and %d; got %d", DefaultMaxDecimals, t.Decimals)
	}

	if t.MinProviderCount < DefaultMinProviderCount {
		return fmt.Errorf("min provider count must be at least %d; got %d", DefaultMinProviderCount, t.MinProviderCount)
	}

	return nil
}

// ValidateBasic performs basic validation on the TickerConfig.
func (tc *TickerConfig) ValidateBasic() error {
	if err := tc.Ticker.ValidateBasic(); err != nil {
		return err
	}

	if len(tc.OffChainTicker) == 0 {
		return fmt.Errorf("off chain ticker cannot be empty")
	}

	return nil
}
