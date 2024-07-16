package config

import (
	"fmt"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	// MaxInterval is the maximum time between each price update request.
	MaxInterval = 1 * time.Minute
	// MaxMaxAge is the maximum maximum age of the latest price response before it is
	// considered stale.
	MaxMaxAge = 1 * time.Minute
)

const (
	// DefaultConfigTemplate should be utilized in the app.toml file.
	// This template configures the application to connect to the
	// oracle sidecar and exposes instrumentation for the oracle client
	// and the interaction between the oracle and the app.
	DefaultConfigTemplate = `

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Enabled indicates whether the oracle is enabled.
enabled = "{{ .Oracle.Enabled }}"

# Oracle Address is the URL of the out of process oracle sidecar. This is used to
# connect to the oracle sidecar when the application boots up. Note that the address
# can be modified at any point, but will only take effect after the application is
# restarted. This can be the address of an oracle container running on the same
# machine or a remote machine.
oracle_address = "{{ .Oracle.OracleAddress }}"

# Client Timeout is the time that the client is willing to wait for responses from 
# the oracle before timing out.
client_timeout = "{{ .Oracle.ClientTimeout }}"

# MetricsEnabled determines whether oracle metrics are enabled. Specifically
# this enables instrumentation of the oracle client and the interaction between
# the oracle and the app.
metrics_enabled = "{{ .Oracle.MetricsEnabled }}"

# MaxAge is the maximum age of the latest price response before it is considered stale. 
# The recommended max age is 30 seconds. If this is greater than 1 minute, the app
# will not start.
max_age = "{{ .Oracle.MaxAge }}"

# Interval is the time between each price update request. The recommended interval
# is the block time of the chain. Otherwise, 2 seconds is a good default. If this
# is greater than 1 minute, the app will not start.
interval = "{{ .Oracle.Interval }}"
`
)

const (
	flagEnabled                 = "oracle.enabled"
	flagOracleAddress           = "oracle.oracle_address"
	flagClientTimeout           = "oracle.client_timeout"
	flagMetricsEnabled          = "oracle.metrics_enabled"
	flagPrometheusServerAddress = "oracle.prometheus_server_address"
	flagMaxAge                  = "oracle.max_age"
	flagInterval                = "oracle.interval"
)

// AppConfig contains the application side oracle configurations that must
// be set in the app.toml file.
type AppConfig struct {
	// Enabled indicates whether the oracle is enabled.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`

	// OracleAddress is the URL of the out of process oracle sidecar. This is
	// used to connect to the oracle sidecar.
	OracleAddress string `mapstructure:"oracle_address" toml:"oracle_address"`

	// ClientTimeout is the time that the client is willing to wait for responses
	// from the oracle before timing out.
	ClientTimeout time.Duration `mapstructure:"client_timeout" toml:"client_timeout"`

	// MetricsEnabled is a flag that determines whether oracle metrics are enabled.
	MetricsEnabled bool `mapstructure:"metrics_enabled" toml:"metrics_enabled"`

	// MaxAge is the maximum age of the latest price response before it is considered stale.
	MaxAge time.Duration `mapstructure:"max_age" toml:"max_age"`

	// Interval is the time between each price update request.
	Interval time.Duration `mapstructure:"interval" toml:"interval"`
}

// ValidateBasic performs basic validation of the app config.
func (c *AppConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if len(c.OracleAddress) == 0 {
		return fmt.Errorf("oracle address must not be empty")
	}

	if c.ClientTimeout <= 0 {
		return fmt.Errorf("oracle client timeout must be greater than 0")
	}

	if c.MaxAge <= 0 || c.MaxAge > MaxMaxAge {
		return fmt.Errorf("oracle max age must be between 0 and %s", MaxMaxAge)
	}

	if c.Interval <= 0 || c.Interval > MaxInterval {
		return fmt.Errorf("oracle interval must be between 0 and %s", MaxInterval)
	}

	if c.Interval >= c.MaxAge {
		return fmt.Errorf("oracle interval must be strictly less than max age")
	}

	return nil
}

// ReadConfigFromFile reads a config from a file and returns the config.
func ReadConfigFromFile(path string) (AppConfig, error) {
	var config AppConfig

	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := config.ValidateBasic(); err != nil {
		return config, err
	}

	return config, nil
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (AppConfig, error) {
	var (
		cfg AppConfig
		err error
	)

	// determine if the oracle is enabled
	if v := opts.Get(flagEnabled); v != nil {
		if cfg.Enabled, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	// get the oracle address
	if v := opts.Get(flagOracleAddress); v != nil {
		if cfg.OracleAddress, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	// get the client timeout
	if v := opts.Get(flagClientTimeout); v != nil {
		if cfg.ClientTimeout, err = cast.ToDurationE(v); err != nil {
			return cfg, err
		}
	}

	// get the metrics enabled
	if v := opts.Get(flagMetricsEnabled); v != nil {
		if cfg.MetricsEnabled, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	// get the max age
	if v := opts.Get(flagMaxAge); v != nil {
		if cfg.MaxAge, err = cast.ToDurationE(v); err != nil {
			return cfg, err
		}
	}

	// get the interval
	if v := opts.Get(flagInterval); v != nil {
		if cfg.Interval, err = cast.ToDurationE(v); err != nil {
			return cfg, err
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return cfg, err
	}

	return cfg, err
}
