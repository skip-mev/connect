package metrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// OracleSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	OracleSubsystem = "oracle"
)

// Metrics is an interface that defines the API for oracle metrics.
//
//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// AddTick increments the number of ticks, this can represent a liveness counter. This
	// is incremented once every interval (which is defined by the oracle config).
	AddTick()
	//
	// TODO: Add more metrics here in later PRs.
}

// OracleMetricsImpl is a Metrics implementation that does nothing.
type OracleMetricsImpl struct {
	ticks metrics.Counter
}

// NewMetricsFromConfig returns a oracle Metrics implementation based on the provided
// config.
func NewMetricsFromConfig(config config.OracleMetricsConfig) Metrics {
	if config.Enabled {
		return NewMetrics()
	}
	return NewNopMetrics()
}

// NewMetrics returns a Metrics implementation that exposes metrics to Prometheus.
func NewMetrics() Metrics {
	m := &OracleMetricsImpl{
		ticks: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "ticks",
			Help:      "Number of ticks with a successful oracle update.",
		}, []string{}),
	}

	return m
}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &OracleMetricsImpl{
		ticks: discard.NewCounter(),
	}
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *OracleMetricsImpl) AddTick() {
	m.ticks.Add(1)
}
