package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// OracleConfig is the over-arching config for the oracle sidecar and instrumentation. The
// oracle is configured via a set of data providers (i.e. coinbase, binance, etc.) and a set
// of currency pairs (i.e. BTC/USD, ETH/USD, etc.). The oracle will fetch prices from the
// data providers for the currency pairs at the specified update interval.
type OracleConfig struct {
	// UpdateInterval is the interval at which the oracle will fetch prices from providers.
	UpdateInterval time.Duration `json:"updateInterval"`

	// Providers is the list of providers that the oracle will fetch prices from.
	Providers []ProviderConfig `json:"providers"`

	// Production specifies whether the oracle is running in production mode. This is used to
	// determine whether the oracle should be run in debug mode or not.
	Production bool `json:"production"`

	// Metrics is the metrics configurations for the oracle.
	Metrics MetricsConfig `json:"metrics"`
}

// ValidateBasic performs basic validation on the oracle config.
func (c *OracleConfig) ValidateBasic() error {
	if c.UpdateInterval <= 0 {
		return fmt.Errorf("oracle update interval must be greater than 0")
	}

	for _, p := range c.Providers {
		if err := p.ValidateBasic(); err != nil {
			return fmt.Errorf("provider is not formatted correctly: %w", err)
		}
	}

	return c.Metrics.ValidateBasic()
}

// ReadOracleConfigFromFile reads a config from a file and returns the config.
func ReadOracleConfigFromFile(path string) (OracleConfig, error) {
	// Read in config file.
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return OracleConfig{}, err
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
