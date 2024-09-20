package types

import (
	"fmt"
	"strings"

	"github.com/skip-mev/connect/v2/pkg/json"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

const (
	// DefaultMaxDecimals is the maximum number of decimals allowed for a ticker.
	DefaultMaxDecimals = 36
	// DefaultMinProviderCount is the minimum number of providers required for a
	// ticker to be considered valid.
	DefaultMinProviderCount = 1
	// MaxMetadataJSONFieldLength is the maximum length of the MetadataJSON field (in bytes).
	MaxMetadataJSONFieldLength = 16384
)

// NewTicker returns a new Ticker instance. A Ticker represents a price feed for
// a given asset pair i.e. BTC/USD. The price feed is scaled to a number of decimal
// places and has a minimum number of providers required to consider the ticker valid.
func NewTicker(base, quote string, decimals, minProviderCount uint64, enabled bool) Ticker {
	return Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
			Base:  strings.ToUpper(base),
			Quote: strings.ToUpper(quote),
		},
		Decimals:         decimals,
		MinProviderCount: minProviderCount,
		Enabled:          enabled,
	}
}

// String returns a string representation of the Ticker.
func (t Ticker) String() string {
	return t.CurrencyPair.String()
}

// ValidateBasic performs basic validation on the Ticker.
func (t *Ticker) ValidateBasic() error {
	if t.Decimals > DefaultMaxDecimals || t.Decimals == 0 {
		return fmt.Errorf("decimals must be between 1 and %d; got %d for %s", DefaultMaxDecimals, t.Decimals, t.CurrencyPair.String())
	}
	if t.MinProviderCount < DefaultMinProviderCount {
		return fmt.Errorf("min provider count must be at least %d; got %d for %s", DefaultMinProviderCount, t.MinProviderCount, t.CurrencyPair.String())
	}

	if err := t.CurrencyPair.ValidateBasic(); err != nil {
		return err
	}

	if len(t.Metadata_JSON) > MaxMetadataJSONFieldLength {
		return fmt.Errorf("metadata json field is longer than maximum length of %d", MaxMetadataJSONFieldLength)
	}

	if err := json.IsValid([]byte(t.Metadata_JSON)); err != nil {
		return fmt.Errorf("invalid ticker metadata json: %w", err)
	}

	return nil
}

// Equal returns true iff the Ticker is equal to the given Ticker.
func (t *Ticker) Equal(other Ticker) bool {
	return t.CurrencyPair.Equal(other.CurrencyPair) &&
		t.Decimals == other.Decimals &&
		t.MinProviderCount == other.MinProviderCount &&
		t.Metadata_JSON == other.Metadata_JSON &&
		t.Enabled == other.Enabled
}
