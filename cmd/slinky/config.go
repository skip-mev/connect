package main

import (
	"fmt"
	"github.com/skip-mev/slinky/providers"
	"github.com/spf13/viper"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// DefaultUpdateInterval is the default value for how frequently slinky updates aggregate price responses.
	DefaultUpdateInterval = 250000000
	// DefaultMaxPriceAge is the default value for the oldest price considered in an aggregate price response by slinky.
	DefaultMaxPriceAge = 120000000000
	// DefaultPrometheusServerAddress is the default value for the prometheus server address in slinky.
	DefaultPrometheusServerAddress = "0.0.0.0:8002"
	// DefaultMetricsEnabled is the default value for enabling prometheus metrics in slinky.
	DefaultMetricsEnabled = true
	// DefaultHost is the default for the slinky oracle server host.
	DefaultHost = "0.0.0.0"
	// DefaultPort is the default for the slinky oracle server port.
	DefaultPort = "8080"
)

type OracleConfig struct {
	// UpdateInterval is the interval at which the oracle will fetch prices from providers.
	UpdateInterval time.Duration `json:"updateInterval"`

	// MaxPriceAge is the maximum age of a price that the oracle will consider valid. If a
	// price is older than this, the oracle will not consider it valid and will not return it in /prices
	// requests.
	MaxPriceAge time.Duration `json:"maxPriceAge"`

	// Providers is the map of provider names to providers that the oracle will fetch prices from.
	Providers map[string]config.ProviderConfig `json:"providers"`

	// Production specifies whether the oracle is running in production mode. This is used to
	// determine whether the oracle should be run in debug mode or not.
	//
	// Deprecated: This field is no longer used.
	Production bool `json:"production"`

	// Metrics is the metrics configurations for the oracle.
	Metrics config.MetricsConfig `json:"metrics"`

	// Host is the host that the oracle will listen on.
	Host string `json:"host"`

	// Port is the port that the oracle will listen on.
	Port string `json:"port"`
}

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

func (c *OracleConfig) ToLegacy() config.OracleConfig {
	providers := make([]config.ProviderConfig, len(c.Providers))
	var i int
	for _, providerConfig := range c.Providers {
		providers[i] = providerConfig
		i++
	}
	return config.OracleConfig{
		UpdateInterval: c.UpdateInterval,
		MaxPriceAge:    c.MaxPriceAge,
		Providers:      providers,
		Production:     c.Production,
		Metrics:        c.Metrics,
		Host:           c.Host,
		Port:           c.Port,
	}
}

func SetDefaults() {
	viper.SetDefault("updateInterval", DefaultUpdateInterval)
	viper.SetDefault("maxPriceAge", DefaultMaxPriceAge)
	viper.SetDefault("metrics.prometheusServerAddress", DefaultPrometheusServerAddress)
	viper.SetDefault("metrics.enabled", DefaultMetricsEnabled)
	viper.SetDefault("host", DefaultHost)
	viper.SetDefault("port", DefaultPort)
	for _, providerConfig := range providers.ProviderDefaults {
		viper.SetDefault(fmt.Sprintf("providers.%s", providerConfig.Name), providerConfig)
	}
}

func GetLegacyOracleConfig(path string) (config.OracleConfig, error) {
	SetDefaults()
	var oracleCfg OracleConfig
	var err error
	if path != "" {
		oracleCfg, err = ReadOracleConfigFromFile(path)
	} else {
		err = viper.Unmarshal(&oracleCfg)
	}
	if err != nil {
		return config.OracleConfig{}, err
	}
	return oracleCfg.ToLegacy(), nil
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
