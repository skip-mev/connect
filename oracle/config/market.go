package config

import (
	"fmt"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// MarketConfig represents the provider specific configurations for different
// currency pairs and the corresponding markets they are traded on.
type MarketConfig struct {
	// Name identifies which provider this config is for.
	Name string `mapstructure:"name" toml:"name"`

	// CurrencyPairToMarketConfigs is the config the provider uses to create mappings
	// between on-chain and off-chain currency pairs. In particular, this config
	// maps the on-chain currency pair representation (i.e. BITCOIN/USD) to the
	// off-chain currency pair representation (i.e. BTC/USD).
	CurrencyPairToMarketConfigs map[string]CurrencyPairMarketConfig `mapstructure:"currency_pair_to_market_configs" toml:"currency_pair_to_market_configs"`

	// MarketToCurrencyPairConfigs is the config the provider uses to create mappings
	// between off-chain and on-chain currency pairs. In particular, this config
	// maps the off-chain currency pair representation (i.e. BTC/USD) to the
	// on-chain currency pair representation (i.e. BITCOIN/USD).
	MarketToCurrencyPairConfigs map[string]CurrencyPairMarketConfig
}

// CurrencyPairMarketConfig is the config the provider uses to create mappings
// between on-chain and off-chain currency pairs.
type CurrencyPairMarketConfig struct {
	// Ticker is the ticker symbol for the currency pair off-chain.
	Ticker string `mapstructure:"ticker" toml:"ticker"`

	// CurrencyPair is the on-chain representation of the currency pair.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`
}

// NewMarketConfig returns a new MarketConfig instance.
func NewMarketConfig() MarketConfig {
	return MarketConfig{
		CurrencyPairToMarketConfigs: make(map[string]CurrencyPairMarketConfig),
	}
}

// Invert returns the inverted currency pair market config. This is used to
// create the inverse currency pair market config for the provider.
func (c *MarketConfig) Invert() map[string]CurrencyPairMarketConfig {
	marketToCPConfig := make(map[string]CurrencyPairMarketConfig)

	for _, marketConfig := range c.CurrencyPairToMarketConfigs {
		marketToCPConfig[marketConfig.Ticker] = marketConfig
	}

	return marketToCPConfig
}

// ValidateBasic performs basic validation of the market config.
func (c *MarketConfig) ValidateBasic() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	if len(c.CurrencyPairToMarketConfigs) == 0 {
		return fmt.Errorf("market config must have at least one currency pair")
	}

	for cpStr, marketConfig := range c.CurrencyPairToMarketConfigs {
		cp, err := oracletypes.CurrencyPairFromString(cpStr)
		if err != nil {
			return fmt.Errorf("currency pair is not formatted correctly %w", err)
		}

		if err := marketConfig.ValidateBasic(); err != nil {
			return fmt.Errorf("market config is not formatted correctly %w", err)
		}

		// Update the correctly formatted currency pair string.
		delete(c.CurrencyPairToMarketConfigs, cpStr)
		c.CurrencyPairToMarketConfigs[cp.ToString()] = marketConfig
	}

	// Invert the currency pair market config.
	c.MarketToCurrencyPairConfigs = c.Invert()

	return nil
}

// ValidateBasic performs basic validation of the currency pair market config.
func (c *CurrencyPairMarketConfig) ValidateBasic() error {
	if len(c.Ticker) == 0 {
		return fmt.Errorf("ticker cannot be empty")
	}

	return c.CurrencyPair.ValidateBasic()
}
