package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// ProviderLabel is a label for the provider name.
	ProviderLabel = "provider"
	// ProviderTypeLabel is a label for the type of provider (WS, API, etc.)
	ProviderTypeLabel = "type"
	// PairIDLabel is the
	PairIDLabel = "pair"
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

	// UpdatePrice price updates the price for the given pairID for the provider.
	UpdatePrice(name, handlerType, pairID string, price float64)

	// UpdateAggregatePrice updates the aggregated price for the given pairID.
	UpdateAggregatePrice(pairID string, price float64)
}

// OracleMetricsImpl is a Metrics implementation that does nothing.
type OracleMetricsImpl struct {
	ticks           prometheus.Counter
	prices          *prometheus.GaugeVec
	aggregatePrices *prometheus.GaugeVec
}

// NewMetricsFromConfig returns a oracle Metrics implementation based on the provided
// config.
func NewMetricsFromConfig(config config.MetricsConfig) Metrics {
	if config.Enabled {
		return NewMetrics()
	}
	return NewNopMetrics()
}

// NewMetrics returns a Metrics implementation that exposes metrics to Prometheus.
func NewMetrics() Metrics {
	m := &OracleMetricsImpl{
		ticks: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "ticks",
			Help:      "Number of ticks with a successful oracle update.",
		}),
		prices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "provider_price",
			Help:      "Price gauge for a given currency pair on a provider",
		}, []string{ProviderLabel, ProviderTypeLabel, PairIDLabel}),
		aggregatePrices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "aggregate_price",
			Help:      "Aggregate price for a given currency pair",
		}, []string{PairIDLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.ticks)
	prometheus.MustRegister(m.prices)
	prometheus.MustRegister(m.aggregatePrices)

	return m
}

type noOpOracleMetrics struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &noOpOracleMetrics{}
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *noOpOracleMetrics) AddTick() {
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *noOpOracleMetrics) UpdatePrice(_, _, _ string, _ float64) {
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *noOpOracleMetrics) UpdateAggregatePrice(_ string, _ float64) {
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *OracleMetricsImpl) AddTick() {
	m.ticks.Add(1)
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *OracleMetricsImpl) UpdatePrice(providerName, handlerType, pairID string, price float64) {
	m.prices.With(prometheus.Labels{
		ProviderLabel:     providerName,
		ProviderTypeLabel: handlerType,
		PairIDLabel:       pairID,
	},
	).Set(price)
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *OracleMetricsImpl) UpdateAggregatePrice(pairID string, price float64) {
	m.aggregatePrices.With(prometheus.Labels{
		PairIDLabel: pairID,
	},
	).Set(price)
}
