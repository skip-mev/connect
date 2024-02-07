package types

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultDecimals  = 8
	EthereumDecimals = 18
)

// NewCurrencyPair returns a new CurrencyPair with the given base and quote strings.
func NewCurrencyPair(base, quote string, decimals int64) CurrencyPair {
	return CurrencyPair{
		Base:     base,
		Quote:    quote,
		Decimals: decimals,
	}
}

// ValidateBasic checks that the Base / Quote strings in the CurrencyPair are formatted correctly, i.e
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

	if cp.Decimals <= 0 {
		return fmt.Errorf("decimals must be greater than 0")
	}

	return nil
}

// String returns a string representation of the CurrencyPair, in the following form "ETH/BTC" where 8 is the number of decimals.
func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s/%d", cp.Base, cp.Quote, cp.Decimals)
}

// Ticker returns a string representation of the CurrencyPair, in the following form "ETH/BTC" where 8 is the number of decimals.
// This function differs from String() which also includes the decimal information.
func (cp CurrencyPair) Ticker() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
}

// CurrencyPairStringToTicker takes a CurrencyPair string and returns the Ticker representation of it.
func CurrencyPairStringToTicker(cpStr string) (string, error) {
	cp, err := CurrencyPairFromString(cpStr)
	if err != nil {
		split := strings.Split(cpStr, "/")
		if len(split) == 2 {
			return cpStr, nil
		}

		return "", fmt.Errorf("invalid string provided: %w", err)
	}

	return cp.Ticker(), nil
}

// CurrencyPairString constructs and returns the string representation of a currency pair.
func CurrencyPairString(base, quote string, decimals int64) string {
	cp := NewCurrencyPair(base, quote, decimals)
	return cp.String()
}

func CurrencyPairFromString(s string) (CurrencyPair, error) {
	split := strings.Split(s, "/")
	if len(split) != 3 {
		return CurrencyPair{}, fmt.Errorf("incorrectly formatted CurrencyPair: %s", s)
	}

	decimals, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return CurrencyPair{}, fmt.Errorf("incorrectly formatted CurrencyPair: %s", s)
	}

	cp := CurrencyPair{
		Base:     strings.ToUpper(split[0]),
		Quote:    strings.ToUpper(split[1]),
		Decimals: decimals,
	}

	return cp, cp.ValidateBasic()
}

func CurrencyPairFromTicker(s string, decimals int64) (CurrencyPair, error) {
	split := strings.Split(s, "/")
	if len(split) != 2 {
		return CurrencyPair{}, fmt.Errorf("incorrectly formatted CurrencyPair: %s", s)
	}

	cp := CurrencyPair{
		Base:     strings.ToUpper(split[0]),
		Quote:    strings.ToUpper(split[1]),
		Decimals: decimals,
	}

	return cp, cp.ValidateBasic()
}

// NewCurrencyPairState returns a new CurrencyPairState given an ID, nonce, and QuotePrice.
func NewCurrencyPairState(id, nonce uint64, quotePrice *QuotePrice, decimals int64) CurrencyPairState {
	return CurrencyPairState{
		Id:       id,
		Nonce:    nonce,
		Price:    quotePrice,
		Decimals: decimals,
	}
}

// ValidateBasic checks that the CurrencyPairState is valid, i.e the nonce is zero if the QuotePrice is nil, and non-zero
// otherwise.
func (cps *CurrencyPairState) ValidateBasic() error {
	// check that the nonce is zero if the QuotePrice is nil
	if cps.Price == nil && cps.Nonce != 0 {
		return fmt.Errorf("invalid nonce, no price update but non-zero nonce: %v", cps.Nonce)
	}

	// check that the nonce is non-zero if the QuotePrice is non-nil
	if cps.Price != nil && cps.Nonce == 0 {
		return fmt.Errorf("invalid nonce, price update but zero nonce: %v", cps.Nonce)
	}

	if cps.Decimals <= 0 {
		return fmt.Errorf("decimals must be greater than 0")
	}

	return nil
}
