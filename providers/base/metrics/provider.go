package metrics

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
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
	Status       string
	ProviderType string
)

const (
	Success    Status       = "success"
	Failure    Status       = "failure"
	WebSockets ProviderType = "websockets"
	API        ProviderType = "api"
)

// ProviderMetrics is an interface that defines the API for metrics collection for providers. The
// base provider utilizes this interface to collect metrics, whether the underlying implementation
// is API or web socket based.
//
//go:generate mockery --name ProviderMetrics --filename mock_metrics.go
type ProviderMetrics interface {
	// AddProviderResponseByID increments the number of ticks with a fully successful provider update
	// for a given provider and ID (i.e. currency pair).
	AddProviderResponseByID(providerName, id string, status Status, providerType ProviderType)

	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	AddProviderResponse(providerName string, status Status, providerType ProviderType)

	// LastUpdated updates the last time a given ID (i.e. currency pair) was updated.
	LastUpdated(providerName, id string, providerType ProviderType)
}

// ProviderMetricsImpl contains metrics exposed by this package.
type ProviderMetricsImpl struct {
	// Number of provider successes by ID.
	responseStatusPerProviderByID metrics.Counter

	// Number of provider successes.
	responseStatusPerProvider metrics.Counter

	// Last time a given ID (i.e. currency pair) was updated.
	lastUpdatedPerProvider metrics.Gauge
}

// NewProviderMetricsFromConfig returns a new Metrics struct given the main oracle metrics config.
func NewProviderMetricsFromConfig(config oracleconfig.OracleMetricsConfig) ProviderMetrics {
	if config.Enabled {
		return NewProviderMetrics()
	}
	return NewNopProviderMetrics()
}

// NewProviderMetrics returns a Provider Metrics implementation that uses Prometheus.
func NewProviderMetrics() ProviderMetrics {
	m := &ProviderMetricsImpl{
		responseStatusPerProviderByID: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "provider_status_responses_per_id",
			Help:      "Number of provider successes with a given ID.",
		}, []string{ProviderLabel, IDLabel, StatusLabel, ProviderTypeLabel}),
		responseStatusPerProvider: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "provider_status_responses",
			Help:      "Number of provider successes.",
		}, []string{ProviderLabel, StatusLabel, ProviderTypeLabel}),
		lastUpdatedPerProvider: prometheus.NewGaugeFrom(stdprom.GaugeOpts{
			Namespace: oraclemetrics.OracleSubsystem,
			Name:      "provider_last_updated_id",
			Help:      "Last time a given ID (i.e. currency pair) was updated.",
		}, []string{ProviderLabel, IDLabel, ProviderTypeLabel}),
	}

	return m
}

// NewNopProviderMetrics returns a Provider Metrics implementation that does not collect metrics.
func NewNopProviderMetrics() ProviderMetrics {
	return &ProviderMetricsImpl{
		responseStatusPerProviderByID: discard.NewCounter(),
		responseStatusPerProvider:     discard.NewCounter(),
		lastUpdatedPerProvider:        discard.NewGauge(),
	}
}

// AddProviderResponseByID increments the number of ticks with a fully successful provider update
// for a given provider and ID (i.e. currency pair).
func (m *ProviderMetricsImpl) AddProviderResponseByID(providerName, id string, status Status, providerType ProviderType) {
	m.responseStatusPerProviderByID.With(
		ProviderLabel, providerName,
		IDLabel, id,
		StatusLabel, string(status),
		ProviderTypeLabel, string(providerType),
	).Add(1)
}

// AddProviderResponse increments the number of ticks with a fully successful provider update.
func (m *ProviderMetricsImpl) AddProviderResponse(providerName string, status Status, providerType ProviderType) {
	m.responseStatusPerProvider.With(
		ProviderLabel, providerName,
		StatusLabel, string(status),
		ProviderTypeLabel, string(providerType),
	).Add(1)
}

// LastUpdated updates the last time a given ID (i.e. currency pair) was updated.
func (m *ProviderMetricsImpl) LastUpdated(providerName, id string, providerType ProviderType) {
	now := time.Now().UTC()
	m.lastUpdatedPerProvider.With(
		ProviderLabel, providerName,
		IDLabel, id,
		ProviderTypeLabel, string(providerType),
	).Set(float64(now.Unix()))
}
