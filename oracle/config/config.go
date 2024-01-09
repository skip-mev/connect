package config

import (
	"fmt"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
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
# Oracle path is the path for the config file for the oracle.
oracle_path = "{{ .Oracle.OraclePath }}"

# Metrics path is the path for the config file for the metrics.
metrics_path = "{{ .Oracle.MetricsPath }}"
`
)

const (
	// Flags utilized to retrieve the config from baseapp options.
	flagOraclePath  = "oracle.oracle_path"
	flagMetricsPath = "oracle.metrics_path"
)

// Config is the over-arching config for the oracle sidecar and instrumentation. It expects
// a config file with the path to the oracle and metrics config files.
type Config struct {
	// OraclePath is the path for the config file for the oracle.
	OraclePath string `mapstructure:"oracle_path" toml:"oracle_path"`

	// MetricsPath is the path for the config file for the metrics.
	MetricsPath string `mapstructure:"metrics_path" toml:"metrics_path"`
}

// ValidateBasic performs basic validation of the config.
func (c *Config) ValidateBasic() error {
	if len(c.OraclePath) == 0 {
		return fmt.Errorf("oracle path cannot be empty")
	}

	if len(c.MetricsPath) == 0 {
		return fmt.Errorf("metrics path cannot be empty")
	}

	return nil
}

// ReadConfigFromFile reads a config from a file and returns the config.
func ReadConfigFromFile(path string) (Config, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	// Check required fields
	requiredFields := []string{"oracle_path", "metrics_path"}
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			return Config{}, fmt.Errorf("required field %s is missing in config", field)
		}
	}

	// unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	if err := config.ValidateBasic(); err != nil {
		return Config{}, err
	}

	return config, nil
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (Config, error) {
	var (
		cfg Config
		err error
	)

	// get the path to the oracle config
	if v := opts.Get(flagOraclePath); v != nil {
		if cfg.OraclePath, err = cast.ToStringE(v); err != nil {
			return Config{}, err
		}
	}

	// get the path to the metrics config
	if v := opts.Get(flagMetricsPath); v != nil {
		if cfg.MetricsPath, err = cast.ToStringE(v); err != nil {
			return Config{}, err
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return Config{}, err
	}

	return cfg, err
}
