package types

import (
	"fmt"
	"time"

	"github.com/holiman/uint256"
)

// QuotePrice defines price information for a given CurrencyPair provided
// by a price provider.
type QuotePrice struct {
	// Price tracks the quote price for a given CurrencyPair.
	Price *uint256.Int

	// Timestamp tracks the time at which the price was fetched.
	Timestamp time.Time
}

// NewQuotePrice returns a new QuotePrice.
func NewQuotePrice(price *uint256.Int, timestamp time.Time) (QuotePrice, error) {
	if price == nil {
		return QuotePrice{}, fmt.Errorf("last price cannot be nil")
	}

	return QuotePrice{Price: price, Timestamp: timestamp}, nil
}
