package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/connect/v2/oracle/config"
	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
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
	connectionStatusPerProvider *prometheus.CounterVec

	// Number of data handler successes.
	dataHandlerStatusPerProvider *prometheus.CounterVec

	// Histogram paginated by provider, measuring the latency between invocation and collection.
	responseTimePerProvider *prometheus.HistogramVec
}

// NewWebSocketMetricsFromConfig returns a new Metrics struct given the main oracle metrics config.
func NewWebSocketMetricsFromConfig(config config.MetricsConfig) WebSocketMetrics {
	if config.Enabled {
		return NewWebSocketMetrics()
	}
	return NewNopWebSocketMetrics()
}

// NewWebSocketMetrics returns a Provider Metrics implementation that uses Prometheus.
func NewWebSocketMetrics() WebSocketMetrics {
	m := &WebSocketMetricsImpl{
		connectionStatusPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "web_socket_connection_status",
			Help:      "Statuses associated with the underlying web socket connection.",
		}, []string{providermetrics.ProviderLabel, StatusLabel}),
		dataHandlerStatusPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "web_socket_data_handler_status",
			Help:      "Statuses associated with parsing/sending web socket messages from/to a web socket connection.",
		}, []string{providermetrics.ProviderLabel, StatusLabel}),
		responseTimePerProvider: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "web_socket_response_time",
			Help:      "Response time per web socket provider.",
			Buckets:   []float64{50, 100, 250, 500, 1000, 2000},
		}, []string{providermetrics.ProviderLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.connectionStatusPerProvider)
	prometheus.MustRegister(m.dataHandlerStatusPerProvider)
	prometheus.MustRegister(m.responseTimePerProvider)

	return m
}

type noOpWebSocketMetricsImpl struct{}

// NewNopWebSocketMetrics returns a Provider Metrics implementation that does not collect metrics.
func NewNopWebSocketMetrics() WebSocketMetrics {
	return &noOpWebSocketMetricsImpl{}
}

func (m *noOpWebSocketMetricsImpl) AddWebSocketConnectionStatus(_ string, _ ConnectionStatus) {
}

func (m *noOpWebSocketMetricsImpl) AddWebSocketDataHandlerStatus(_ string, _ HandlerStatus) {
}

func (m *noOpWebSocketMetricsImpl) ObserveWebSocketLatency(_ string, _ time.Duration) {
}

// AddWebSocketConnectionStatus adds a method / status response to the metrics collector for the
// given provider. Specifically, this tracks various connection related errors.
func (m *WebSocketMetricsImpl) AddWebSocketConnectionStatus(provider string, status ConnectionStatus) {
	m.connectionStatusPerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: provider,
		StatusLabel:                   status.String(),
	},
	).Add(1)
}

// AddWebSocketDataHandlerStatus adds a method / status response to the metrics collector for the
// given provider. Specifically, this tracks various data handler related errors.
func (m *WebSocketMetricsImpl) AddWebSocketDataHandlerStatus(provider string, status HandlerStatus) {
	m.dataHandlerStatusPerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: provider,
		StatusLabel:                   status.String(),
	},
	).Add(1)
}

// ObserveWebSocketLatency adds a latency observation to the metrics collector for the given provider.
func (m *WebSocketMetricsImpl) ObserveWebSocketLatency(provider string, duration time.Duration) {
	m.responseTimePerProvider.With(prometheus.Labels{
		providermetrics.ProviderLabel: provider,
	},
	).Observe(float64(duration.Milliseconds()))
}
