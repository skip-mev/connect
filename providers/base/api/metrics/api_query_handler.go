package metrics

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"
	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
)

const (
	// ProviderLabel is a label for the provider name.
	ProviderLabel = "provider"
	// StatusLabel is a label for the status of a provider response.
	StatusLabel = "status"
	// IDLabel is a label for the ID of a provider response.
	IDLabel = "id"
)

// APIMetrics is an interface that defines the API for metrics collection for providers
// that implement the APIQueryHandler.
//
//go:generate mockery --name APIMetrics --filename mock_metrics.go
type APIMetrics interface {
	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	// This increments the number of responses by provider, id (i.e. currency pair), and status.
	AddProviderResponse(providerName, id string, status Status)

	// ObserveProviderResponseTime records the time it took for a provider to respond for
	// within a single interval. Note that if the provider is not atomic, this will be the
	// time it took for all of the requests to complete.
	ObserveProviderResponseLatency(providerName string, duration time.Duration)
}

// APIMetricsImpl contains metrics exposed by this package.
type APIMetricsImpl struct {
	// Number of provider successes.
	responseStatusPerProvider metrics.Counter

	// Histogram paginated by provider, measuring the latency between invocation and collection.
	responseTimePerProvider metrics.Histogram
}

// NewAPIMetricsFromConfig returns a new Metrics struct given the main oracle metrics config.
func NewAPIMetricsFromConfig(config config.OracleMetricsConfig) APIMetrics {
	if config.Enabled {
		return NewAPIMetrics()
	}
	return NewNopAPIMetrics()
}

// NewAPIMetrics returns a Provider Metrics implementation that uses Prometheus.
func NewAPIMetrics() APIMetrics {
	m := &APIMetricsImpl{
		responseStatusPerProvider: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "response_status_per_provider",
			Help:      "Number of provider successes.",
		}, []string{ProviderLabel, IDLabel, StatusLabel}),
		responseTimePerProvider: prometheus.NewHistogramFrom(stdprom.HistogramOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "response_time_per_provider",
			Help:      "Response time per provider.",
			Buckets:   []float64{50, 100, 250, 500, 1000},
		}, []string{ProviderLabel}),
	}

	return m
}

// NewNopAPIMetrics returns a Provider Metrics implementation that does nothing.
func NewNopAPIMetrics() APIMetrics {
	return &APIMetricsImpl{
		responseStatusPerProvider: discard.NewCounter(),
		responseTimePerProvider:   discard.NewHistogram(),
	}
}

// AddProviderResponse increments the number of requests by provider and status.
func (m *APIMetricsImpl) AddProviderResponse(providerName string, id string, status Status) {
	m.responseStatusPerProvider.With(ProviderLabel, providerName, IDLabel, id, StatusLabel, status.String()).Add(1)
}

// ObserveProviderResponseTime records the time it took for a provider to respond.
func (m *APIMetricsImpl) ObserveProviderResponseLatency(providerName string, duration time.Duration) {
	m.responseTimePerProvider.With(ProviderLabel, providerName).Observe(float64(duration.Milliseconds()))
}
