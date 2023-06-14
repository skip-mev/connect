package types

import (
	"fmt"
	"strings"
)

// CurrencyPair defines a currency exchange pair consisting of a base and a quote.
type CurrencyPair struct {
	Base  string
	Quote string
}

func NewCurrencyPair(ticker string) (CurrencyPair, error) {
	tokens := strings.Split(ticker, "/")
	if len(tokens) != 2 {
		return CurrencyPair{}, fmt.Errorf("invalid ticker %s", ticker)
	}

	return CurrencyPair{
		Base:  strings.ToUpper(tokens[0]),
		Quote: strings.ToUpper(tokens[1]),
	}, nil
}

// String implements the Stringer interface and defines a ticker symbol for
// querying the exchange rate.
func (cp CurrencyPair) String() string {
	return cp.Base + "/" + cp.Quote
}

// GetBase returns the base currency.
func (cp CurrencyPair) GetBase() string {
	return cp.Base
}

// GetQuote returns the quote currency.
func (cp CurrencyPair) GetQuote() string {
	return cp.Quote
}
