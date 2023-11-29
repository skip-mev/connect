package aggregator

import (
	"fmt"
	"math/big"
)

// QuotePrice defines price information for a given CurrencyPair provided
// by a price provider.
type QuotePrice struct {
	// Price tracks the quote price for a given CurrencyPair.
	Price *big.Int
}

// NewQuotePrice returns a new QuotePrice.
func NewQuotePrice(price *big.Int) (QuotePrice, error) {
	if price == nil {
		return QuotePrice{}, fmt.Errorf("last price cannot be nil")
	}

	return QuotePrice{Price: price}, nil
}
