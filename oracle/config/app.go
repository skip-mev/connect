package config

import (
	"fmt"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	// DefaultConfigTemplate should be utilized in the app.toml file.
	DefaultConfigTemplate = `

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Oracle Address is the URL of the out of process oracle sidecar. This is used to
# connect to the oracle sidecar.
oracle_address = "{{ .Oracle.OracleAddress }}"

# Client Timeout is the time that the client is willing to wait for responses from 
# the oracle before timing out.
client_timeout = "{{ .Oracle.ClientTimeout }}"

# MetricsEnabled determines whether oracle metrics are enabled. Specifically
# this enables intsrumentation of the oracle client and the interaction between
# the oracle and the app.
metrics_enabled = "{{ .Oracle.MetricsEnabled }}"

# PrometheusServerAddress is the address of the prometheus server that the oracle 
# will expose metrics to
prometheus_server_address = "{{ .Oracle.PrometheusServerAddress }}"

# ValidatorConsAddress is the validator's consensus address. Validator's must register their
# consensus address in order to enable app side metrics.
validator_cons_address = "{{ .Oracle.ValidatorConsAddress }}"
`
)

const (
	// Flags utilized to retrieve the config from the config file.
	flagOracleAddress = "oracle.oracle_address"

	// Flags utilized to retrieve the config from the config file.
	flagClientTimeout = "oracle.client_timeout"

	// Flags utilized to retrieve the config from the config file.
	flagMetricsEnabled = "oracle.metrics_enabled"

	// Flags utilized to retrieve the config from the config file.
	flagPrometheusServerAddress = "oracle.prometheus_server_address"

	// Flags utilized to retrieve the config from the config file.
	flagValidatorConsAddress = "oracle.validator_cons_address"
)

// AppConfig contains the application side oracle configurations that must
// be set in the app.toml file.
type AppConfig struct {
	// OracleAddress is the URL of the out of process oracle sidecar. This is used to
	// connect to the oracle sidecar.
	OracleAddress string `mapstructure:"oracle_address" toml:"oracle_address"`

	// ClientTimeout is the time that the client is willing to wait for responses from the oracle before timing out.
	ClientTimeout time.Duration `mapstructure:"client_timeout" toml:"client_timeout"`

	// MetricsEnabled is the address of the prometheus server that the oracle will expose metrics to
	MetricsEnabled bool `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`

	// ValidatorConsAddress is the validator's consensus address. Validator's must register their
	// consensus address in order to enable app side metrics.
	ValidatorConsAddress string `mapstructure:"validator_cons_address" toml:"validator_cons_address"`
}

// ConsAddress returns the validator's consensus address
func (c *AppConfig) ConsAddress() (sdk.ConsAddress, error) {
	if len(c.ValidatorConsAddress) != 0 {
		return sdk.ConsAddressFromBech32(c.ValidatorConsAddress)
	}

	return nil, nil
}

// ValidateBasic performs basic validation of the config.
func (c *AppConfig) ValidateBasic() error {
	if len(c.OracleAddress) == 0 {
		return fmt.Errorf("oracle address cannot be empty; please set oracle address in config")
	}

	if c.ClientTimeout <= 0 {
		return fmt.Errorf("oracle client timeout must be greater than 0")
	}

	if c.MetricsEnabled && len(c.PrometheusServerAddress) == 0 {
		return fmt.Errorf("must supply a prometheus server address if metrics are enabled")
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

	// Check required fields
	fields := []string{"oracle_address", "client_timeout"}
	for _, field := range fields {
		if !viper.IsSet(field) {
			return config, fmt.Errorf("required oracle field is missing in config")
		}
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
