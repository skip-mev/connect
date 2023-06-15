package types

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

// Nonce returns the nonce for a given QuotePriceWithNonce
func (q *QuotePriceWithNonce) Nonce() uint64 {
	return q.nonce
}
