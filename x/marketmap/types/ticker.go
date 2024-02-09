package types

import (
	"fmt"
	"strings"
)

const (
	// DefaultMaxDecimals is the maximum number of decimals allowed for a ticker.
	DefaultMaxDecimals = 36
	// DefaultMinProviderCount is the minimum number of providers required for a
	// ticker to be considered valid.
	DefaultMinProviderCount = 1
)

// NewTicker returns a new Ticker instance. A Ticker represents a price feed for
// a given asset pair i.e. BTC/USD. The price feed is scaled to a number of decimal
// places and has a minimum number of providers required to consider the ticker valid.
func NewTicker(base, quote string, decimals, minProviderCount uint64) (Ticker, error) {
	t := Ticker{
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

// String returns a string representation of the Ticker.
func (t *Ticker) String() string {
	return fmt.Sprintf("%s/%s", t.Base, t.Quote)
}

// ValidateBasic performs basic validation on the Ticker.
func (t *Ticker) ValidateBasic() error {
	if len(t.Base) == 0 {
		return fmt.Errorf("base cannot be empty")
	}
	if len(t.Quote) == 0 {
		return fmt.Errorf("quote cannot be empty")
	}

	// Ensure the base and quote are upper case.
	if strings.ToUpper(t.Base) != t.Base {
		return fmt.Errorf("base must be upper case; got %s", t.Base)
	}
	if strings.ToUpper(t.Quote) != t.Quote {
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

// NewTickerConfig returns a new TickerConfig instance. The TickerConfig is
// the config the provider uses to create mappings between on-chain and off-chain
// price feeds. The ticker is considered the canonical representation of the price
// feed and the off-chain ticker is the provider specific representation.
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
