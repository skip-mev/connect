package config

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/skip-mev/connect/v2/cmd/constants"
	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	mmtypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
)

const (
	// DefaultUpdateInterval is the default value for how frequently connect updates aggregate price responses.
	DefaultUpdateInterval = 250000000
	// DefaultMaxPriceAge is the default value for the oldest price considered in an aggregate price response by connect.
	DefaultMaxPriceAge = 120000000000
	// DefaultPrometheusServerAddress is the default value for the prometheus server address in connect.
	DefaultPrometheusServerAddress = "0.0.0.0:8002"
	// DefaultMetricsEnabled is the default value for enabling prometheus metrics in connect.
	DefaultMetricsEnabled = true
	// DefaultTelemetryDisabled is the default value for disabling telemetry.
	DefaultTelemetryDisabled = false
	// DefaultHost is the default for the connect oracle server host.
	DefaultHost = "0.0.0.0"
	// DefaultPort is the default for the connect oracle server port.
	DefaultPort = "8080"
	// jsonFieldDelimiter is the delimiter used to separate fields in the JSON output.
	jsonFieldDelimiter = "."
	// ConnectConfigEnvironmentPrefix is the prefix for environment variables that override the connect config.
	ConnectConfigEnvironmentPrefix = "CONNECT_CONFIG"
	// TelemetryPushAddress is the value for the publication endpoint.
	TelemetryPushAddress = "connect-statsd-data.dev.skip.money:9125"
)

// DefaultOracleConfig returns the default configuration for the connect oracle.
func DefaultOracleConfig() config.OracleConfig {
	cfg := config.OracleConfig{
		UpdateInterval: DefaultUpdateInterval,
		MaxPriceAge:    DefaultMaxPriceAge,
		Metrics: config.MetricsConfig{
			PrometheusServerAddress: DefaultPrometheusServerAddress,
			Enabled:                 DefaultMetricsEnabled,
			Telemetry: config.TelemetryConfig{
				Disabled:    DefaultTelemetryDisabled,
				PushAddress: TelemetryPushAddress,
			},
		},
		Providers: make(map[string]config.ProviderConfig),
		Host:      DefaultHost,
		Port:      DefaultPort,
	}

	for _, provider := range append(constants.Providers, constants.AlternativeMarketMapProviders...) {
		cfg.Providers[provider.Name] = provider
	}

	if err := cfg.ValidateBasic(); err != nil {
		panic(fmt.Sprintf("default oracle config is invalid: %s", err))
	}

	return cfg
}

func SetDefaults() {
	setViperDefaultsForDataStructure("", DefaultOracleConfig())
}

