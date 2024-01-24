package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
)

const (
	// StatusLabel is a label for the status of a provider API response.
	StatusLabel = "status"
)

// APIMetrics is an interface that defines the API for metrics collection for providers
// that implement the APIQueryHandler.
//
//go:generate mockery --name APIMetrics --filename mock_metrics.go
type APIMetrics interface {
	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	// This increments the number of responses by provider, id (i.e. currency pair), and status.
	AddProviderResponse(providerName, id string, status Status)

	// ObserveProviderResponseLatency records the time it took for a provider to respond for
	// within a single interval. Note that if the provider is not atomic, this will be the
	// time it took for all the requests to complete.
	ObserveProviderResponseLatency(providerName string, duration time.Duration)
}

// APIMetricsImpl contains metrics exposed by this package.
type APIMetricsImpl struct {
	// Number of provider successes.
	responseStatusPerProvider *prometheus.CounterVec

	// Histogram paginated by provider, measuring the latency between invocation and collection.
	responseTimePerProvider *prometheus.HistogramVec
}

// NewAPIMetricsFromConfig returns a new Metrics struct given the main oracle metrics config.
func NewAPIMetricsFromConfig(config config.MetricsConfig) APIMetrics {
	if config.Enabled {
		return NewAPIMetrics()
	}
	return NewNopAPIMetrics()
}

// NewAPIMetrics returns a Provider Metrics implementation that uses Prometheus.
func NewAPIMetrics() APIMetrics {
	m := &APIMetricsImpl{
		responseStatusPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "api_response_status_per_provider",
			Help:      "Number of API provider successes.",
		}, []string{providermetrics.ProviderLabel, providermetrics.IDLabel, StatusLabel}),
		responseTimePerProvider: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "api_response_time_per_provider",
			Help:      "Response time per API provider.",
			Buckets:   []float64{50, 100, 250, 500, 1000},
		}, []string{providermetrics.ProviderLabel}),
	}

	return m
}

type noOpAPIMetricsImpl struct{}

// NewNopAPIMetrics returns a Provider Metrics implementation that does nothing.
func NewNopAPIMetrics() APIMetrics {
	return &noOpAPIMetricsImpl{}
}

func (m *noOpAPIMetricsImpl) AddProviderResponse(_ string, _ string, _ Status)         {}
func (m *noOpAPIMetricsImpl) ObserveProviderResponseLatency(_ string, _ time.Duration) {}

// AddProviderResponse increments the number of requests by provider and status.
func (m *APIMetricsImpl) AddProviderResponse(providerName string, id string, status Status) {
	m.responseStatusPerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: providerName,
		providermetrics.IDLabel:       id,
		StatusLabel:                   status.String(),
	},
	).Add(1)
}

// ObserveProviderResponseLatency records the time it took for a provider to respond.
func (m *APIMetricsImpl) ObserveProviderResponseLatency(providerName string, duration time.Duration) {
	m.responseTimePerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: providerName,
	},
	).Observe(float64(duration.Milliseconds()))
}
