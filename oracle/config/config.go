package config

import (
	"fmt"
	"time"

	"cosmossdk.io/log"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/coinbase"
	"github.com/skip-mev/slinky/providers/coingecko"
	"github.com/skip-mev/slinky/providers/coinmarketcap"
	"github.com/skip-mev/slinky/providers/evm/erc4626"
	"github.com/skip-mev/slinky/providers/evm/erc4626sharepriceoracle"
	"github.com/skip-mev/slinky/providers/mock"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	oracle_metrics "github.com/skip-mev/slinky/oracle/metrics"
	service_metrics "github.com/skip-mev/slinky/service/metrics"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	flagOracleInProcess      = "oracle.in_process"
	flagOracleRemoteAddress  = "oracle.remote_address"
	flagOracleUpdateInterval = "oracle.update_interval"
	flagOracleProviders      = "oracle.providers"
	flagOracleCurrencyPairs  = "oracle.currency_pairs"
	flagTimeout              = "oracle.timeout"
	flagMetricsServerAddress = "oracle.metrics.prometheus_server_address"
	flagOracleMetricsEnabled = "oracle.metrics.oracle_metrics.enabled"
	flagAppMetricsEnabled    = "oracle.metrics.app_metrics.enabled"
	flagAppMetricsValidator  = "oracle.metrics.app_metrics.validator_cons_address"
)

// Config is the over-arching config for the oracle sidecar and instrumentation. It configures the
// in-process oracle server, out-of-process oracle client and the instrumentation for both.
type Config struct {
	Oracle  `mapstructure:"oracle" toml:"oracle"`
	Metrics `mapstructure:"metrics" toml:"metrics"`
}