func setViperDefaultsForDataStructure(keyPrefix string, config interface{}) {
	val := reflect.ValueOf(config)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		jsonFieldTag := typ.Field(i).Tag.Get("json")

		// the fully-qualified key for this field
		fullKey := keyPrefix + jsonFieldTag

		switch field.Kind() {
		case reflect.Struct:
			// set viper defaults for struct via recursion
			setViperDefaultsForDataStructure(fullKey+jsonFieldDelimiter, field.Interface())
		case reflect.Map:
			// set viper defaults for map
			for _, key := range field.MapKeys() {
				setViperDefaultsForDataStructure(
					fullKey+jsonFieldDelimiter+key.String()+jsonFieldDelimiter,
					field.MapIndex(key).Interface(),
				)
			}
		default:
			viper.SetDefault(fullKey, field.Interface())
		}
	}

	// set the environment prefix
	viper.SetEnvPrefix(ConnectConfigEnvironmentPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// ReadOracleConfigWithOverrides reads a config from a file and returns the config.
func ReadOracleConfigWithOverrides(path string, marketMapProvider string) (config.OracleConfig, error) {
	// if the path is non-nil read data from a file\
	SetDefaults()
	if path != "" {
		viper.SetConfigFile(path)
		viper.SetConfigType("json")

		if err := viper.ReadInConfig(); err != nil {
			return config.OracleConfig{}, err
		}
	}

	cfg, err := oracleConfigFromViper()
	if err != nil {
		return config.OracleConfig{}, err
	}

	// filter the market-map providers
	if _, ok := constants.MarketMapProviderNames[marketMapProvider]; !ok {
		return config.OracleConfig{}, fmt.Errorf("market map provider %s not found", marketMapProvider)
	}

	// filter the unused market-map providers
	for name, provider := range cfg.Providers {
		if provider.Type == mmtypes.ConfigType {
			if name != marketMapProvider {
				delete(cfg.Providers, name)
			}
		}
	}

	return cfg, cfg.ValidateBasic()
}

// oracleConfigFromViper unmarshals an oracle config from viper, validates it, and returns it.
func oracleConfigFromViper() (config.OracleConfig, error) {
	var cfg config.OracleConfig
	unmarshalMetadata := mapstructure.Metadata{}
	if err := viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.ErrorUnused = true
		c.Metadata = &unmarshalMetadata
	}); err != nil {
		return config.OracleConfig{}, err
	}

	// for each api-provider, we'll have to manually fill the endpoints
	for _, provider := range cfg.Providers {
		// if a provider was not unmarshaled correctly, surface that error
		if provider.Name == "" {
			if len(unmarshalMetadata.Unset) > 0 {
				return config.OracleConfig{}, fmt.Errorf("overridden key %s does not correspond to a known provider", unmarshalMetadata.Unset[0])
			}
		}

		// Update API endpoints
		for i, endpoint := range provider.API.Endpoints {
			provider.API.Endpoints[i], _ = updateEndpointFromEnvironment(endpoint, provider.Name, i, "api")
		}

		firstEndpointFromViperIndex := len(provider.API.Endpoints)
		for found := true; found; firstEndpointFromViperIndex++ {
			var endpoint config.Endpoint
			endpoint, found = updateEndpointFromEnvironment(config.Endpoint{}, provider.Name, firstEndpointFromViperIndex, "api")
			if found {
				provider.API.Endpoints = append(provider.API.Endpoints, endpoint)
			}
		}

		// Update WebSocket endpoints
		for i, endpoint := range provider.WebSocket.Endpoints {
			provider.WebSocket.Endpoints[i], _ = updateEndpointFromEnvironment(endpoint, provider.Name, i, "webSocket")
		}

		firstEndpointFromViperIndex = len(provider.WebSocket.Endpoints)
		for found := true; found; firstEndpointFromViperIndex++ {
			var endpoint config.Endpoint
			endpoint, found = updateEndpointFromEnvironment(config.Endpoint{}, provider.Name, firstEndpointFromViperIndex, "webSocket")
			if found {
				provider.WebSocket.Endpoints = append(provider.WebSocket.Endpoints, endpoint)
			}
		}

		// update the provider with the updated endpoints
		cfg.Providers[provider.Name] = provider
	}

	if err := cfg.ValidateBasic(); err != nil {
		return config.OracleConfig{}, err
	}

	return cfg, nil
}

// updateEndpointFromEnvironment returns an updated endpoint with the values from the environment. If viper is not aware of
// any overrides variables for the endpoint, the original endpoint is returned with false.
func updateEndpointFromEnvironment(endpoint config.Endpoint, providerName string, idx int, configType string) (config.Endpoint, bool) {
	// check if an environment variable exists for this endpoint
	endpointURL := viper.Get(fmt.Sprintf("providers.%s.%s.endpoints.%d.url", providerName, configType, idx))
	endpointAPIKey := viper.Get(fmt.Sprintf("providers.%s.%s.endpoints.%d.authentication.apiKey", providerName, configType, idx))
	endpointAPIKeyHeader := viper.Get(fmt.Sprintf("providers.%s.%s.endpoints.%d.authentication.apiKeyHeader", providerName, configType, idx))

	// if the environment variable exists, set the endpoint to the value of the environment variable
	if endpointURL != nil {
		endpoint.URL = endpointURL.(string)
	}

	if endpointAPIKey != nil {
		endpoint.Authentication.APIKey = endpointAPIKey.(string)
	}

	if endpointAPIKeyHeader != nil {
		endpoint.Authentication.APIKeyHeader = endpointAPIKeyHeader.(string)
	}

	return endpoint, endpointURL != nil || endpointAPIKey != nil || endpointAPIKeyHeader != nil
}

func GetNodeEndpointFromConfig(cfg config.OracleConfig) (config.Endpoint, error) {
	for _, provider := range cfg.Providers {
		if provider.Type == mmtypes.ConfigType {
			isAlternativeMMProvider := slices.IndexFunc(constants.AlternativeMarketMapProviders, func(c config.ProviderConfig) bool {
				return c.Name == provider.Name
			}) >= 0

			if isAlternativeMMProvider {
				if provider.Name == dydx.SwitchOverAPIHandlerName {
					// grpc endpoint is always the 2nd entry in the dydx
					// migration provider
					return provider.API.Endpoints[1], nil
				}
			} else {
				// normal mm providers are grpc endpoints
				return provider.API.Endpoints[0], nil
			}
		}
	}

	return config.Endpoint{}, fmt.Errorf("could not find marketmap endpoint")
}
