package config

import (
	"fmt"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// ProviderMarketConfig represents the provider specific configurations for different
// currency pairs and the corresponding markets they are traded on.
type ProviderMarketConfig struct {
	// Name identifies which provider this config is for.
	Name string `mapstructure:"name" toml:"name"`

	// CurrencyPairToMarketConfigs is the config the provider uses to create mappings
	// between on-chain and off-chain currency pairs. In particular, this config
	// maps the on-chain currency pair representation (i.e. BITCOIN/USD) to the
	// off-chain currency pair representation (i.e. BTC/USD).
	CurrencyPairToMarketConfigs map[string]CurrencyPairMarketConfig `mapstructure:"currency_pair_to_market_configs" toml:"currency_pair_to_market_configs"`
}

// CurrencyPairMarketConfig is the config the provider uses to create mappings
// between on-chain and off-chain currency pairs.
type CurrencyPairMarketConfig struct {
	// Ticker is the ticker symbol for the currency pair.
	Ticker string `mapstructure:"ticker" toml:"ticker"`

	// CurrencyPair is the on-chain representation of the currency pair.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`
}

func (c *ProviderMarketConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	for cp, marketConfig := range c.CurrencyPairToMarketConfigs {
		if _, err := oracletypes.CurrencyPairFromString(cp); err != nil {
			return fmt.Errorf("currency pair is not formatted correctly %w", err)
		}

		if err := marketConfig.ValidateBasic(); err != nil {
			return fmt.Errorf("market config is not formatted correctly %w", err)
		}
	}

	return nil
}

func (c *CurrencyPairMarketConfig) ValidateBasic() error {
	if len(c.Ticker) == 0 {
		return fmt.Errorf("ticker cannot be empty")
	}

	return nil
}
