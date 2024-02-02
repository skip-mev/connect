package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

const (
	// ProviderLabel is a label for the provider name.
	ProviderLabel = "provider"
	// IDLabel is a label for the ID of a provider response.
	IDLabel = "id"
	// ProviderTypeLabel is a label for the type of provider (WS, API, etc.)
	ProviderTypeLabel = "type"
	// StatusLabel is a label for the status of a provider response.
	StatusLabel = "status"
)

type (
	Status string
)

const (
	Success Status = "success"
	Failure Status = "failure"
)

// ProviderMetrics is an interface that defines the API for metrics collection for providers. The
// base provider utilizes this interface to collect metrics, whether the underlying implementation
// is API or websocket based.
//
//go:generate mockery --name ProviderMetrics --filename mock_metrics.go
type ProviderMetrics interface {
	// AddProviderResponseByID increments the number of ticks with a fully successful provider update
	// for a given provider and ID (i.e. currency pair).
	AddProviderResponseByID(providerName, id string, status Status, providerType providertypes.ProviderType)

	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	AddProviderResponse(providerName string, status Status, providerType providertypes.ProviderType)

	// LastUpdated updates the last time a given ID (i.e. currency pair) was updated.
	LastUpdated(providerName, id string, providerType providertypes.ProviderType)
}

// ProviderMetricsImpl contains metrics exposed by this package.
type ProviderMetricsImpl struct {
	// Number of provider successes by ID.
	responseStatusPerProviderByID *prometheus.CounterVec

	// Number of provider successes.
	responseStatusPerProvider *prometheus.CounterVec

	// Last time a given ID (i.e. currency pair) was updated.
	lastUpdatedPerProvider *prometheus.GaugeVec
}

// NewProviderMetricsFromConfig returns a new Metrics struct given the main oracle metrics config.
func NewProviderMetricsFromConfig(config config.MetricsConfig) ProviderMetrics {
	if config.Enabled {
		return NewProviderMetrics()
	}
	return NewNopProviderMetrics()
}

// NewProviderMetrics returns a Provider Metrics implementation that uses Prometheus.
func NewProviderMetrics() ProviderMetrics {
	m := &ProviderMetricsImpl{
		responseStatusPerProviderByID: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "provider_status_responses_per_id",
			Help:      "Number of provider successes with a given ID.",
		}, []string{ProviderLabel, IDLabel, StatusLabel, ProviderTypeLabel}),
		responseStatusPerProvider: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "provider_status_responses",
			Help:      "Number of provider successes.",
		}, []string{ProviderLabel, StatusLabel, ProviderTypeLabel}),
		lastUpdatedPerProvider: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "provider_last_updated_id",
			Help:      "Last time a given ID (i.e. currency pair) was updated.",
		}, []string{ProviderLabel, IDLabel, ProviderTypeLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.responseStatusPerProviderByID)
	prometheus.MustRegister(m.responseStatusPerProvider)
	prometheus.MustRegister(m.lastUpdatedPerProvider)

	return m
}

type noOpProviderMetricsImpl struct{}

// NewNopProviderMetrics returns a Provider Metrics implementation that does not collect metrics.
func NewNopProviderMetrics() ProviderMetrics {
	return &noOpProviderMetricsImpl{}
}

func (m *noOpProviderMetricsImpl) AddProviderResponseByID(_, _ string, _ Status, _ providertypes.ProviderType) {
}

func (m *noOpProviderMetricsImpl) AddProviderResponse(_ string, _ Status, _ providertypes.ProviderType) {
}
func (m *noOpProviderMetricsImpl) LastUpdated(_, _ string, _ providertypes.ProviderType) {}

// AddProviderResponseByID increments the number of ticks with a fully successful provider update
// for a given provider and ID (i.e. currency pair).
func (m *ProviderMetricsImpl) AddProviderResponseByID(providerName, id string, status Status, providerType providertypes.ProviderType) {
	m.responseStatusPerProviderByID.With(prometheus.Labels{
		ProviderLabel:     providerName,
		IDLabel:           id,
		StatusLabel:       string(status),
		ProviderTypeLabel: string(providerType),
	},
	).Add(1)
}

// AddProviderResponse increments the number of ticks with a fully successful provider update.
func (m *ProviderMetricsImpl) AddProviderResponse(providerName string, status Status, providerType providertypes.ProviderType) {
	m.responseStatusPerProvider.With(prometheus.Labels{
		ProviderLabel:     providerName,
		StatusLabel:       string(status),
		ProviderTypeLabel: string(providerType),
	},
	).Add(1)
}

// LastUpdated updates the last time a given ID (i.e. currency pair) was updated.
func (m *ProviderMetricsImpl) LastUpdated(providerName, id string, providerType providertypes.ProviderType) {
	now := time.Now().UTC()
	m.lastUpdatedPerProvider.With(prometheus.Labels{
		ProviderLabel:     providerName,
		IDLabel:           id,
		ProviderTypeLabel: string(providerType),
	},
	).Set(float64(now.Unix()))
}
