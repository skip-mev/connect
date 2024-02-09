package types

import "fmt"

// NewMarketConfig returns a new MarketConfig instance.
func NewMarketConfig(provider string, configs map[string]TickerConfig) MarketConfig {
	return MarketConfig{
		Name:          provider,
		TickerConfigs: configs,
	}
}

// Tickers returns all of the tickers in the MarketConfig.
func (c MarketConfig) Tickers() []Ticker {
	tickers := make([]Ticker, 0, len(c.TickerConfigs))
	for _, cfg := range c.TickerConfigs {
		tickers = append(tickers, cfg.Ticker)
	}

	return tickers
}

// ValidateBasic performs basic validation on the MarketConfig.
func (c *MarketConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	if len(c.TickerConfigs) == 0 {
		return fmt.Errorf("ticker configs cannot be empty")
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
