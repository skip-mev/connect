package metrics

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
)

const (
	// StatusLabel is the label used for the status of a provider response.
	StatusLabel = "status"
)

// WebSocketMetrics is an interface that defines the API for metrics collection for providers
// that implement the WebSocketQueryHandler.
//
//go:generate mockery --name WebSocketMetrics --filename mock_metrics.go
type WebSocketMetrics interface {
	// AddWebSocketConnectionStatus adds a method / status response to the metrics collector for the
	// given provider. Specifically, this tracks various connection related errors.
	AddWebSocketConnectionStatus(provider string, status ConnectionStatus)

	// AddWebSocketDataHandlerStatus adds a method / status response to the metrics collector for the
	// given provider. Specifically, this tracks various data handler related errors.
	AddWebSocketDataHandlerStatus(provider string, status HandlerStatus)

	// ObserveWebSocketLatency adds a latency observation to the metrics collector for the
	// given provider.
	ObserveWebSocketLatency(provider string, duration time.Duration)
}

// WebSocketMetricsImpl contains metrics exposed by this package.
type WebSocketMetricsImpl struct {
	// Number of connection successes.
	connectionStatusPerProvider metrics.Counter

	// Number of data handler successes.
	dataHandlerStatusPerProvider metrics.Counter

	// Histogram paginated by provider, measuring the latency between invocation and collection.
	responseTimePerProvider metrics.Histogram
}

// NewWebSocketMetricsFromConfig returns a new Metrics struct given the main oracle metrics config.
func NewWebSocketMetricsFromConfig(config config.OracleMetricsConfig) WebSocketMetrics {
	if config.Enabled {
		return NewWebSocketMetrics()
	}
	return NewNopWebSocketMetrics()
}

// NewWebSocketMetrics returns a Provider Metrics implementation that uses Prometheus.
func NewWebSocketMetrics() WebSocketMetrics {
	m := &WebSocketMetricsImpl{
		connectionStatusPerProvider: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "web_socket_connection_status_per_provider",
			Help:      "Number of web socket connection successes.",
		}, []string{providermetrics.ProviderLabel, StatusLabel}),
		dataHandlerStatusPerProvider: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "web_socket_data_handler_status_per_provider",
			Help:      "Number of web socket data handler successes.",
		}, []string{providermetrics.ProviderLabel, StatusLabel}),
		responseTimePerProvider: prometheus.NewHistogramFrom(stdprom.HistogramOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "web_socket_response_time_per_provider",
			Help:      "Response time per API provider.",
			Buckets:   []float64{50, 100, 250, 500, 1000},
		}, []string{providermetrics.ProviderLabel}),
	}

	return m
}

// NewNopWebSocketMetrics returns a Provider Metrics implementation that does not collect metrics.
func NewNopWebSocketMetrics() WebSocketMetrics {
	return &WebSocketMetricsImpl{
		connectionStatusPerProvider:  discard.NewCounter(),
		dataHandlerStatusPerProvider: discard.NewCounter(),
		responseTimePerProvider:      discard.NewHistogram(),
	}
}

// AddWebSocketConnectionStatus adds a method / status response to the metrics collector for the
// given provider. Specifically, this tracks various connection related errors.
func (m *WebSocketMetricsImpl) AddWebSocketConnectionStatus(provider string, status ConnectionStatus) {
	m.connectionStatusPerProvider.With(providermetrics.ProviderLabel, provider, StatusLabel, status.String()).Add(1)
}

// AddWebSocketDataHandlerStatus adds a method / status response to the metrics collector for the
// given provider. Specifically, this tracks various data handler related errors.
func (m *WebSocketMetricsImpl) AddWebSocketDataHandlerStatus(provider string, status HandlerStatus) {
	m.dataHandlerStatusPerProvider.With(providermetrics.ProviderLabel, provider, StatusLabel, status.String()).Add(1)
}

// ObserveWebSocketLatency adds a latency observation to the metrics collector for the given provider.
func (m *WebSocketMetricsImpl) ObserveWebSocketLatency(provider string, duration time.Duration) {
	m.responseTimePerProvider.With(providermetrics.ProviderLabel, provider).Observe(float64(duration.Milliseconds()))
}
