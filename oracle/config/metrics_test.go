package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
)

func TestMetricsConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.MetricsConfig
		expectedErr bool
	}{
		{
			name: "good config with metrics",
			config: config.MetricsConfig{
				Enabled:                 true,
				PrometheusServerAddress: "localhost:9090",
				Telemetry: config.TelemetryConfig{
					Disabled:    false,
					PushAddress: "localhost:9125",
				},
			},
			expectedErr: false,
		},
		{
			name: "bad config with no prometheus server address",
			config: config.MetricsConfig{
				Enabled:                 true,
				PrometheusServerAddress: "",
			},
			expectedErr: true,
		},
		{
			name: "bad config with no telemetry push address",
			config: config.MetricsConfig{
				Enabled: true,
				Telemetry: config.TelemetryConfig{
					Disabled:    false,
					PushAddress: "",
				},
			},
			expectedErr: true,
		},
		{
			name: "no metrics enabled",
			config: config.MetricsConfig{
				Enabled:                 false,
				PrometheusServerAddress: "",
			},
			expectedErr: false,
		},
		{
			name: "telemetry disabled",
			config: config.MetricsConfig{
				Enabled:                 false,
				PrometheusServerAddress: "",
				Telemetry:               config.TelemetryConfig{},
			},
			expectedErr: false,
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
