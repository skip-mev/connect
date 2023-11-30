package types

import (
	"fmt"
	"strings"
)

const (
	ethereum = "ETHEREUM"
)

// NewCurrencyPair returns a new CurrencyPair with the given base and quote strings.
func NewCurrencyPair(base, quote string) CurrencyPair {
	return CurrencyPair{
		Base:  base,
		Quote: quote,
	}
}

// ValidateBasic checks that the Base / Quote strings in the CurrencyPair are formatted correctly, i.e
// Base + Quote are non-empty, and are in upper-case.
func (cp CurrencyPair) ValidateBasic() error {
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
	return nil
}

// ToString returns a string representation of the CurrencyPair, in the following form "ETH/BTC".
//
// NOTICE: prefer ToString over the default String method, as the ToString method is used for marshalling
// currency-pairs into vote-extensions.
func (cp CurrencyPair) ToString() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
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

// Decimals returns the number of decimals that the quote will be reported to. If the quote is Ethereum, then
// the number of decimals is 18. Otherwise, the decimals will be reorted to 8.
func (cp CurrencyPair) Decimals() int {
	if strings.ToUpper(cp.Quote) == ethereum {
		return 18
	}
	return 8
}

// NewCurrencyPairState returns a new CurrencyPairState given an Id, nonce, and QuotePrice.
func NewCurrencyPairState(id uint64, nonce uint64, quotePrice *QuotePrice) CurrencyPairState {
	return CurrencyPairState{
		Id:    id,
		Nonce: nonce,
		Price: quotePrice,
	}
}

// ValidateBasic checks that the CurrencyPairState is valid, i.e the nonce is zero if the QuotePrice is nil, and non-zero
// otherwise.
func (cps CurrencyPairState) ValidateBasic() error {
	// check that the nonce is zero if the QuotePrice is nil
	if cps.Price == nil && cps.Nonce != 0 {
		return fmt.Errorf("invalid nonce, no price update but non-zero nonce: %v", cps.Nonce)
	}

	// check that the nonce is non-zero if the QuotePrice is non-nil
	if cps.Price != nil && cps.Nonce == 0 {
		return fmt.Errorf("invalid nonce, price update but zero nonce: %v", cps.Nonce)
	}

	return nil
}
