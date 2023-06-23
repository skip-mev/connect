package oracle

import (
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	otypes "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/coinbase"
	"github.com/skip-mev/slinky/providers/coingecko"
	"github.com/spf13/viper"
)

type Config struct {
	UpdateInterval time.Duration           `mapstructure:"update_interval"`
	Providers      []otypes.ProviderConfig `mapstructure:"providers"`
	CurrencyPairs  []otypes.CurrencyPair   `mapstructure:"currency_pairs"`
}

func ReadConfigFromFile(path string) (*Config, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func providerFromProviderConfig(cfg otypes.ProviderConfig, cps []otypes.CurrencyPair, l log.Logger) (otypes.Provider, error) {
	switch cfg.Name {
	case "coingecko":
		return coingecko.NewProvider(l, cps), nil
	case "coinbase":
		return coinbase.NewProvider(l, cps), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
}

// Providers returns the set of providers that this config determines from the provider configs
func (c *Config) GetProviders(l log.Logger) ([]otypes.Provider, error) {
	providers := make([]otypes.Provider, len(c.Providers))
	var err error

	for i, p := range c.Providers {
		providers[i], err = providerFromProviderConfig(p, c.CurrencyPairs, l)
		if err != nil {
			return nil, err
		}
	}

	return providers, nil
}
