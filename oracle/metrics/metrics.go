package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	OracleSubsystem = "oracle"
)

type Config struct {
	// Enabled indicates whether metrics should be enabled
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}

//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// AddTick increments the number of ticks, this can represent a liveness counter. This metric is paginated by status.
	AddTick()
}

type nopMetricsImpl struct{}

func NewNopMetrics() Metrics {
	return &nopMetricsImpl{}
}

func (m *nopMetricsImpl) AddTick() {}

func NewMetrics() Metrics {
	m := &metricsImpl{
		ticks: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "ticks",
			Help:      "Number of ticks with a fully successful Oracle update (all providers returned).",
		}),
	}

	// register the metrics
	prometheus.MustRegister(m.ticks)

	return m
}

// Metrics contains metrics exposed by this package.
type metricsImpl struct {
	// Number of ticks with a fully successful Oracle update (all providers returned).
	ticks prometheus.Counter
}

func (m *metricsImpl) AddTick() {
	m.ticks.Add(1)
}
