package types

import (
	"fmt"
	"strings"
)

const (
	ethereum         = "ETHEREUM"
	MaxCPFieldLength = 128
)

// NewCurrencyPair returns a new CurrencyPair with the given base and quote strings.
func NewCurrencyPair(base, quote string) CurrencyPair {
	return CurrencyPair{
		Base:  base,
		Quote: quote,
	}
}

// ValidateBasic checks that the Base / Quote strings in the CurrencyPair are formatted correctly, i.e.
// Base + Quote are non-empty, and are in upper-case.
func (cp *CurrencyPair) ValidateBasic() error {
	// strings must be valid
	if cp.Base == "" || cp.Quote == "" {
		return fmt.Errorf("empty quote or base string")
	}
	// check formatting of base / quote
	if strings.ToUpper(cp.Base) != cp.Base {
		return fmt.Errorf("incorrectly formatted base string, expected: %s got: %s", strings.ToUpper(cp.Base), cp.Base)
	}
	if strings.ToUpper(cp.Quote) != cp.Quote {
		return fmt.Errorf("incorrectly formatted quote string, expected: %s got: %s", strings.ToUpper(cp.Quote), cp.Quote)
	}

	if len(cp.Base) > MaxCPFieldLength || len(cp.Quote) > MaxCPFieldLength {
		return fmt.Errorf("string field exceeds maximum length of %d", MaxCPFieldLength)
	}

	return nil
}

// Invert returns an inverted version of cp (where the Base and Quote are swapped).
func (cp *CurrencyPair) Invert() CurrencyPair {
	return CurrencyPair{
		Base:  cp.Quote,
		Quote: cp.Base,
	}
}

// String returns a string representation of the CurrencyPair, in the following form "ETH/BTC".
func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
}

// CurrencyPairString constructs and returns the string representation of a currency pair.
func CurrencyPairString(base, quote string) string {
	cp := NewCurrencyPair(base, quote)
	return cp.String()
}

func CurrencyPairFromString(s string) (CurrencyPair, error) {
	split := strings.Split(s, "/")
	if len(split) != 2 {
		return CurrencyPair{}, fmt.Errorf("incorrectly formatted CurrencyPair: %s", s)
	}
	cp := CurrencyPair{
		Base:  strings.ToUpper(split[0]),
		Quote: strings.ToUpper(split[1]),
	}

	return cp, cp.ValidateBasic()
}

// LegacyDecimals returns the number of decimals that the quote will be reported to. If the quote is Ethereum, then
// the number of decimals is 18. Otherwise, the decimals will be reorted to 8.
func (cp *CurrencyPair) LegacyDecimals() int {
	if strings.ToUpper(cp.Quote) == ethereum {
		return 18
	}
	return 8
}

// Equal returns true iff the CurrencyPair is equal to the given CurrencyPair.
func (cp *CurrencyPair) Equal(other CurrencyPair) bool {
	return cp.Base == other.Base && cp.Quote == other.Quote
}
