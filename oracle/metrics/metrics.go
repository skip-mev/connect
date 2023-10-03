package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	OracleSubsystem = "oracle"
	ProviderLabel   = "provider"
	StatusLabel     = "status"
)

type Config struct {
	// Enabled indicates whether metrics should be enabled
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}

type Status int

const (
	StatusFailure Status = iota
	StatusSuccess
)

func (s Status) String() string {
	switch s {
	case StatusFailure:
		return "failure"
	case StatusSuccess:
		return "success"
	default:
		return "unknown"
	}
}

func StatusFromError(err error) Status {
	if err == nil {
		return StatusSuccess
	}
	return StatusFailure
}

//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// AddProviderResponse increments the number of ticks with a fully successful Oracle update (all providers returned).
	AddProviderResponse(providerName string, status Status)

	// AddTick increments the number of ticks, this can represent a liveness counter. This metric is paginated by status.
	AddTick()

	// ObserveProviderResponseTime records the time it took for a provider to respond
	ObserveProviderResponseLatency(providerName string, duration time.Duration)
}

type nopMetricsImpl struct{}

func NewNopMetrics() Metrics {
	return &nopMetricsImpl{}
}

func (m *nopMetricsImpl) AddProviderResponse(_ string, _ Status)                   {}
func (m *nopMetricsImpl) AddTick()                                                 {}
func (m *nopMetricsImpl) ObserveProviderResponseLatency(_ string, _ time.Duration) {}

func NewMetrics() Metrics {
	m := &metricsImpl{
		ticks: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "ticks",
			Help:      "Number of ticks with a fully successful Oracle update (all providers returned).",
		}),
		responseStatusPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "response_status_per_provider",
			Help:      "Number of provider successes.",
		}, []string{ProviderLabel, StatusLabel}),
		responseTimePerProvider: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: OracleSubsystem,
			Name:      "response_time_per_provider",
			Help:      "ResponseTimePerProvider",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}, []string{ProviderLabel}),
	}

	// register the metrics
	prometheus.MustRegister(m.ticks)
	prometheus.MustRegister(m.responseStatusPerProvider)
	prometheus.MustRegister(m.responseTimePerProvider)

	return m
}

// Metrics contains metrics exposed by this package.
type metricsImpl struct {
	// Number of ticks with a fully successful Oracle update (all providers returned).
	ticks prometheus.Counter

	// Number of provider successes.
	responseStatusPerProvider *prometheus.CounterVec

	// histogram paginated by provider, measuring the latency between invocation and collection (of all responses)
	responseTimePerProvider *prometheus.HistogramVec
}

func (m *metricsImpl) AddProviderResponse(providerName string, status Status) {
	m.responseStatusPerProvider.WithLabelValues(providerName, status.String()).Add(1)
}

func (m *metricsImpl) AddTick() {
	m.ticks.Add(1)
}

func (m *metricsImpl) ObserveProviderResponseLatency(providerName string, duration time.Duration) {
	m.responseTimePerProvider.WithLabelValues(providerName).Observe(float64(duration.Milliseconds()))
}
