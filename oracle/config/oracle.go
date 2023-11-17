package config

import (
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
	// InProcess specifies whether the oracle configured, is currently running as a remote grpc-server, or will be run in process
	InProcess bool `mapstructure:"in_process" toml:"in_process"`

	// RemoteAddress is the address of the remote oracle server (if it is running out-of-process)
	RemoteAddress string `mapstructure:"remote_address" toml:"remote_address"`

	// Timeout is the time that the client is willing to wait for responses from the oracle before timing out.
	Timeout time.Duration `mapstructure:"timeout" toml:"timeout"`

	// UpdateInterval is the interval at which the oracle will fetch prices from providers
	UpdateInterval time.Duration `mapstructure:"update_interval" toml:"update_interval"`

	// Providers is the list of providers that the oracle will fetch prices from.
	Providers []ProviderConfig `mapstructure:"providers" toml:"providers"`

	// CurrencyPairs is the list of currency pairs that the oracle will fetch prices for.
	CurrencyPairs []oracletypes.CurrencyPair `mapstructure:"currency_pairs" toml:"currency_pairs"`
}

// ReadOracleConfigFromFile reads a config from a file and returns the config.
func ReadOracleConfigFromFile(path string) (OracleConfig, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return OracleConfig{}, err
	}

	// unmarshal config
	var config OracleConfig
	if err := viper.Unmarshal(&config); err != nil {
		return OracleConfig{}, err
	}

	return config, nil
}
