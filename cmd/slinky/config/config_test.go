package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/skip-mev/slinky/cmd/slinky/config"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	mmtypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestValidateBasic(t *testing.T) {
	tcs := []struct {
		name        string
		config      config.OracleConfig
		expectedErr bool
	}{
		{
			name: "good config",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:     oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:    oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:   oracleconfig.DefaultEnableCompression,
							ReadTimeout:         oracleconfig.DefaultReadTimeout,
							WriteTimeout:        oracleconfig.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: false,
		},
		{
			name: "bad config w/ bad provider",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "testa",
							ReadBufferSize:      oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:     oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:    oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:   oracleconfig.DefaultEnableCompression,
							ReadTimeout:         oracleconfig.DefaultReadTimeout,
							WriteTimeout:        oracleconfig.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name: "bad config w/ no max-price-age",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:     oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:    oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:   oracleconfig.DefaultEnableCompression,
							ReadTimeout:         oracleconfig.DefaultReadTimeout,
							WriteTimeout:        oracleconfig.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name:        "bad config with no update interval",
			config:      config.OracleConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with bad metrics",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:     oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:    oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:   oracleconfig.DefaultEnableCompression,
							ReadTimeout:         oracleconfig.DefaultReadTimeout,
							WriteTimeout:        oracleconfig.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Metrics: oracleconfig.MetricsConfig{
					Enabled: true,
				},
				Host: "localhost",
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name: "bad config with missing host",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:     oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:    oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:   oracleconfig.DefaultEnableCompression,
							ReadTimeout:         oracleconfig.DefaultReadTimeout,
							WriteTimeout:        oracleconfig.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Port: "8080",
			},
			expectedErr: true,
		},
		{
			name: "bad config with missing port",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:     oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:    oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:   oracleconfig.DefaultEnableCompression,
							ReadTimeout:         oracleconfig.DefaultReadTimeout,
							WriteTimeout:        oracleconfig.DefaultWriteTimeout,
						},
						Type: "price_provider",
					},
				},
				Host: "localhost",
			},
			expectedErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()

			// error should be nil if the test case expects the config to pass validation
			legacyCfg := tc.config.ToLegacy()
			if legacyErr, shouldBeNil := legacyCfg.ValidateBasic(), err == nil; (legacyErr == nil) != shouldBeNil {
				t.Errorf("expected legacy error to be nil, got %v", legacyErr)
			}

			// expect errors if necessary
			if !tc.expectedErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tc.expectedErr && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestReadOracleConfigWithOverrides(t *testing.T) {
	updateIntervalOverride := 100 * time.Second
	endpointOverride := oracleconfig.Endpoint{
		URL: "wss://test.com",
		Authentication: oracleconfig.Authentication{
			APIKey:       "test",
			APIKeyHeader: "test",
		},
	}
	prometheusServerOverride := "0.0.0.0:8081"

	expectedConfig := filterMarketMapProvidersFromOracleConfig(config.DefaultOracleConfig(), marketmap.Name)
	expectedConfig.UpdateInterval = updateIntervalOverride
	provider := expectedConfig.Providers[raydium.Name]
	provider.API.Endpoints = append(expectedConfig.Providers[raydium.Name].API.Endpoints, endpointOverride)
	expectedConfig.Providers[raydium.Name] = provider
	expectedConfig.Metrics.PrometheusServerAddress = prometheusServerOverride

	t.Run("overriding variables from environment", func(t *testing.T) {
		os.Setenv(config.SlinkyConfigEnvironmentPrefix+"_UPDATEINTERVAL", updateIntervalOverride.String())
		os.Setenv(config.SlinkyConfigEnvironmentPrefix+"_METRICS_PROMETHEUSSERVERADDRESS", prometheusServerOverride)
		os.Setenv(config.SlinkyConfigEnvironmentPrefix+"_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_URL", endpointOverride.URL)
		os.Setenv(config.SlinkyConfigEnvironmentPrefix+"_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_AUTHENTICATION_APIKEY", endpointOverride.Authentication.APIKey)
		os.Setenv(config.SlinkyConfigEnvironmentPrefix+"_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_AUTHENTICATION_APIKEYHEADER", endpointOverride.Authentication.APIKeyHeader)

		defer func() {
			os.Unsetenv(config.SlinkyConfigEnvironmentPrefix + "_UPDATEINTERVAL")
			os.Unsetenv(config.SlinkyConfigEnvironmentPrefix + "_METRICS_PROMETHEUSSERVERADDRESS")
			os.Unsetenv(config.SlinkyConfigEnvironmentPrefix + "_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_URL")
			os.Unsetenv(config.SlinkyConfigEnvironmentPrefix + "_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_AUTHENTICATION_APIKEY")
			os.Unsetenv(config.SlinkyConfigEnvironmentPrefix + "_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_AUTHENTICATION_APIKEYHEADER")
		}()

		cfg, err := config.ReadOracleConfigWithOverrides("", marketmap.Name)
		require.NoError(t, err)

		require.ElementsMatch(t, expectedConfig.ToLegacy().Providers, cfg.Providers)
		require.Equal(t, expectedConfig.UpdateInterval, cfg.UpdateInterval)
		require.Equal(t, expectedConfig.Metrics.PrometheusServerAddress, cfg.Metrics.PrometheusServerAddress)
	})

	t.Run("overriding variables via config", func(t *testing.T) {
		// create a temp file in the current directory
		tmpfile, err := os.CreateTemp("", "slinky-config-*.json")
		require.NoError(t, err)

		defer os.Remove(tmpfile.Name())

		overrides := fmt.Sprintf(`
		{
			"updateInterval": "%s",
			"metrics": {
				"prometheusServerAddress": "%s"
			},
			"providers": {
				"%s": {
					"api": {
						"endpoints": [
							{
								"url": "%s"
							},
							{
								"url": "%s",
								"authentication": {
									"apiKey": "%s",
									"apiKeyHeader": "%s"
								}
							}
						]
					}
				}
			}
		}
		`, updateIntervalOverride, prometheusServerOverride, raydium.Name, raydium.DefaultAPIConfig.Endpoints[0].URL, endpointOverride.URL, endpointOverride.Authentication.APIKey, endpointOverride.Authentication.APIKeyHeader)
		tmpfile.Write([]byte(overrides))

		cfg, err := config.ReadOracleConfigWithOverrides(tmpfile.Name(), marketmap.Name)
		require.NoError(t, err)

		require.ElementsMatch(t, expectedConfig.ToLegacy().Providers, cfg.Providers)
		require.Equal(t, expectedConfig.UpdateInterval, cfg.UpdateInterval)
		require.Equal(t, expectedConfig.Metrics.PrometheusServerAddress, cfg.Metrics.PrometheusServerAddress)
	})
}

func filterMarketMapProvidersFromOracleConfig(cfg config.OracleConfig, mmProvider string) config.OracleConfig {
	// filter out providers that are not in the market map
	for name, provider := range cfg.Providers {
		if provider.Type == mmtypes.ConfigType {
			if name != mmProvider {
				delete(cfg.Providers, name)
			}
		}
	}

	return cfg
}
