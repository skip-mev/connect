package types

import (
	"fmt"
)

// QuotePriceWithNonce is a wrapper around the QuotePrice object which also contains a nonce.
// The nonce is meant to represent the number of times that the QuotePrice has been updated for a given
// CurrencyPair.
type QuotePriceWithNonce struct {
	QuotePrice
	nonce uint64
}

func NewQuotePriceWithNonce(qp QuotePrice, nonce uint64) QuotePriceWithNonce {
	return QuotePriceWithNonce{
		qp,
		nonce,
	}
}

// Nonce returns the nonce for a given QuotePriceWithNonce.
func (q *QuotePriceWithNonce) Nonce() uint64 {
	return q.nonce
}

// ValidateBasic validates that the QuotePrice is valid, i.e. that the price is non-negative.
func (qp *QuotePrice) ValidateBasic() error {
	// Check that the price is non-negative
	if qp.Price.IsNegative() {
		return fmt.Errorf("price cannot be negative: %s", qp.Price)
	}

	return nil
}

// ValidateBasic validates that the QuotePriceWithNonce is valid, i.e that the underlying QuotePrice is valid.
func (q *QuotePriceWithNonce) ValidateBasic() error {
	return q.QuotePrice.ValidateBasic()
}
