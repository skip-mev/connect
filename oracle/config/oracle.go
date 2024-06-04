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

	// MaxPriceAge is the maximum age of a price that the oracle will consider valid. If a
	// price is older than this, the oracle will not consider it valid and will not return it in /prices
	// requests.
	MaxPriceAge time.Duration `json:"maxPriceAge"`

	// Providers is the list of providers that the oracle will fetch prices from.
	Providers map[string]ProviderConfig `json:"providers"`

	// Metrics is the metrics configurations for the oracle.
	Metrics MetricsConfig `json:"metrics"`

	// Host is the host that the oracle will listen on.
	Host string `json:"host"`

	// Port is the port that the oracle will listen on.
	Port string `json:"port"`
}

// ValidateBasic performs basic validation on the oracle config.
func (c *OracleConfig) ValidateBasic() error {
	if c.UpdateInterval <= 0 {
		return fmt.Errorf("oracle update interval must be greater than 0")
	}

	if c.MaxPriceAge <= 0 {
		return fmt.Errorf("oracle max price age must be greater than 0")
	}

	for _, p := range c.Providers {
		if err := p.ValidateBasic(); err != nil {
			return fmt.Errorf("provider is not formatted correctly: %w", err)
		}
	}

	if len(c.Host) == 0 {
		return fmt.Errorf("oracle host cannot be empty")
	}

	if len(c.Port) == 0 {
		return fmt.Errorf("oracle port cannot be empty")
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
