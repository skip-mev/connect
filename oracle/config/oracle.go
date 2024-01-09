package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// OracleConfig is the base config for both out-of-process and in-process oracles.
// If the oracle is to be configured out-of-process in base-app, a grpc-client of
// the grpc-server running at RemoteAddress is instantiated, otherwise, an in-process
// local client oracle is instantiated. Note, that you can only have one oracle
// running at a time.
type OracleConfig struct {
	// Enabled specifies whether the side-car oracle needs to be run.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`

	// InProcess specifies whether the oracle configured, is currently running as a remote grpc-server, or will be run in process
	InProcess bool `mapstructure:"in_process" toml:"in_process"`

	// RemoteAddress is the address of the remote oracle server (if it is running out-of-process)
	RemoteAddress string `mapstructure:"remote_address" toml:"remote_address"`

	// ClientTimeout is the time that the client is willing to wait for responses from the oracle before timing out.
	ClientTimeout time.Duration `mapstructure:"client_timeout" toml:"client_timeout"`

	// UpdateInterval is the interval at which the oracle will fetch prices from providers
	UpdateInterval time.Duration `mapstructure:"update_interval" toml:"update_interval"`

	// Providers is the list of providers that the oracle will fetch prices from.
	Providers []ProviderConfig `mapstructure:"providers" toml:"providers"`

	// CurrencyPairs is the list of currency pairs that the oracle will fetch prices for.
	CurrencyPairs []oracletypes.CurrencyPair `mapstructure:"currency_pairs" toml:"currency_pairs"`

	// Production specifies whether the oracle is running in production mode. This is used to
	// determine whether the oracle should be run in debug mode or not.
	Production bool `mapstructure:"production" toml:"production"`
}

// ValidateBasic performs basic validation on the oracle config.
func (c *OracleConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if !c.InProcess && len(c.RemoteAddress) == 0 {
		return fmt.Errorf("must supply a remote address if the oracle is running out of process")
	}

	if c.UpdateInterval <= 0 || c.ClientTimeout <= 0 {
		return fmt.Errorf("oracle update interval and client timeout must be greater than 0")
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

	return nil
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
	requiredFields := []string{"enabled", "in_process", "remote_address", "client_timeout", "update_interval", "production"}
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
