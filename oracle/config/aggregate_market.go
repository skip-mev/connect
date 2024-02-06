package config

import (
	"fmt"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// AggregateMarketConfig represents the market configurations for how currency pairs will
// be resolved to a final price. Each currency pair can have a list of convertable markets
// that will be used to convert the price of the currency pair to a common currency pair.
//
// Price aggregation broadly follows the following steps:
//  1. Fetch prices for each currency pair from the price feeds.
//  2. Calculate the provider-weighted median price for each currency pair.
//  3. Convert the price of each currency pair to a common currency pair using the
//     AggregatedFeeds field. If there are multiple convertable feeds for a given currency
//     pair, calculate a median price for each final currency pair, weighted by the
//     weights of the convertable feeds.
//
// For example, the oracle may be configured with the feeds:
//   - BTC/USDT
//   - USDT/USD
//   - BTC/USDC
//   - USDC/USD
//
// The aggregated feeds may be:
//   - BTC/USD: (calculate a median price from the following convertable markets)
//     1. BTC/USDT -> USDT/USD = BTC/USD
//     2. BTC/USDC -> USDC/USD = BTC/USD
type AggregateMarketConfig struct {
	// Feeds is a map of all of the price feeds that the oracle will fetch prices for along
	// with the resolution configurations for each feed.
	Feeds map[string]FeedConfig `mapstructure:"currency_pairs" toml:"currency_pairs"`

	// AggregatedFeeds is the list of valid markets (i.e. price feeds) that can be
	// used to convert the price of the currency pair to a common currency pair. For
	// example, if the oracle receives a price for BTC/USDT and USDT/USD, it can use
	// the conversion market to convert the BTC/USDT price to BTC/USD. These must be
	// provided in a topologically sorted order that resolve to the same currency pair
	// defined in the CurrencyPair field.
	AggregatedFeeds map[string][][]Conversion `mapstructure:"aggregated_feeds" toml:"aggregated_feeds"`
}

// FeedConfig represents the configurations for a given price feed. Each currency pair
// will have its own feed configuration.
type FeedConfig struct {
	// CurrencyPair is the currency pair that the oracle will fetch prices for.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`
}

// Conversion represents a price feed that can be used to convert to a final common
// currency pair.
type Conversion struct {
	// CurrencyPair is the feed that will be used in the conversion.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`

	// Invert is a flag that indicates if the feed should be inverted prior to being used
	// in the conversion.
	Invert bool `mapstructure:"invert" toml:"invert"`
}

// GetCurrencyPairs returns the set of currency pairs in the aggregate market config.
func (c *AggregateMarketConfig) GetCurrencyPairs() []oracletypes.CurrencyPair {
	var currencyPairs []oracletypes.CurrencyPair

	for _, cpConfig := range c.Feeds {
		currencyPairs = append(currencyPairs, cpConfig.CurrencyPair)
	}

	return currencyPairs
}

// ValidateBasic performs basic validation on the AggregateMarketConfig.
func (c *AggregateMarketConfig) ValidateBasic() error {
	// Verify the configurations of all price feeds.
	for marketString, feedConfig := range c.Feeds {
		cp, err := oracletypes.CurrencyPairFromID(marketString, feedConfig.CurrencyPair.Decimals)
		if err != nil {
			return err
		}

		// The currency pair in the config must match the key.
		if cp != feedConfig.CurrencyPair {
			return fmt.Errorf("currency pair %s does not match the currency pair in the config", marketString)
		}

		if err := feedConfig.ValidateBasic(); err != nil {
			return err
		}

		// Upper case the currency pair string since toml may not preserve the case.
		delete(c.Feeds, marketString)
		c.Feeds[cp.Ticker()] = feedConfig
	}

	// Ensure that all convertable feeds are valid. We consider it valid if the
	// currency pair can be found in the feeds map and the convertable market is topologically
	// sorted.
	for marketString, convertableFeedsForCP := range c.AggregatedFeeds {
		// check validity of market string
		checkCP, err := oracletypes.CurrencyPairFromID(marketString, oracletypes.DefaultDecimals)
		if err != nil {
			return err
		}

		if len(convertableFeedsForCP) == 0 {
			return fmt.Errorf("no convertable markets provided for %s", marketString)
		}

		for _, feeds := range convertableFeedsForCP {
			for _, conversion := range feeds {
				if _, ok := c.Feeds[conversion.CurrencyPair.Ticker()]; !ok {
					return fmt.Errorf("convertable market %s does not exist in the feeds", conversion.CurrencyPair)
				}
			}

			if err := checkSort(checkCP, feeds); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateBasic performs basic validation on the FeedConfig.
func (c *FeedConfig) ValidateBasic() error {
	if err := c.CurrencyPair.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// checkSort checks if the given list of convertable markets is topologically sorted.
func checkSort(pair oracletypes.CurrencyPair, feeds []Conversion) error {
	// Check that order is topologically sorted for each market. For example, if the oracle
	// receives a price for BTC/USDT and USDT/USD, the order must be BTC/USDT -> USDT/USD.
	// Alternatively, if the oracle receives a price for BTC/USDT and USD/USDT, the order must
	// be BTC/USDT -> USD/USDT (inverted == true).
	if len(feeds) == 0 {
		return fmt.Errorf("at least one markets must be provided in order for a viable conversion to occur")
	}

	if err := feeds[0].CurrencyPair.ValidateBasic(); err != nil {
		return err
	}

	// Basic check to ensure the base and quote denom of the first and last market match the
	// currency pair in the config.
	base := feeds[0].CurrencyPair.Base
	if feeds[0].Invert {
		base = feeds[0].CurrencyPair.Quote
	}

	quote := feeds[len(feeds)-1].CurrencyPair.Quote
	if feeds[len(feeds)-1].Invert {
		quote = feeds[len(feeds)-1].CurrencyPair.Base
	}

	if base != pair.Base || quote != pair.Quote {
		return fmt.Errorf("invalid convertable market; expected %s but got base %s, quote %s", pair.Ticker(), base, quote)
	}

	// Check that the order is topologically sorted.
	quote = feeds[0].CurrencyPair.Quote
	if feeds[0].Invert {
		quote = feeds[0].CurrencyPair.Base
	}
	for _, cpConfig := range feeds[1:] {
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

	return nil
}
