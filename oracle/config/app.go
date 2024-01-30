package config

import (
	"fmt"
	"net/url"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
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
# restarted.
oracle_address = "{{ .Oracle.OracleAddress }}"

# Client Timeout is the time that the client is willing to wait for responses from 
# the oracle before timing out.
client_timeout = "{{ .Oracle.ClientTimeout }}"

# MetricsEnabled determines whether oracle metrics are enabled. Specifically
# this enables intsrumentation of the oracle client and the interaction between
# the oracle and the app.
metrics_enabled = "{{ .Oracle.MetricsEnabled }}"

# PrometheusServerAddress is the address of the prometheus server that metrics will be
# exposed to.
prometheus_server_address = "{{ .Oracle.PrometheusServerAddress }}"

# ValidatorConsAddress is the validator's consensus address. This is optional. If supplied,
# the oracle will expose metrics for how often the validator's vote extensions are included
# in blocks and more.
validator_cons_address = "{{ .Oracle.ValidatorConsAddress }}"
`
)

const (
	flagEnabled                 = "oracle.enabled"
	flagOracleAddress           = "oracle.oracle_address"
	flagClientTimeout           = "oracle.client_timeout"
	flagMetricsEnabled          = "oracle.metrics_enabled"
	flagPrometheusServerAddress = "oracle.prometheus_server_address"
	flagValidatorConsAddress    = "oracle.validator_cons_address"
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

	// PrometheusServerAddress is the address of the prometheus server that the oracle
	// will expose metrics to.
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// ValidatorConsAddress is the validator's consensus address.
	ValidatorConsAddress string `mapstructure:"validator_cons_address" toml:"validator_cons_address"`
}

// ConsAddress returns the validator's consensus address.
func (c *AppConfig) ConsAddress() (sdk.ConsAddress, error) {
	if len(c.ValidatorConsAddress) != 0 {
		return sdk.ConsAddressFromBech32(c.ValidatorConsAddress)
	}

	return nil, nil
}

// ValidateBasic performs basic validation of the app config.
func (c *AppConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if _, err := url.ParseRequestURI(c.OracleAddress); err != nil {
		return fmt.Errorf("oracle address must be valid: %w", err)
	}

	if c.ClientTimeout <= 0 {
		return fmt.Errorf("oracle client timeout must be greater than 0")
	}

	if c.MetricsEnabled {
		if c.PrometheusServerAddress == "" {
			return fmt.Errorf("must supply a non-empty prometheus server address if metrics are enabled")
		}
	}

	if _, err := c.ConsAddress(); err != nil {
		return err
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

	// get the prometheus server address
	if v := opts.Get(flagPrometheusServerAddress); v != nil {
		if cfg.PrometheusServerAddress, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	// get the validator consensus address
	if v := opts.Get(flagValidatorConsAddress); v != nil {
		if cfg.ValidatorConsAddress, err = cast.ToStringE(v); err != nil {
			return cfg, err
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return cfg, err
	}

	return cfg, err
}
