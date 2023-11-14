package config

import (
	oracle_metrics "github.com/skip-mev/slinky/oracle/metrics"
	service_metrics "github.com/skip-mev/slinky/service/metrics"
	"github.com/spf13/viper"
)

// MetricsConfig is the metrics configurations for the oracle. This configuration object allows for
// metrics tracking of the oracle and the interaction between the oracle and the app.
type MetricsConfig struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// OracleMetrics is the config for the oracle metrics
	OracleMetrics oracle_metrics.Config `mapstructure:"oracle_metrics" toml:"oracle_metrics"`

	// AppMetrics is the config for the app metrics
	AppMetrics service_metrics.Config `mapstructure:"app_metrics" toml:"app_metrics"`
}

// ReadMetricsConfigFromFile reads a config from a file and returns the config.
func ReadMetricsConfigFromFile(path string) (MetricsConfig, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return MetricsConfig{}, err
	}

	// unmarshal config
	var config MetricsConfig
	if err := viper.Unmarshal(&config); err != nil {
		return MetricsConfig{}, err
	}

	return config, nil
}
