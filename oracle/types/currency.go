package types

import (
	"fmt"
	"strconv"
	"strings"
)

// CurrencyPair defines a currency exchange pair consisting of a base and a quote.
type CurrencyPair struct {
	// Base defines the base currency.
	Base string `mapstructure:"base"`

	// Quote defines the quote currency i.e. the currency that the base currency
	// is being exchanged for.
	Quote string `mapstructure:"quote"`

	// QuoteDecimals defines the number of decimals for the quote currency.
	QuoteDecimals int `mapstructure:"quote_decimals"`
}

func NewCurrencyPair(base, quote string, decimals int) CurrencyPair {
	return CurrencyPair{
		Base:          base,
		Quote:         quote,
		QuoteDecimals: decimals,
	}
}

// NewCurrencyPairFromString returns a new CurrencyPair from a string. The string
// must be in the form of "base/quote/quote_decimals" i.e. "BTC/USD/8".
func NewCurrencyPairFromString(asset string) (CurrencyPair, error) {
	pair := strings.Split(asset, "/")
	if len(pair) != 3 {
		return CurrencyPair{}, fmt.Errorf("invalid currency pair: %s", pair)
	}

	decimals, err := strconv.Atoi(pair[2])
	if err != nil {
		return CurrencyPair{}, fmt.Errorf("failed to retrieve quote decimals: %s", err)
	}

	return CurrencyPair{
		Base:          pair[0],
		Quote:         pair[1],
		QuoteDecimals: decimals,
	}, nil
}

// String implements the Stringer interface and defines a ticker symbol for
// querying the exchange rate.
func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s/%d", cp.Base, cp.Quote, cp.QuoteDecimals)
}
