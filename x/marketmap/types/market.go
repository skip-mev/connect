package types

import "fmt"

// NewMarketConfig returns a new MarketConfig instance.
func NewMarketConfig(provider string, configs map[uint64]TickerConfig) MarketConfig {
	return MarketConfig{
		Name:          provider,
		TickerConfigs: configs,
	}
}

// ValidateBasic performs basic validation on the MarketConfig.
func (c *MarketConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	if len(c.TickerConfigs) == 0 {
		return fmt.Errorf("ticker configs cannot be empty")
	}

	for id, cfg := range c.TickerConfigs {
		if err := cfg.ValidateBasic(); err != nil {
			return err
		}

		if id != cfg.Ticker.Id {
			return fmt.Errorf("id %d does not match the id in the config", id)
		}
	}

	return nil
}
