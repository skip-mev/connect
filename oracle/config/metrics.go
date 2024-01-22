package config

import (
	"fmt"
	"net/url"
)

// MetricsConfig is the metrics configurations for the oracle. This configuration object specifically
// exposes metrics pertaining to the oracle side car. To enable app side metrics, please see the app
// configuration.
type MetricsConfig struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose
	// metrics to.
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// Enabled indicates whether metrics should be enabled.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}

// ValidateBasic performs basic validation of the config.
func (c *MetricsConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if _, err := url.ParseRequestURI(c.PrometheusServerAddress); err != nil {
		return fmt.Errorf("must supply a valid prometheus server address if metrics are enabled: %w", err)
	}

	return nil
}
