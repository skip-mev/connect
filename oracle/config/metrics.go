package config

import (
	"fmt"

	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MetricsConfig is the metrics configurations for the oracle. This configuration object allows for
// metrics tracking of the oracle and the interaction between the oracle and the app.
type MetricsConfig struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// OracleMetrics is the config for the oracle metrics
	OracleMetrics OracleMetricsConfig `mapstructure:"oracle_metrics" toml:"oracle_metrics"`

	// AppMetrics is the config for the app metrics
	AppMetrics AppMetricsConfig `mapstructure:"app_metrics" toml:"app_metrics"`
}

// ValidateBasic performs basic validation of the config.
func (c *MetricsConfig) ValidateBasic() error {
	if (c.OracleMetrics.Enabled || c.AppMetrics.Enabled) && len(c.PrometheusServerAddress) == 0 {
		return fmt.Errorf("must supply a prometheus server address if metrics are enabled")
	}

	return c.AppMetrics.ValidateBasic()
}

// OracleMetricsConfig is the config for the oracle metrics. These are also utililized to enable
// provider metrics.
type OracleMetricsConfig struct {
	// Enabled indicates whether metrics should be enabled.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}

// AppMetricsConfig is the config for the app metrics. Specifically, this is used to enable validator
// metrics as proposals and vote extensions are being constructed. Nodes that are not participating
// in consensus should not enable app metrics.
type AppMetricsConfig struct {
	// Enabled indicates whether app side metrics should be enabled.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`

	// ValidatorConsAddress is the validator's consensus address. Validator's must register their
	// consensus address in order to enable app side metrics.
	ValidatorConsAddress string `mapstructure:"validator_cons_address" toml:"validator_cons_address"`
}

// ValidateBasic performs basic validation of the config.
func (c *AppMetricsConfig) ValidateBasic() error {
	if c.Enabled {
		_, err := sdk.ConsAddressFromBech32(c.ValidatorConsAddress)
		return err
	}

	return nil
}

// ConsAddress returns the validator's consensus address
func (c *AppMetricsConfig) ConsAddress() (sdk.ConsAddress, error) {
	if c.Enabled {
		return sdk.ConsAddressFromBech32(c.ValidatorConsAddress)
	}

	return nil, nil
}

// ReadMetricsConfigFromFile reads a config from a file and returns the config.
func ReadMetricsConfigFromFile(path string) (MetricsConfig, error) {
	// Read in config file.
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return MetricsConfig{}, err
	}

	// Check required fields.
	requiredFields := []string{"prometheus_server_address", "oracle_metrics", "app_metrics"}
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			return MetricsConfig{}, fmt.Errorf("required field %s is missing in config", field)
		}
	}

	// Unmarshal the config.
	var config MetricsConfig
	if err := viper.Unmarshal(&config); err != nil {
		return MetricsConfig{}, err
	}

	if err := config.ValidateBasic(); err != nil {
		return MetricsConfig{}, err
	}

	return config, nil
}
