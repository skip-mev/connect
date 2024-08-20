package config

import (
	"fmt"
)

// MetricsConfig is the metrics configurations for the oracle. This configuration object specifically
// exposes metrics pertaining to the oracle sidecar. To enable app side metrics, please see the app
// configuration.
type MetricsConfig struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose
	// metrics to.
	PrometheusServerAddress string `json:"prometheusServerAddress"`

	Telemetry TelemetryConfig `json:"telemetry"`

	// Enabled indicates whether metrics should be enabled.
	Enabled bool `json:"enabled"`
}

type TelemetryConfig struct {
	// Toggle to disable opt-out telemetry
	Disabled bool `json:"disabled"`

	// Address of the remote server to push telemetry to
	PushAddress string `json:"pushAddress"`
}

// ValidateBasic performs basic validation of the config.
func (c *MetricsConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if len(c.PrometheusServerAddress) == 0 {
		return fmt.Errorf("must supply a non-empty prometheus server address if metrics are enabled")
	}

	return c.Telemetry.ValidateBasic()
}

// ValidateBasic performs basic validation of the config.
func (c *TelemetryConfig) ValidateBasic() error {
	if c.Disabled {
		return nil
	}

	if len(c.PushAddress) == 0 {
		return fmt.Errorf("must supply a non-empty push address when telemetry is not disabled")
	}

	return nil
}