// Metrics is the base metrics for both out-of-process and in-process oracles.
type Metrics struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`
	OracleMetrics           oracle_metrics.Config
	AppMetrics              service_metrics.Config
}

// Oracle is the base config for both out-of-process and in-process oracles. If the oracle is to be configured out-of-process in base-app, a
// grpc-client of the grpc-server running at RemoteAddress is instantiated, otherwise, an in-process LocalClient oracle is instantiated.
type Oracle struct {
	// InProcess specifies whether the oracle configured, is currently running as a remote grpc-server, or will be run in process
	InProcess bool `mapstructure:"in_process" toml:"in_process"`

	// Timeout is the time that the client is willing to wait for responses from the oracle
	Timeout time.Duration `mapstructure:"timeout" toml:"timeout"`

	// RemoteAddress is the address of the remote oracle server (if it is running out-of-process)
	RemoteAddress string `mapstructure:"remote_address" toml:"remote_address"`

	// UpdateInterval is the interval at which the oracle will fetch prices from providers
	UpdateInterval time.Duration `mapstructure:"update_interval" toml:"update_interval"`

	// Providers is the set of providers that the oracle will fetch prices from.
	Providers []types.ProviderConfig `mapstructure:"providers" toml:"providers"`

	// CurrencyPairs is the set of currency pairs that the oracle will fetch prices for
	CurrencyPairs []oracletypes.CurrencyPair `mapstructure:"currency_pairs" toml:"currency_pairs"`
}

// ReadConfigFromFile reads a config from a file and returns the config.
func ReadConfigFromFile(path string) (*Config, error) {
	// read in config file
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func providerFromProviderConfig(cfg types.ProviderConfig, cps []oracletypes.CurrencyPair, l log.Logger) (types.Provider, error) {
	switch cfg.Name {
	case "coingecko":
		return coingecko.NewProvider(l, cps, cfg.Apikey), nil
	case "coinbase":
		return coinbase.NewProvider(l, cps), nil
	case "coinmarketcap":
		return coinmarketcap.NewProvider(l, cps, cfg.Apikey, cfg.TokenNameToMetadata), nil
	case "erc4626":
		provider, err := erc4626.NewProvider(l, cps, cfg.Apikey, cfg.TokenNameToMetadata)
		if err != nil {
			return nil, err
		}

		return provider, nil
	case "erc4626-share-price-oracle":
		provider, err := erc4626sharepriceoracle.NewProvider(l, cps, cfg.Apikey, cfg.TokenNameToMetadata)
		if err != nil {
			return nil, err
		}

		return provider, nil
	case "timeout-mock-provider":
		// This will timeout after the configured timeout + 1 second
		return mock.NewTimeoutMockProvider(cfg.ProviderTimeout + time.Second), nil
	case "failing-mock-provider":
		// This will always panic whenever GetPrices is called
		return mock.NewFailingMockProvider(), nil
	case "static-mock-provider":
		// This will return mock prices (randomly generated) for the configured currency pairs
		if len(cfg.TokenNameToMetadata) != 0 {
			return mock.NewStaticMockProviderFromConfig(cfg), nil
		}

		return mock.NewStaticMockProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
}

// Providers returns the set of providers that this config determines from the provider configs
func (c *Config) GetProviders(l log.Logger) ([]types.Provider, error) {
	providers := make([]types.Provider, len(c.Providers))
	var err error

	for i, p := range c.Providers {
		providers[i], err = providerFromProviderConfig(p, c.CurrencyPairs, l)
		if err != nil {
			return nil, err
		}
	}

	return providers, nil
}

// ReadConfigFromAppOpts reads the config parameters from the AppOptions and returns the config.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (cfg *Config, err error) {
	cfg = &Config{}

	// get the in-process flag
	if v := opts.Get(flagOracleInProcess); v != nil {
		if cfg.InProcess, err = cast.ToBoolE(v); err != nil {
			return nil, err
		}
	}

	// get timeout
	if v := opts.Get(flagTimeout); v != nil {
		if cfg.Timeout, err = cast.ToDurationE(v); err != nil {
			return nil, err
		}
	}

	// get the remote address
	if !cfg.InProcess {
		if v := opts.Get(flagOracleRemoteAddress); v != nil {
			if cfg.RemoteAddress, err = cast.ToStringE(v); err != nil {
				return nil, err
			}
		}
	}

	// get the update interval
	if v := opts.Get(flagOracleUpdateInterval); v != nil {
		if cfg.UpdateInterval, err = cast.ToDurationE(v); err != nil {
			return nil, err
		}
	}

	// get the providers
	if v := opts.Get(flagOracleProviders); v != nil {
		iFaces, err := cast.ToSliceE(v)
		if err != nil {
			return nil, err
		}

		// iterate through iterfaces and add to config
		for _, iFace := range iFaces {
			if providerCfg, err := providerConfigFromToml(iFace); err == nil {
				cfg.Providers = append(cfg.Providers, providerCfg)
			}
		}
	}

	// get the currency pairs
	if v := opts.Get(flagOracleCurrencyPairs); v != nil {
		iFaces, err := cast.ToSliceE(v)
		if err != nil {
			return nil, err
		}

		// iterate through iterfaces and add to config
		for _, iFace := range iFaces {
			if currencyPair, err := currencyPairConfigFromToml(iFace); err == nil {
				cfg.CurrencyPairs = append(cfg.CurrencyPairs, currencyPair)
			}
		}
	}

	// get the metrics config
	metrics, err := metricsConfigFromToml(opts)
	if err != nil {
		return nil, err
	}
	cfg.Metrics = metrics

	return cfg, err
}

func metricsConfigFromToml(opts servertypes.AppOptions) (Metrics, error) {
	metrics := Metrics{}
	var err error

	// get prometheus server address
	if v := opts.Get(flagMetricsServerAddress); v != nil {
		if metrics.PrometheusServerAddress, err = cast.ToStringE(v); err != nil {
			return metrics, err
		}
	}

	// get the metrics enabled flag
	if v := opts.Get(flagOracleMetricsEnabled); v != nil {
		if metrics.OracleMetrics.Enabled, err = cast.ToBoolE(v); err != nil {
			return metrics, err
		}
	}

	// get the app metrics enabled flag
	if v := opts.Get(flagAppMetricsEnabled); v != nil {
		if metrics.AppMetrics.Enabled, err = cast.ToBoolE(v); err != nil {
			return metrics, err
		}
	}

	// get the app metrics validator address
	if v := opts.Get(flagAppMetricsValidator); v != nil {
		if metrics.AppMetrics.ValidatorConsAddress, err = cast.ToStringE(v); err != nil {
			return metrics, err
		}
	}

	return metrics, nil
}

func providerConfigFromToml(iface interface{}) (types.ProviderConfig, error) {
	providerCfg := types.ProviderConfig{}

	// convert interface to map
	iFaceMap, ok := iface.(map[string]interface{})
	if !ok {
		return providerCfg, fmt.Errorf("failed to convert interface to map")
	}

	// get the name
	if v, ok := iFaceMap["name"]; ok {
		if providerCfg.Name, ok = v.(string); !ok {
			return providerCfg, fmt.Errorf("failed to convert name to string")
		}
	}

	// get the apikey
	if v, ok := iFaceMap["apikey"]; ok {
		if providerCfg.Apikey, ok = v.(string); !ok {
			return providerCfg, fmt.Errorf("failed to convert apikey to string")
		}
	}

	// get the token name to metadata map
	if v, ok := iFaceMap["token_name_to_metadata"]; ok {
		if tokenNameToMetadata, ok := v.(map[string]interface{}); ok {
			providerCfg.TokenNameToMetadata = make(map[string]types.TokenMetadata)
			for k, v := range tokenNameToMetadata {
				if metadata, ok := v.(types.TokenMetadata); ok {
					providerCfg.TokenNameToMetadata[k] = metadata
				}
			}
		}
	}

	// get the provider timeout
	if v, ok := iFaceMap["provider_timeout"]; ok {
		var err error
		if providerCfg.ProviderTimeout, err = cast.ToDurationE(v); err != nil {
			return providerCfg, fmt.Errorf("failed to convert provider timeout to duration")
		}
	}

	return providerCfg, nil
}

func currencyPairConfigFromToml(iface interface{}) (oracletypes.CurrencyPair, error) {
	currencyPair := oracletypes.CurrencyPair{}

	// convert interface to map
	iFaceMap, ok := iface.(map[string]interface{})
	if !ok {
		return currencyPair, fmt.Errorf("failed to convert interface to map")
	}

	// get the base currency
	if v, ok := iFaceMap["base"]; ok {
		if currencyPair.Base, ok = v.(string); !ok {
			return currencyPair, fmt.Errorf("failed to convert base currency to string")
		}
	} else if v, ok := iFaceMap["Base"]; ok {
		if currencyPair.Base, ok = v.(string); !ok {
			return currencyPair, fmt.Errorf("failed to convert base currency to string")
		}
	}

	// get the quote currency
	if v, ok := iFaceMap["quote"]; ok {
		if currencyPair.Quote, ok = v.(string); !ok {
			return currencyPair, fmt.Errorf("failed to convert quote currency to string")
		}
	} else if v, ok := iFaceMap["Quote"]; ok {
		if currencyPair.Quote, ok = v.(string); !ok {
			return currencyPair, fmt.Errorf("failed to convert quote currency to string")
		}
	}

	return currencyPair, nil
}
