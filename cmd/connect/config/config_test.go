package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cmdconfig "github.com/skip-mev/connect/v2/cmd/connect/config"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	"github.com/skip-mev/connect/v2/providers/websockets/coinbase"
	mmtypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
)

func TestValidateBasic(t *testing.T) {
	tcs := []struct {
		name        string
		config      oracleconfig.OracleConfig
		expectedErr bool
	}{
		{
			name: "good config",
			config: oracleconfig.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []oracleconfig.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:                     "test",
							ReadBufferSize:           oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:          oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:         oracleconfig.DefaultHandshakeTimeout,
							EnableCompression:        oracleconfig.DefaultEnableCompression,
							ReadTimeout:              oracleconfig.DefaultReadTimeout,
							WriteTimeout:             oracleconfig.DefaultWriteTimeout,
							MaxSubscriptionsPerBatch: oracleconfig.DefaultMaxSubscriptionsPerBatch,
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
			config: oracleconfig.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []oracleconfig.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "testa",
							ReadBufferSize:    oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:   oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:  oracleconfig.DefaultHandshakeTimeout,
							EnableCompression: oracleconfig.DefaultEnableCompression,
							ReadTimeout:       oracleconfig.DefaultReadTimeout,
							WriteTimeout:      oracleconfig.DefaultWriteTimeout,
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
			config: oracleconfig.OracleConfig{
				UpdateInterval: time.Second,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []oracleconfig.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:   oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:  oracleconfig.DefaultHandshakeTimeout,
							EnableCompression: oracleconfig.DefaultEnableCompression,
							ReadTimeout:       oracleconfig.DefaultReadTimeout,
							WriteTimeout:      oracleconfig.DefaultWriteTimeout,
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
			config:      oracleconfig.OracleConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with bad metrics",
			config: oracleconfig.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []oracleconfig.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:   oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:  oracleconfig.DefaultHandshakeTimeout,
							EnableCompression: oracleconfig.DefaultEnableCompression,
							ReadTimeout:       oracleconfig.DefaultReadTimeout,
							WriteTimeout:      oracleconfig.DefaultWriteTimeout,
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
			config: oracleconfig.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []oracleconfig.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:   oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:  oracleconfig.DefaultHandshakeTimeout,
							EnableCompression: oracleconfig.DefaultEnableCompression,
							ReadTimeout:       oracleconfig.DefaultReadTimeout,
							WriteTimeout:      oracleconfig.DefaultWriteTimeout,
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
			config: oracleconfig.OracleConfig{
				UpdateInterval: time.Second,
				MaxPriceAge:    time.Minute,
				Providers: map[string]oracleconfig.ProviderConfig{
					"test": {
						Name: "test",
						WebSocket: oracleconfig.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							Endpoints: []oracleconfig.Endpoint{
								{
									URL: "wss://test.com",
								},
							},
							Name:              "test",
							ReadBufferSize:    oracleconfig.DefaultReadBufferSize,
							WriteBufferSize:   oracleconfig.DefaultWriteBufferSize,
							HandshakeTimeout:  oracleconfig.DefaultHandshakeTimeout,
							EnableCompression: oracleconfig.DefaultEnableCompression,
							ReadTimeout:       oracleconfig.DefaultReadTimeout,
							WriteTimeout:      oracleconfig.DefaultWriteTimeout,
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
			if legacyErr, shouldBeNil := tc.config.ValidateBasic(), err == nil; (legacyErr == nil) != shouldBeNil {
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

	expectedConfig := filterMarketMapProvidersFromOracleConfig(cmdconfig.DefaultOracleConfig(), marketmap.Name)
	require.NoError(t, expectedConfig.ValidateBasic())
	expectedConfig.UpdateInterval = updateIntervalOverride
	provider := expectedConfig.Providers[raydium.Name]
	provider.API.Endpoints = append(provider.API.Endpoints, endpointOverride)
	expectedConfig.Providers[raydium.Name] = provider
	expectedConfig.Metrics.PrometheusServerAddress = prometheusServerOverride

	coinbase := expectedConfig.Providers[coinbase.Name]
	coinbase.WebSocket.Endpoints = []oracleconfig.Endpoint{{URL: endpointOverride.URL}}
	expectedConfig.Providers[coinbase.Name] = coinbase

	t.Run("overriding variables from environment", func(t *testing.T) {
		// set the environment variables
		t.Setenv(cmdconfig.ConnectConfigEnvironmentPrefix+"_UPDATEINTERVAL", updateIntervalOverride.String())
		t.Setenv(cmdconfig.ConnectConfigEnvironmentPrefix+"_METRICS_PROMETHEUSSERVERADDRESS", prometheusServerOverride)
		t.Setenv(cmdconfig.ConnectConfigEnvironmentPrefix+"_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_URL", endpointOverride.URL)
		t.Setenv(cmdconfig.ConnectConfigEnvironmentPrefix+"_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_AUTHENTICATION_APIKEY", endpointOverride.Authentication.APIKey)
		t.Setenv(cmdconfig.ConnectConfigEnvironmentPrefix+"_PROVIDERS_RAYDIUM_API_API_ENDPOINTS_1_AUTHENTICATION_APIKEYHEADER", endpointOverride.Authentication.APIKeyHeader)
		t.Setenv(cmdconfig.ConnectConfigEnvironmentPrefix+"_PROVIDERS_COINBASE_WS_WEBSOCKET_ENDPOINTS_0_URL", endpointOverride.URL)

		cfg, err := cmdconfig.ReadOracleConfigWithOverrides("", marketmap.Name)
		require.NoError(t, err)

		require.Equal(t, expectedConfig.Providers, cfg.Providers)
		require.Equal(t, expectedConfig.UpdateInterval, cfg.UpdateInterval)
		require.Equal(t, expectedConfig.Metrics.PrometheusServerAddress, cfg.Metrics.PrometheusServerAddress)
	})

	t.Run("overriding variables via config", func(t *testing.T) {
		// create a temp file in the current directory
		tmpfile, err := os.CreateTemp("", "connect-config-*.json")
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
				},
				"%s": {
					"webSocket": {
						"endpoints": [
							{
								"url": "%s"
							}
						]
					}
				}
			}
		}
		`,
			updateIntervalOverride,
			prometheusServerOverride,
			raydium.Name,
			raydium.DefaultAPIConfig.Endpoints[0].URL,
			endpointOverride.URL,
			endpointOverride.Authentication.APIKey,
			endpointOverride.Authentication.APIKeyHeader,
			coinbase.Name,
			endpointOverride.URL,
		)
		tmpfile.Write([]byte(overrides))

		cfg, err := cmdconfig.ReadOracleConfigWithOverrides(tmpfile.Name(), marketmap.Name)
		require.NoError(t, err)

		require.Equal(t, expectedConfig.Providers, cfg.Providers)
		require.Equal(t, expectedConfig.UpdateInterval, cfg.UpdateInterval)
		require.Equal(t, expectedConfig.Metrics.PrometheusServerAddress, cfg.Metrics.PrometheusServerAddress)
	})

	t.Run("overriding a nonexistent provider via config fails", func(t *testing.T) {
		// create a temp file in the current directory
		tmpfile, err := os.CreateTemp("", "connect-config-*.json")
		require.NoError(t, err)

		defer os.Remove(tmpfile.Name())

		overrides := fmt.Sprintf(`
		{
			"updateInterval": "%s",
			"metrics": {
				"prometheusServerAddress": "%s"
			},
			"providers": {
				"doesNotExist": {
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
		`,
			updateIntervalOverride,
			prometheusServerOverride,
			raydium.DefaultAPIConfig.Endpoints[0].URL,
			endpointOverride.URL,
			endpointOverride.Authentication.APIKey,
			endpointOverride.Authentication.APIKeyHeader,
		)
		tmpfile.Write([]byte(overrides))

		_, err = cmdconfig.ReadOracleConfigWithOverrides(tmpfile.Name(), marketmap.Name)
		require.ErrorContains(t, err, "overridden key")
	})
}

func TestOracleConfigWithExtraKeys(t *testing.T) {
	t.Run("an oracle config with extraneous keys", func(t *testing.T) {
		// create a temp file in the current directory
		tmpfile, err := os.CreateTemp("", "connect-config-*.json")
		require.NoError(t, err)

		defer os.Remove(tmpfile.Name())

		overrides := `
		{
			"providers": {
				"raydium_api": {
					"api": {
						"endpoints": [
							{
								"url": "http://somewhere",
								"some_field_that_is_not_relevant": ""
							}
						]
					}
				}
			}
		}
		`
		tmpfile.Write([]byte(overrides))

		_, err = cmdconfig.ReadOracleConfigWithOverrides(tmpfile.Name(), marketmap.Name)
		require.Error(t, err)
	})
}

func filterMarketMapProvidersFromOracleConfig(cfg oracleconfig.OracleConfig, mmProvider string) oracleconfig.OracleConfig {
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
