package types

import fmt "fmt"

// NewAggregateMarketConfig returns a new AggregateMarketConfig instance.
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
	for name, market := range c.MarketConfigs {
		if err := market.ValidateBasic(); err != nil {
			return err
		}

		if name != market.Name {
			return fmt.Errorf("market config key does not match market value; expected %s, got %s", name, market.Name)
		}
	}

	for ticker, cfg := range c.TickerConfigs {
		if err := cfg.ValidateBasic(); err != nil {
			return err
		}

		if ticker != cfg.Ticker.String() {
			return fmt.Errorf("ticker config key does not match ticker value; expected %s, got %s", ticker, cfg.Ticker.String())
		}
	}

	return nil
}
