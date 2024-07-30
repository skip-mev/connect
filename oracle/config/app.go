package config

import (
	"fmt"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

var (
	DefaultOracleEnabled  = false
	DefaultOracleAddress  = "localhost:8080"
	DefaultClientTimeout  = 3 * time.Second
	DefaultMetricsEnabled = false
	DefaultPriceTTL       = 10 * time.Second
	DefaultInterval       = 1500 * time.Millisecond

	MaxInterval = 1 * time.Minute
	MaxPriceTTL = 1 * time.Minute
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
# the oracle before timing out. The recommended timeout is 3 seconds (3000ms).
client_timeout = "{{ .Oracle.ClientTimeout }}"

# MetricsEnabled determines whether oracle metrics are enabled. Specifically
# this enables instrumentation of the oracle client and the interaction between
# the oracle and the app.
metrics_enabled = "{{ .Oracle.MetricsEnabled }}"

# PriceTTL is the maximum age of the latest price response before it is considered stale. 
# The recommended max age is 10 seconds (10s). If this is greater than 1 minute (1m), the app
# will not start.
price_ttl = "{{ .Oracle.PriceTTL }}"

# Interval is the time between each price update request. The recommended interval
# is the block time of the chain. Otherwise, 1.5 seconds (1500ms) is a good default. If this
# is greater than 1 minute (1m), the app will not start.
interval = "{{ .Oracle.Interval }}"
`
)

// NewDefaultAppConfig returns a default application side oracle configuration.
func NewDefaultAppConfig() AppConfig {
	return AppConfig{
		Enabled:        DefaultOracleEnabled,
		OracleAddress:  DefaultOracleAddress,
		ClientTimeout:  DefaultClientTimeout,
		MetricsEnabled: DefaultMetricsEnabled,
		PriceTTL:       DefaultPriceTTL,
		Interval:       DefaultInterval,
	}
}

const (
	flagEnabled                 = "oracle.enabled"
	flagOracleAddress           = "oracle.oracle_address"
	flagClientTimeout           = "oracle.client_timeout"
	flagMetricsEnabled          = "oracle.metrics_enabled"
	flagPrometheusServerAddress = "oracle.prometheus_server_address"
	flagPriceTTL                = "oracle.price_ttl"
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

	// PriceTTL is the maximum age of the latest price response before it is considered
	// stale.
	PriceTTL time.Duration `mapstructure:"price_ttl" toml:"price_ttl"`

	// Interval is the time between each price update request.
	Interval time.Duration `mapstructure:"interval" toml:"interval"`
}

// ValidateBasic performs basic validation of the app config.
func (c *AppConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if len(c.OracleAddress) == 0 {
		return fmt.Errorf("poorly formatted app.toml (oracle subsection): oracle address must not be empty")
	}

	if c.ClientTimeout <= 0 {
		return fmt.Errorf("poorly formatted app.toml (oracle subsection): oracle client timeout must be greater than 0")
	}

	if c.PriceTTL <= 0 || c.PriceTTL > MaxPriceTTL {
		return fmt.Errorf("poorly formatted app.toml (oracle subsection): oracle price time to live (price_ttl) must be between 0 and %s", MaxPriceTTL)
	}

	if c.Interval <= 0 || c.Interval > MaxInterval {
		return fmt.Errorf("poorly formatted app.toml (oracle subsection): oracle interval must be between 0 and %s", MaxInterval)
	}

	if c.Interval >= c.PriceTTL {
		return fmt.Errorf("poorly formatted app.toml (oracle subsection): oracle interval must be strictly less than max age")
	}

	return nil
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (AppConfig, error) {
	var (
		cfg = NewDefaultAppConfig()
		err error
	)

	// determine if the oracle is enabled
	if v := opts.Get(flagEnabled); v != nil {
		if cfg.Enabled, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	if !cfg.Enabled {
		return cfg, nil
	}

	// get the oracle address
	if v := opts.Get(flagOracleAddress); v != nil {
		address, err := cast.ToStringE(v)
		if err != nil {
			return cfg, fmt.Errorf("oracle address must be a non-empty string")
		}

		// only update the address if it is non-empty
		if len(address) > 0 {
			cfg.OracleAddress = address
		}
	}

	// get the client timeout
	if v := opts.Get(flagClientTimeout); v != nil {
		clientTimeout, err := cast.ToDurationE(v)
		if err != nil {
			return cfg, fmt.Errorf("client timeout must be a positive duration")
		}

		// only update the client timeout if it is positive
		if clientTimeout > 0 {
			cfg.ClientTimeout = clientTimeout
		}
	}

	// get the metrics enabled
	if v := opts.Get(flagMetricsEnabled); v != nil {
		if cfg.MetricsEnabled, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}

	// get the price ttl
	if v := opts.Get(flagPriceTTL); v != nil {
		priceTTL, err := cast.ToDurationE(v)
		if err != nil {
			return cfg, fmt.Errorf("price ttl must be a positive duration")
		}

		// only update the price ttl if it is positive
		if priceTTL > 0 {
			cfg.PriceTTL = priceTTL
		}
	}

	// get the interval
	if v := opts.Get(flagInterval); v != nil {
		interval, err := cast.ToDurationE(v)
		if err != nil {
			return cfg, fmt.Errorf("interval must be a positive duration")
		}

		// only update the interval if it is positive
		if interval > 0 {
			cfg.Interval = interval
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return cfg, err
	}

	return cfg, err
}

// String implements the stringer interface for the AppConfig.
func (c AppConfig) String() string {
	return fmt.Sprintf(`Oracle Config:
  Enabled: %v
  Oracle Address: %s
  Client Timeout: %s
  Metrics Enabled: %v
  Price TTL: %s
  Interval: %s`,
		c.Enabled, c.OracleAddress, c.ClientTimeout, c.MetricsEnabled, c.PriceTTL, c.Interval)
}
