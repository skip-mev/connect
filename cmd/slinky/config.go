package main

import (
	"fmt"
	"github.com/skip-mev/slinky/providers"
	"github.com/spf13/viper"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
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
	viper.Set("updateInterval", 250000000)
	viper.Set("maxPriceAge", 120000000000)
	viper.Set("metrics.prometheusServerAddress", "0.0.0.0:8002")
	viper.Set("metrics.enabled", true)
	viper.Set("host", "0.0.0.0")
	viper.Set("port", "8080")
	for _, providerConfig := range providers.ProviderDefaults {
		viper.Set(fmt.Sprintf("providers.%s", providerConfig.Name), providerConfig)
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
