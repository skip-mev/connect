package main

import (
	"fmt"
	"reflect"
	"time"

	binanceapi "github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	krakenapi "github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/volatile"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	"github.com/skip-mev/slinky/providers/websockets/bitstamp"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	"github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/gate"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	"github.com/skip-mev/slinky/providers/websockets/kraken"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
	"github.com/spf13/viper"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
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
	// jsonFieldDelimiter is the delimiter used to separate fields in the JSON output.
	jsonFieldDelimiter = "."
)

// DefaultOracleConfig returns the default configuration for the slinky oracle
var DefaultOracleConfig = OracleConfig{
	UpdateInterval: DefaultUpdateInterval,
	MaxPriceAge:    DefaultMaxPriceAge,
	Providers: map[string]config.ProviderConfig{
		coinbaseapi.Name: config.ProviderConfig{
			Name: coinbaseapi.Name,
			API: coinbaseapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		binanceapi.Name: config.ProviderConfig{
			Name: binanceapi.Name,
			API: binanceapi.DefaultNonUSAPIConfig,
			Type: types.ConfigType,
		},
		raydium.Name: config.ProviderConfig{
			Name: raydium.Name,
			API: raydium.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		uniswapv3.ProviderNames[constants.ETHEREUM]: config.ProviderConfig{
			Name: uniswapv3.ProviderNames[constants.ETHEREUM],
			API: uniswapv3.DefaultETHAPIConfig,
			Type: types.ConfigType,
		},
		krakenapi.Name: config.ProviderConfig{
			Name: krakenapi.Name,
			API: krakenapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		marketmap.Name: config.ProviderConfig{
			Name: marketmap.Name,
			API: marketmap.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		volatile.Name: config.ProviderConfig{
			Name: volatile.Name,
			API: volatile.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		bitfinex.Name: config.ProviderConfig{
			Name: bitfinex.Name,
			WebSocket: bitfinex.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		bitstamp.Name: config.ProviderConfig{
			Name: bitstamp.Name,
			WebSocket: bitstamp.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		bybit.Name: config.ProviderConfig{
			Name: bybit.Name,
			WebSocket: bybit.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		coinbase.Name: config.ProviderConfig{
			Name: coinbase.Name,
			WebSocket: coinbase.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		cryptodotcom.Name: config.ProviderConfig{
			Name: cryptodotcom.Name,
			WebSocket: cryptodotcom.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		gate.Name: config.ProviderConfig{
			Name: gate.Name,
			WebSocket: gate.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		huobi.Name: config.ProviderConfig{
			Name: huobi.Name,
			WebSocket: huobi.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		kraken.Name: config.ProviderConfig{
			Name: kraken.Name,
			WebSocket: kraken.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		kucoin.Name: config.ProviderConfig{
			Name: kucoin.Name,
			WebSocket: kucoin.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		mexc.Name: config.ProviderConfig{
			Name: mexc.Name,
			WebSocket: mexc.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
		okx.Name: config.ProviderConfig{
			Name: okx.Name,
			WebSocket: okx.DefaultWebSocketConfig,
			Type: types.ConfigType,
		},
	},
	Metrics: config.MetricsConfig{
		PrometheusServerAddress: DefaultPrometheusServerAddress,
		Enabled:                DefaultMetricsEnabled,
	},
	Host: DefaultHost,
	Port: DefaultPort,
}

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
			return fmt.Errorf("provider %s is not formatted correctly: %w", p.Name, err)
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
	setViperDefaultsForDataStructure("", DefaultOracleConfig)
}

func setViperDefaultsForDataStructure(keyPrefix string, config interface{}) {
	val := reflect.ValueOf(config)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		jsonFieldTag := typ.Field(i).Tag.Get("json")

		// the fully-qualified key for this field
		fullKey := keyPrefix + jsonFieldTag

		if field.Kind() == reflect.Struct {
			setViperDefaultsForDataStructure(fullKey + jsonFieldDelimiter, field.Interface())
		} else if field.Kind() == reflect.Map {
			// set viper defaults for map
			for _, key := range field.MapKeys() {
				setViperDefaultsForDataStructure(
					fullKey+jsonFieldDelimiter+key.String() + jsonFieldDelimiter, 
					field.MapIndex(key).Interface(),
				)
			}
		} else {
			viper.SetDefault(fullKey, field.Interface())
		}
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
	SetDefaults()

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
