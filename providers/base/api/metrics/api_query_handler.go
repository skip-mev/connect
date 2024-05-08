package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

const (
	// StatusLabel is a label for the status of a provider API response.
	StatusLabel = "internal_status"

	// StatusCodeLabel is a label for the status code of a provider API response.
	StatusCodeLabel = "status_code"
)

// APIMetrics is an interface that defines the API for metrics collection for providers
// that implement the APIQueryHandler.
//
//go:generate mockery --name APIMetrics --filename mock_metrics.go
type APIMetrics interface {
	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	// This increments the number of responses by provider, id (i.e. currency pair), and status.
	AddProviderResponse(providerName, id string, errorCode providertypes.ErrorCode)

	// AddHTTPStatusCode increments the number of responses by provider and status.
	// This is used to track the number of responses by provider and status.
	AddHTTPStatusCode(providerName string, resp *http.Response)

	// ObserveProviderResponseLatency records the time it took for a provider to respond for
	// within a single interval. Note that if the provider is not atomic, this will be the
	// time it took for all the requests to complete.
	ObserveProviderResponseLatency(providerName string, duration time.Duration)
}

// APIMetricsImpl contains metrics exposed by this package.
type APIMetricsImpl struct {
	// Number of provider successes.
	apiResponseStatusPerProvider *prometheus.CounterVec

	// Number of provider responses by grouped status code.
	apiHTTPStatusCodePerProvider *prometheus.CounterVec

	// Number of provider responses by status code with the exact status code.
	apiHTTPStatusCodeExactPerProvider *prometheus.CounterVec

	// Histogram paginated by provider, measuring the latency between invocation and collection.
	apiResponseTimePerProvider *prometheus.HistogramVec
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
		apiResponseStatusPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "api_response_internal_status",
			Help:      "Number of API provider successes.",
		}, []string{providermetrics.ProviderLabel, providermetrics.IDLabel, StatusLabel}),
		apiHTTPStatusCodePerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "api_http_status_code",
			Help:      "Number of API provider responses by status code grouped by category (2XX, 3XX, etc.).",
		}, []string{providermetrics.ProviderLabel, StatusCodeLabel}),
		apiHTTPStatusCodeExactPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "api_http_status_code_granular",
			Help:      "Number of API provider responses by status code with granularity on the exact code.",
		}, []string{providermetrics.ProviderLabel, StatusCodeLabel}),
		apiResponseTimePerProvider: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "api_response_latency",
			Help:      "Response time per API provider.",
			Buckets:   []float64{50, 100, 250, 500, 1000, 2000},
		}, []string{providermetrics.ProviderLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.apiResponseStatusPerProvider)
	prometheus.MustRegister(m.apiHTTPStatusCodePerProvider)
	prometheus.MustRegister(m.apiResponseTimePerProvider)
	prometheus.MustRegister(m.apiHTTPStatusCodeExactPerProvider)

	return m
}

type noOpAPIMetricsImpl struct{}

// NewNopAPIMetrics returns a Provider Metrics implementation that does nothing.
func NewNopAPIMetrics() APIMetrics {
	return &noOpAPIMetricsImpl{}
}

func (m *noOpAPIMetricsImpl) AddProviderResponse(_ string, _ string, _ providertypes.ErrorCode) {}
func (m *noOpAPIMetricsImpl) AddHTTPStatusCode(_ string, _ *http.Response)                      {}
func (m *noOpAPIMetricsImpl) ObserveProviderResponseLatency(_ string, _ time.Duration)          {}

// AddProviderResponse increments the number of requests by provider and status.
func (m *APIMetricsImpl) AddProviderResponse(providerName string, id string, err providertypes.ErrorCode) {
	var status string
	if err.Error() == nil {
		status = "success"
	} else {
		status = err.Error().Error()
	}

	m.apiResponseStatusPerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: providerName,
		providermetrics.IDLabel:       id,
		StatusLabel:                   status,
	},
	).Add(1)
}

// AddHTTPStatusCode increments the http status code by provider and response.
func (m *APIMetricsImpl) AddHTTPStatusCode(providerName string, resp *http.Response) {
	var (
		status      string
		statusExact string
	)
	switch {
	case resp == nil || resp.StatusCode >= 500:
		status = "5XX"
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		status = "2XX"
	case resp.StatusCode >= 300 && resp.StatusCode < 400:
		status = "3XX"
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		status = "4XX"
	}

	if resp != nil {
		statusExact = fmt.Sprintf("%d", resp.StatusCode)
	} else {
		statusExact = "500"
	}

	m.apiHTTPStatusCodePerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: providerName,
		StatusCodeLabel:               status,
	}).Add(1)

	m.apiHTTPStatusCodeExactPerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: providerName,
		StatusCodeLabel:               statusExact,
	}).Add(1)
}

// ObserveProviderResponseLatency records the time it took for a provider to respond.
func (m *APIMetricsImpl) ObserveProviderResponseLatency(providerName string, duration time.Duration) {
	m.apiResponseTimePerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: providerName,
	},
	).Observe(float64(duration.Milliseconds()))
}
