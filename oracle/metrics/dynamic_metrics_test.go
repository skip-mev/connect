package metrics

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/metrics/mocks"
)

func TestDetermineMetricsType(t *testing.T) {
	tcs := []struct {
		name    string
		metrics Metrics
		mType   ImplType
	}{
		{
			name:    "oracle metrics type",
			metrics: &OracleMetricsImpl{},
			mType:   OracleMetricsType,
		},
		{
			name:    "noop metrics type",
			metrics: &noOpOracleMetrics{},
			mType:   NoOpMetricsType,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			typ := determineMetricsType(tc.metrics)
			require.Equal(t, tc.mType, typ)
		})
	}
}

// TestDynamicMetrics_Switches tests that the metrics impl will switch if it can communicate with the node.
func TestDynamicMetrics_Switches(t *testing.T) {
	ctx := context.Background()
	cfg := config.MetricsConfig{
		PrometheusServerAddress: "",
		Telemetry: config.TelemetryConfig{
			Disabled:    false,
			PushAddress: "",
		},
		Enabled: false,
	}
	node := mocks.NewNodeClient(t)

	// i'm using a mutex to block the DeriveNodeIdentifier call so i can prevent the for/select from looping.
	mtx := sync.Mutex{}
	mtx.Lock()

	// it gets called once in the loop where it checks the node,
	// and again in NewMetricsFromConfig.
	node.EXPECT().DeriveNodeIdentifier().Run(func() {
		mtx.Lock()
		mtx.Unlock() //nolint:staticcheck
	}).Return("foobar", nil).Twice()

	dyn := NewDynamicMetrics(ctx, cfg, node)
	dynImpl, ok := dyn.(*dynamicMetrics)
	require.True(t, ok)
	require.Equal(t, dynImpl.metricsType, NoOpMetricsType)

	dynImpl.cfg.Enabled = true

	mtx.Unlock()

	valid := false
	for range 10 {
		if dynImpl.metricsType == OracleMetricsType {
			valid = true
			break
		}
		time.Sleep(time.Millisecond * 500)
	}

	require.True(t, valid, "the metrics type did not change after 500ms", "metricsType", dynImpl.metricsType)
}
