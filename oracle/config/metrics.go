package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// MetricsConfig is the metrics configurations for the oracle. This configuration object allows for
// metrics tracking of the oracle and the interaction between the oracle and the app.
type MetricsConfig struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// Enabled indicates whether metrics should be enabled.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}

// ValidateBasic performs basic validation of the config.
func (c *MetricsConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if len(c.PrometheusServerAddress) == 0 {
		return fmt.Errorf("must supply a prometheus server address if metrics are enabled")
	}

	return nil
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
	requiredFields := []string{"prometheus_server_address", "enabled"}
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
