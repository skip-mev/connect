package metrics

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/cmd/build"
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
	// SuccessLabel is a label for a successful operation.
	SuccessLabel = "success"
	// Version is a label for the Slinky version.
	Version = "version"
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
	AddProviderTick(providerName, pairID string, success bool)

	// AddProviderCountForMarket increments the number of providers that were utilized
	// to calculate the final price for a given market.
	AddProviderCountForMarket(market string, count int)

	// SetSlinkyBuildInfo sets the build information for the Slinky binary.
	SetSlinkyBuildInfo()
}

// OracleMetricsImpl is a Metrics implementation that does nothing.
type OracleMetricsImpl struct {
	ticks           prometheus.Counter
	tickerTicks     *prometheus.CounterVec
	prices          *prometheus.GaugeVec
	aggregatePrices *prometheus.GaugeVec
	providerTick    *prometheus.CounterVec
	providerCount   *prometheus.GaugeVec
	slinkyBuildInfo *prometheus.GaugeVec
}

// NewMetricsFromConfig returns an oracle Metrics implementation based on the provided
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
			Name:      "health_check_system_updates_total",
			Help:      "Number of ticks with a successful oracle update.",
		}),
		tickerTicks: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "health_check_ticker_updates_total",
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
			Name:      "health_check_provider_updates_total",
			Help:      "Number of ticks with a successful provider update.",
		}, []string{ProviderLabel, PairIDLabel, SuccessLabel}),
		providerCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "health_check_market_providers",
			Help:      "Number of providers that were utilized to calculate the final price for a given market.",
		}, []string{PairIDLabel}),
		slinkyBuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      "slinky_build_info",
			Help:      "Information about the slinky build",
		}, []string{Version}),
	}

	// register the above metrics
	prometheus.MustRegister(m.ticks)
	prometheus.MustRegister(m.tickerTicks)
	prometheus.MustRegister(m.prices)
	prometheus.MustRegister(m.aggregatePrices)
	prometheus.MustRegister(m.providerTick)
	prometheus.MustRegister(m.providerCount)
	prometheus.MustRegister(m.slinkyBuildInfo)

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
func (m *noOpOracleMetrics) AddProviderTick(_, _ string, _ bool) {
}

// AddProviderCountForMarket increments the number of providers that were utilized
// to calculate the final price for a given market.
func (m *noOpOracleMetrics) AddProviderCountForMarket(string, int) {
}

// SetSlinkyBuildInfo sets the build information for the Slinky binary.
func (m *noOpOracleMetrics) SetSlinkyBuildInfo() {}

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
func (m *OracleMetricsImpl) AddProviderTick(providerName, pairID string, success bool) {
	m.providerTick.With(prometheus.Labels{
		ProviderLabel: strings.ToLower(providerName),
		PairIDLabel:   strings.ToLower(pairID),
		SuccessLabel:  fmt.Sprintf("%t", success),
	},
	).Add(1)
}

// AddProviderCountForMarket increments the number of providers that were utilized
// to calculate the final price for a given market.
func (m *OracleMetricsImpl) AddProviderCountForMarket(market string, count int) {
	m.providerCount.With(prometheus.Labels{
		PairIDLabel: strings.ToLower(market),
	},
	).Set(float64(count))
}

// SetSlinkyBuildInfo sets the build information for the Slinky binary. The version exported
// is determined by the build time version in accordance with the build pkg.
func (m *OracleMetricsImpl) SetSlinkyBuildInfo() {
	m.slinkyBuildInfo.With(prometheus.Labels{
		Version: build.Build,
	}).Set(1)
}
