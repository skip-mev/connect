package types

import (
	"fmt"
)

// NewCurrencyPairState returns a new CurrencyPairState given an ID, nonce, and QuotePrice.
func NewCurrencyPairState(id, nonce uint64, quotePrice *QuotePrice) CurrencyPairState {
	return CurrencyPairState{
		Id:    id,
		Nonce: nonce,
		Price: quotePrice,
	}
}

// ValidateBasic checks that the CurrencyPairState is valid, i.e. the nonce is zero if the QuotePrice is nil, and non-zero
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

	return nil
}
