package config_test

import (
	"testing"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
)

func TestValidateBasic(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.AppConfig
		expectedErr bool
	}{
		{
			name:        "good config with a disabled oracle",
			config:      config.AppConfig{},
			expectedErr: false,
		},
		{
			name: "good config with no metrics",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				PriceTTL:      time.Second * 2,
			},
			expectedErr: false,
		},
		{
			name: "good config with metrics",
			config: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8080",
				ClientTimeout:  time.Second,
				MetricsEnabled: true,
				Interval:       time.Second,
				PriceTTL:       time.Second * 2,
			},
			expectedErr: false,
		},
		{
			name: "bad config with no oracle address",
			config: config.AppConfig{
				Enabled:       true,
				ClientTimeout: time.Second,
				Interval:      time.Second,
				PriceTTL:      time.Second * 2,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no client timeout",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				Interval:      time.Second,
				PriceTTL:      time.Second * 2,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no max age",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
			},
			expectedErr: true,
		},
		{
			name: "bad config with no interval",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				PriceTTL:      time.Second * 2,
			},
			expectedErr: true,
		},
		{
			name: "bad config with max age being less than interval",
			config: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				PriceTTL:      time.Millisecond,
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfigFromAppOptions(t *testing.T) {
	testCases := []struct {
		name        string
		config      servertypes.AppOptions
		res         config.AppConfig
		expectedErr bool
	}{
		{
			name: "good config",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "5s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "20s",
				"oracle.interval":        "10s",
			},
			res: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8081",
				ClientTimeout:  5 * time.Second,
				MetricsEnabled: true,
				PriceTTL:       20 * time.Second,
				Interval:       10 * time.Second,
			},
			expectedErr: false,
		},
		{
			name:        "good config with no fields configured",
			config:      sims.AppOptionsMap{},
			res:         config.NewDefaultAppConfig(),
			expectedErr: false,
		},
		{
			name: "good config with no oracle address specified",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.client_timeout":  "5s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "20s",
				"oracle.interval":        "10s",
			},
			res: config.AppConfig{
				Enabled:        true,
				OracleAddress:  config.DefaultOracleAddress,
				ClientTimeout:  5 * time.Second,
				MetricsEnabled: true,
				PriceTTL:       20 * time.Second,
				Interval:       10 * time.Second,
			},
			expectedErr: false,
		},
		{
			name: "good config with no client timeout specified",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "20s",
				"oracle.interval":        "10s",
			},
			res: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8081",
				ClientTimeout:  config.DefaultClientTimeout,
				MetricsEnabled: true,
				PriceTTL:       20 * time.Second,
				Interval:       10 * time.Second,
			},
			expectedErr: false,
		},
		{
			name: "good config with no metrics enabled specified",
			config: sims.AppOptionsMap{
				"oracle.enabled":        true,
				"oracle.oracle_address": "localhost:8081",
				"oracle.client_timeout": "5s",
				"oracle.price_ttl":      "20s",
				"oracle.interval":       "10s",
			},
			res: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8081",
				ClientTimeout:  5 * time.Second,
				MetricsEnabled: config.DefaultMetricsEnabled,
				PriceTTL:       20 * time.Second,
				Interval:       10 * time.Second,
			},
			expectedErr: false,
		},
		{
			name: "good config with no price ttl specified",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "5s",
				"oracle.metrics_enabled": true,
				"oracle.interval":        "2s",
			},
			res: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8081",
				ClientTimeout:  5 * time.Second,
				MetricsEnabled: true,
				PriceTTL:       config.DefaultPriceTTL,
				Interval:       2 * time.Second,
			},
			expectedErr: false,
		},
		{
			name: "good config with no interval specified",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "5s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "20s",
			},
			res: config.AppConfig{
				Enabled:        true,
				OracleAddress:  "localhost:8081",
				ClientTimeout:  5 * time.Second,
				MetricsEnabled: true,
				PriceTTL:       20 * time.Second,
				Interval:       config.DefaultInterval,
			},
			expectedErr: false,
		},
		{
			name: "bad config with bad client timeout",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "mogged",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "20s",
				"oracle.interval":        "10s",
			},
			res:         config.AppConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with bad price ttl",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "1s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "mogged",
				"oracle.interval":        "10s",
			},
			res:         config.AppConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with bad interval",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "1s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "10s",
				"oracle.interval":        "mogged",
			},
			res:         config.AppConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with price ttl < interval",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "1s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "1s",
				"oracle.interval":        "10s",
			},
			res:         config.AppConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with price ttl exceeding max",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "1s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "100m",
				"oracle.interval":        "10s",
			},
			res:         config.AppConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with interval exceeding max",
			config: sims.AppOptionsMap{
				"oracle.enabled":         true,
				"oracle.oracle_address":  "localhost:8081",
				"oracle.client_timeout":  "1s",
				"oracle.metrics_enabled": true,
				"oracle.price_ttl":       "10s",
				"oracle.interval":        "100m",
			},
			res:         config.AppConfig{},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := config.ReadConfigFromAppOpts(tc.config)
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
