package types

// CurrencyPair defines a currency exchange pair consisting of a base and a quote.
type CurrencyPair struct {
	// Base defines the base currency.
	Base string

	// Quote defines the quote currency i.e. the currency that the base currency
	// is being exchanged for.
	Quote string

	// QuoteDecimals defines the number of decimals for the quote currency.
	QuoteDecimals int
}

func NewCurrencyPair(base, quote string, decimals int) CurrencyPair {
	return CurrencyPair{
		Base:          base,
		Quote:         quote,
		QuoteDecimals: decimals,
	}
}

// String implements the Stringer interface and defines a ticker symbol for
// querying the exchange rate.
func (cp CurrencyPair) String() string {
	return cp.Base + "/" + cp.Quote
}
