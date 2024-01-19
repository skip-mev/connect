package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// OracleConfig is the over-arching config for the oracle sidecar and instrumentation. The
// oracle is configured via a set of data providers (i.e. coinbase, binance, etc.) and a set
// of currency pairs (i.e. BTC/USD, ETH/USD, etc.). The oracle will fetch prices from the
// data providers for the currency pairs at the specified update interval.
type OracleConfig struct {
	// UpdateInterval is the interval at which the oracle will fetch prices from providers
	UpdateInterval time.Duration `mapstructure:"update_interval" toml:"update_interval"`

	// Providers is the list of providers that the oracle will fetch prices from.
	Providers []ProviderConfig `mapstructure:"providers" toml:"providers"`

	// CurrencyPairs is the list of currency pairs that the oracle will fetch prices for.
	CurrencyPairs []oracletypes.CurrencyPair `mapstructure:"currency_pairs" toml:"currency_pairs"`

	// Production specifies whether the oracle is running in production mode. This is used to
	// determine whether the oracle should be run in debug mode or not.
	Production bool `mapstructure:"production" toml:"production"`

	// MetricsConfig is the metrics configurations for the oracle. This configuration object allows for
	// metrics tracking of the oracle and the interaction between the oracle and the app.
	Metrics MetricsConfig `mapstructure:"metrics" toml:"metrics"`
}

// ValidateBasic performs basic validation on the oracle config.
func (c *OracleConfig) ValidateBasic() error {
	if c.UpdateInterval <= 0 {
		return fmt.Errorf("oracle update interval must be greater than 0")
	}

	for _, cp := range c.CurrencyPairs {
		if err := cp.ValidateBasic(); err != nil {
			return fmt.Errorf("currency pair is not formatted correctly %w", err)
		}
	}

	for _, p := range c.Providers {
		if err := p.ValidateBasic(); err != nil {
			return fmt.Errorf("provider is not formatted correctly %w", err)
		}
	}

	return c.Metrics.ValidateBasic()
}

// ReadOracleConfigFromFile reads a config from a file and returns the config.
func ReadOracleConfigFromFile(path string) (OracleConfig, error) {
	// Read in config file.
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return OracleConfig{}, err
	}

	// Check required fields.
	requiredFields := []string{"update_interval", "providers", "currency_pairs"}
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			return OracleConfig{}, fmt.Errorf("required field %s is missing in config", field)
		}
	}

	// Unmarshal the config.
	var config OracleConfig
	if err := viper.Unmarshal(&config); err != nil {
		return OracleConfig{}, err
	}

	if err := config.ValidateBasic(); err != nil {
		return OracleConfig{}, err
	}

	return config, nil
}
