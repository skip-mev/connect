package metrics

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// ProviderLabel is a label for the provider name.
	ProviderLabel = "provider"
	// PairIDLabel is the currency pair for which the metric applies.
	PairIDLabel = "id"
	// DecimalsLabel is the number of decimal points associated with the price.
	DecimalsLabel = "decimals"
	// OracleSubsystem is a subsystem shared by all metrics exposed by this package.
	OracleSubsystem = "side_car"
)

// Metrics is an interface that defines the API for oracle metrics.
//
//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// AddTick increments the number of ticks, this can represent a liveness counter. This
	// is incremented once every interval (which is defined by the oracle config).
	AddTick()

	// AddTickerTick increments the number of ticks for a given ticker. Specifically, this
	// is used to track the number of times a ticker was updated.
	AddTickerTick(ticker string)

	// UpdatePrice price updates the price for the given pairID for the provider.
	UpdatePrice(name, pairID string, decimals uint64, price float64)

	// UpdateAggregatePrice updates the aggregated price for the given pairID.
	UpdateAggregatePrice(pairID string, decimals uint64, price float64)

	// AddProviderTick increments the number of ticks for a given provider. Specifically,
	// this is used to track the number of times a provider included a price update that
	// was used in the aggregation.
	AddProviderTick(providerName, pairID string)
}

// OracleMetricsImpl is a Metrics implementation that does nothing.
type OracleMetricsImpl struct {
	ticks           prometheus.Counter
	tickerTicks     *prometheus.CounterVec
	prices          *prometheus.GaugeVec
	aggregatePrices *prometheus.GaugeVec
	providerTick    *prometheus.CounterVec
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
			Name:      "health_check_system",
			Help:      "Number of ticks with a successful oracle update.",
		}),
		tickerTicks: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "health_check_ticker",
			Help:      "Number of ticks with a successful ticker update.",
		}, []string{PairIDLabel}),
		prices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "provider_price",
			Help:      "Price gauge for a given currency pair on a provider",
		}, []string{ProviderLabel, PairIDLabel, DecimalsLabel}),
		aggregatePrices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "aggregated_price",
			Help:      "Aggregate price for a given currency pair",
		}, []string{PairIDLabel, DecimalsLabel}),
		providerTick: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "health_check_provider",
			Help:      "Number of ticks with a successful provider update.",
		}, []string{ProviderLabel, PairIDLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.ticks)
	prometheus.MustRegister(m.tickerTicks)
	prometheus.MustRegister(m.prices)
	prometheus.MustRegister(m.aggregatePrices)
	prometheus.MustRegister(m.providerTick)

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

// AddTickerTick increments the number of ticks for a given ticker. Specifically, this
// is used to track the number of times a ticker was updated.
func (m *noOpOracleMetrics) AddTickerTick(_ string) {
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *noOpOracleMetrics) UpdatePrice(_, _ string, _ uint64, _ float64) {
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *noOpOracleMetrics) UpdateAggregatePrice(string, uint64, float64) {
}

// AddProviderTick increments the number of ticks for a given provider. Specifically,
// this is used to track the number of times a provider included a price update that
// was used in the aggregation.
func (m *noOpOracleMetrics) AddProviderTick(providerName, pairID string) {
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *OracleMetricsImpl) AddTick() {
	m.ticks.Add(1)
}

// AddTickerTick increments the number of ticks for a given ticker. Specifically, this
// is used to track the number of times a ticker was updated.
func (m *OracleMetricsImpl) AddTickerTick(ticker string) {
	m.tickerTicks.With(prometheus.Labels{
		PairIDLabel: strings.ToLower(ticker),
	},
	).Add(1)
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *OracleMetricsImpl) UpdatePrice(
	providerName, pairID string,
	decimals uint64,
	price float64,
) {
	m.prices.With(prometheus.Labels{
		ProviderLabel: strings.ToLower(providerName),
		PairIDLabel:   strings.ToLower(pairID),
		DecimalsLabel: fmt.Sprintf("%d", decimals),
	},
	).Set(price)
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *OracleMetricsImpl) UpdateAggregatePrice(
	pairID string,
	decimals uint64,
	price float64,
) {
	m.aggregatePrices.With(prometheus.Labels{
		PairIDLabel:   strings.ToLower(pairID),
		DecimalsLabel: fmt.Sprintf("%d", decimals),
	},
	).Set(price)
}

// AddProviderTick increments the number of ticks for a given provider. Specifically,
// this is used to track the number of times a provider included a price update that
// was used in the aggregation.
func (m *OracleMetricsImpl) AddProviderTick(providerName, pairID string) {
	m.providerTick.With(prometheus.Labels{
		ProviderLabel: strings.ToLower(providerName),
		PairIDLabel:   strings.ToLower(pairID),
	},
	).Add(1)
}
