package types

import (
	fmt "fmt"
)

// NewAggregateMarketConfig returns a new AggregateMarketConfig instance. The
// AggregateMarketConfig represents the global set of market configurations for
// all providers that will be utilized off-chain as well as how tickers will be
// resolved to a final price. Each ticker can have a list of convertable markets
// that will be used to convert the prices of a set of tickers to a common
// ticker.
//
// Price aggregation broadly follows the following steps:
//  1. Fetch prices for each ticker from the providers.
//  2. Calculate the final price for each ticker. This is dependent on the
//     aggregation
//     strategy used by the oracle. The oracle may use a median price, a
//     weighted average price, etc.
//  3. Convert the price of each ticker to a common ticker using the aggreagted
//     ticker
//     configurations by default. The oracle may use a different aggregation
//     strategy to convert the price of a ticker to a common ticker.
//
// For example, the oracle may be configured with the feeds:
//   - BTC/USDT
//   - USDT/USD
//   - BTC/USDC
//   - USDC/USD
//
// The aggregated ticker may be:
//   - BTC/USD: (calculate a median price from the following convertable
//     markets)
//     1. BTC/USDT -> USDT/USD = BTC/USD
//     2. BTC/USDC -> USDC/USD = BTC/USD
func NewAggregateMarketConfig(markets map[string]MarketConfig, tickers map[string]PathsConfig) (AggregateMarketConfig, error) {
	c := AggregateMarketConfig{
		MarketConfigs: markets,
		TickerConfigs: tickers,
	}

	if err := c.ValidateBasic(); err != nil {
		return AggregateMarketConfig{}, err
	}

	return c, nil
}

// ValidateBasic performs basic validation on the AggregateMarketConfig.
func (c AggregateMarketConfig) ValidateBasic() error {
	// Track all of the tickers that are supported by all providers to ensure that
	// all market conversions are supported by at least one provider.
	seenTickers := make(map[Ticker]struct{})
	for name, market := range c.MarketConfigs {
		if err := market.ValidateBasic(); err != nil {
			return err
		}

		if name != market.Name {
			return fmt.Errorf("market config key does not match market value; expected %s, got %s", name, market.Name)
		}

		for _, ticker := range market.Tickers() {
			seenTickers[ticker] = struct{}{}
		}
	}

	// Validate all of the conversion paths for each ticker.
	for ticker, cfg := range c.TickerConfigs {
		if err := cfg.ValidateBasic(); err != nil {
			return err
		}

		// The ticker key should match the ticker value.
		if ticker != cfg.Ticker.String() {
			return fmt.Errorf("ticker config key does not match ticker value; expected %s, got %s", ticker, cfg.Ticker.String())
		}

		// Ensure that the target ticker is supported by at least one provider.
		if _, ok := seenTickers[cfg.Ticker]; !ok {
			return fmt.Errorf("ticker not found in market configs: %s", cfg.Ticker.String())
		}

		// Ensure that all of the tickers in the conversion paths are supported by at least one provider.
		for ticker := range cfg.UniqueTickers() {
			if _, ok := seenTickers[ticker]; !ok {
				return fmt.Errorf("ticker not found in market configs: %s", ticker.String())
			}
		}
	}

	return nil
}
