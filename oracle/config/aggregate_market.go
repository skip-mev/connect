package config

import (
	"fmt"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// AggregateMarketConfig represents the market configurations for how currency pairs will
// be resolved to a final price. Each currency pair can have a list of convertable markets
// that will be used to convert the price of the currency pair to a common currency pair.
type AggregateMarketConfig struct {
	// CurrencyPairs is the list of currency pairs that the oracle will fetch aggregate prices
	// for.
	CurrencyPairs map[string]AggregateCurrencyPairConfig `mapstructure:"currency_pairs" toml:"currency_pairs"`
}

// AggregateCurrencyPairConfig represents the configurations for a single currency pair.
type AggregateCurrencyPairConfig struct {
	// CurrencyPair is the currency pair that the oracle will fetch prices for.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`

	// ConvertableMarkets is the list of valid markets (i.e. price feeds) that can be
	// used to convert the price of the currency pair to a common currency pair. For
	// example, if the oracle receives a price for BTC/USDT and USDT/USD, it can use
	// the conversion market to convert the BTC/USDT price to BTC/USD. These must be
	// provided in a topologically sorted order.
	ConvertableMarkets [][]ConvertableMarket `mapstructure:"convertable_markets" toml:"convertable_markets"`
}

// CovertableMarket returns the list of convertable markets for the currency pair.
type ConvertableMarket struct {
	// CurrencyPair is the feed that will be used in the conversion.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`

	// Invert is a flag that indicates if the feed should be inverted
	// prior to being used in the conversion.
	Invert bool `mapstructure:"invert" toml:"invert"`
}

// NewAggregateMarketConfig returns a new AggregateMarketConfig.
func NewAggregateMarketConfig() AggregateMarketConfig {
	return AggregateMarketConfig{
		CurrencyPairs: make(map[string]AggregateCurrencyPairConfig),
	}
}

// GetCurrencyPairs returns the set of currency pairs in the aggregate market config.
func (c AggregateMarketConfig) GetCurrencyPairs() []oracletypes.CurrencyPair {
	var currencyPairs []oracletypes.CurrencyPair

	for _, cpConfig := range c.CurrencyPairs {
		currencyPairs = append(currencyPairs, cpConfig.CurrencyPair)
	}

	return currencyPairs
}

// ValidateBasic performs basic validation on the AggregateMarketConfig.
func (c AggregateMarketConfig) ValidateBasic() error {
	for cpString, cpConfig := range c.CurrencyPairs {
		cp, err := oracletypes.CurrencyPairFromString(cpString)
		if err != nil {
			return err
		}

		// The currency pair in the config must match the key.
		if cp != cpConfig.CurrencyPair {
			return fmt.Errorf("currency pair %s does not match the currency pair in the config", cpString)
		}

		if err := cpConfig.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic performs basic validation on the AggregateCurrencyPairConfig.
func (c AggregateCurrencyPairConfig) ValidateBasic() error {
	if err := c.CurrencyPair.ValidateBasic(); err != nil {
		return err
	}

	// Check that order is topologically sorted for each market. For example, if the oracle
	// receives a price for BTC/USDT and USDT/USD, the order must be BTC/USDT -> USDT/USD.
	// Alternatively, if the oracle receives a price for BTC/USDT and USD/USDT, the order must
	// be BTC/USDT -> USD/USDT (inverted == true).
	for _, conversions := range c.ConvertableMarkets {
		if len(conversions) <= 1 {
			return fmt.Errorf("at least two markets must be provided in order for a viable conversion to occur")
		}

		// Basic check to ensure the base and quote denom of the first and last market match the
		// currency pair in the config.
		base := conversions[0].CurrencyPair.Base
		if conversions[0].Invert {
			base = conversions[0].CurrencyPair.Quote
		}

		quote := conversions[len(conversions)-1].CurrencyPair.Quote
		if conversions[len(conversions)-1].Invert {
			quote = conversions[len(conversions)-1].CurrencyPair.Base
		}

		if base != c.CurrencyPair.Base || quote != c.CurrencyPair.Quote {
			return fmt.Errorf("invalid convertable market; expected %s/%s but got %s/%s", c.CurrencyPair.Base, c.CurrencyPair.Quote, base, quote)
		}

		// Check that the order is topologically sorted.
		quote = conversions[0].CurrencyPair.Quote
		if conversions[0].Invert {
			quote = conversions[0].CurrencyPair.Base
		}
		for _, cpConfig := range conversions[1:] {
			if err := cpConfig.CurrencyPair.ValidateBasic(); err != nil {
				return err
			}

			switch {
			case !cpConfig.Invert && quote == cpConfig.CurrencyPair.Base:
				quote = cpConfig.CurrencyPair.Quote
			case cpConfig.Invert && quote == cpConfig.CurrencyPair.Quote:
				quote = cpConfig.CurrencyPair.Base
			case !cpConfig.Invert && quote != cpConfig.CurrencyPair.Base:
				return fmt.Errorf("invalid convertable market; expected %s but got %s", quote, cpConfig.CurrencyPair.Base)
			case cpConfig.Invert && quote != cpConfig.CurrencyPair.Quote:
				return fmt.Errorf("invalid convertable market; expected %s but got %s", quote, cpConfig.CurrencyPair.Quote)
			}
		}
	}

	return nil
}
