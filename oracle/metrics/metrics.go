package metrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"

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
	//
	// TODO: Add more metrics here in later PRs.

	// UpdatePrice price updates the price for the given pairID for the provider.
	UpdatePrice(name, handlerType, pairID string, price float64)

	// UpdateAggregatePrice rice updates the aggregated price for the given pairID.
	UpdateAggregatePrice(pairID string, price float64)
}

// OracleMetricsImpl is a Metrics implementation that does nothing.
type OracleMetricsImpl struct {
	ticks           metrics.Counter
	prices          metrics.Gauge
	aggregatePrices metrics.Gauge
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
		ticks: prometheus.NewCounterFrom(stdprom.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "ticks",
			Help:      "Number of ticks with a successful oracle update.",
		}, []string{}),
		prices: prometheus.NewGaugeFrom(stdprom.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "provider_price",
			Help:      "Price gauge for a given currency pair on a provider",
		}, []string{ProviderLabel, ProviderTypeLabel, PairIDLabel}),
		aggregatePrices: prometheus.NewGaugeFrom(stdprom.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "aggregate_price",
			Help:      "Aggregate price for a given currency pair",
		}, []string{PairIDLabel}),
	}

	return m
}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &OracleMetricsImpl{
		ticks:           discard.NewCounter(),
		prices:          discard.NewGauge(),
		aggregatePrices: discard.NewGauge(),
	}
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *OracleMetricsImpl) AddTick() {
	m.ticks.Add(1)
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *OracleMetricsImpl) UpdatePrice(providerName, handlerType, pairID string, price float64) {
	m.prices.With(
		ProviderLabel, providerName,
		ProviderTypeLabel, handlerType,
		PairIDLabel, pairID,
	).Add(price)
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *OracleMetricsImpl) UpdateAggregatePrice(pairID string, price float64) {
	m.aggregatePrices.With(
		PairIDLabel, pairID,
	).Add(price)
}
