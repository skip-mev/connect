package types

import "fmt"

// NewMarketConfig returns a new MarketConfig instance. The MarketConfig represents
// the provider specific configurations for different markets and the associated
// markets they are traded on i.e. all of the price feeds that a given provider
// is responsible for maintaining.
func NewMarketConfig(provider string, configs map[string]TickerConfig) MarketConfig {
	return MarketConfig{
		Name:          provider,
		TickerConfigs: configs,
	}
}

// Tickers returns all of the tickers that the provider supports.
func (c *MarketConfig) Tickers() []Ticker {
	tickers := make([]Ticker, 0, len(c.TickerConfigs))

	i := 0
	for _, cfg := range c.TickerConfigs {
		tickers[i] = cfg.Ticker
		i++
	}

	return tickers
}

// ValidateBasic performs basic validation on the MarketConfig.
func (c *MarketConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	// The provider must support at least one ticker.
	if len(c.TickerConfigs) == 0 {
		return fmt.Errorf("ticker configs cannot be empty")
	}

	seen := make(map[string]struct{})
	seenOffChain := make(map[string]struct{})
	for ticker, cfg := range c.TickerConfigs {
		// Validate the ticker configurations.
		if err := cfg.ValidateBasic(); err != nil {
			return err
		}

		// The ticker key should match the ticker value.
		t := cfg.Ticker.String()
		if ticker != t {
			return fmt.Errorf("ticker config key does not match ticker value; expected %s, got %s", ticker, t)
		}

		// Check for duplicate tickers.
		if _, ok := seen[t]; ok {
			return fmt.Errorf("duplicate ticker found: %s", t)
		}
		seen[t] = struct{}{}

		// Check for duplicate off-chain tickers.
		if _, ok := seenOffChain[cfg.OffChainTicker]; ok {
			return fmt.Errorf("duplicate off-chain ticker found: %s", cfg.OffChainTicker)
		}
		seenOffChain[cfg.OffChainTicker] = struct{}{}
	}

	return nil
}
