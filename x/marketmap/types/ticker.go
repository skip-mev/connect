package types

import (
	"fmt"

	"github.com/skip-mev/slinky/pkg/json"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

const (
	// DefaultMaxDecimals is the maximum number of decimals allowed for a ticker.
	DefaultMaxDecimals = 36
	// DefaultMinProviderCount is the minimum number of providers required for a
	// ticker to be considered valid.
	DefaultMinProviderCount = 1
	// MaxPathLength is the maximum length of a path for a ticker conversion.
	MaxPathLength = 2
)

// NewTicker returns a new Ticker instance. A Ticker represents a price feed for
// a given asset pair i.e. BTC/USD. The price feed is scaled to a number of decimal
// places and has a minimum number of providers required to consider the ticker valid.
func NewTicker(base, quote string, decimals, minProviderCount uint64) (Ticker, error) {
	t := Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  base,
			Quote: quote,
		},
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
	return t.CurrencyPair.String()
}

// ValidateBasic performs basic validation on the Ticker.
func (t *Ticker) ValidateBasic() error {
	if t.Decimals > DefaultMaxDecimals || t.Decimals == 0 {
		return fmt.Errorf("decimals must be between 1 and %d; got %d", DefaultMaxDecimals, t.Decimals)
	}
	if t.MinProviderCount < DefaultMinProviderCount {
		return fmt.Errorf("min provider count must be at least %d; got %d", DefaultMinProviderCount, t.MinProviderCount)
	}

	if err := t.CurrencyPair.ValidateBasic(); err != nil {
		return err
	}

	return json.IsValid([]byte(t.Metadata_JSON))
}
