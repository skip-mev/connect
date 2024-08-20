package metrics

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/connect/v2/cmd/build"
	"github.com/skip-mev/connect/v2/oracle/config"
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

	TicksMetricName           = "health_check_system_updates_total"
	TickerTicksMetricName     = "health_check_ticker_updates_total"
	PricesMetricName          = "provider_price"
	AggregatePricesMetricName = "aggregated_price"
	ProviderTickMetricName    = "health_check_provider_updates_total"
	ProviderCountMetricName   = "health_check_market_providers"
	SlinkyBuildInfoMetricName = "slinky_build_info"
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
	promTicks           prometheus.Counter
	promTickerTicks     *prometheus.CounterVec
	promPrices          *prometheus.GaugeVec
	promAggregatePrices *prometheus.GaugeVec
	promProviderTick    *prometheus.CounterVec
	promProviderCount   *prometheus.GaugeVec
	promSlinkyBuildInfo *prometheus.GaugeVec
	statsdClient        statsd.ClientInterface
}

// NewMetricsFromConfig returns an oracle Metrics implementation based on the provided
// config.
func NewMetricsFromConfig(config config.MetricsConfig) Metrics {
	if config.Enabled {
		var telemetryPushAddress string
		if !config.Telemetry.Disabled {
			telemetryPushAddress = config.Telemetry.PushAddress
		}
		return NewMetrics(telemetryPushAddress)
	}
	return NewNopMetrics()
}

// NewMetrics returns a Metrics implementation that exposes metrics to Prometheus.
func NewMetrics(telemetryPushAddress string) Metrics {
	ret := OracleMetricsImpl{}

	if telemetryPushAddress != "" {
		ret.statsdClient, _ = statsd.New(telemetryPushAddress)
	} else {
		ret.statsdClient = &statsd.NoOpClient{}
	}

	ret.promTicks = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: OracleSubsystem,
		Name:      "health_check_system_updates_total",
		Help:      "Number of ticks with a successful oracle update.",
	})
	ret.promTickerTicks = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      "health_check_ticker_updates_total",
			Help:      "Number of ticks with a successful ticker update.",
	}, []string{PairIDLabel})

	ret.promPrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      "provider_price",
		Help:      "Price gauge for a given currency pair on a provider",
	}, []string{ProviderLabel, PairIDLabel, DecimalsLabel})
	ret.promAggregatePrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      "aggregated_price",
		Help:      "Aggregate price for a given currency pair",
	}, []string{PairIDLabel, DecimalsLabel})
	ret.promProviderTick = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: OracleSubsystem,
		Name:      "health_check_provider_updates_total",
		Help:      "Number of ticks with a successful provider update.",
	}, []string{ProviderLabel, PairIDLabel, SuccessLabel})
	ret.promProviderCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      "health_check_market_providers",
		Help:      "Number of providers that were utilized to calculate the final price for a given market.",
	}, []string{PairIDLabel})
	ret.promSlinkyBuildInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      "slinky_build_info",
		Help:      "Information about the slinky build",
	}, []string{Version})

	prometheus.MustRegister(ret.promTicks)
	prometheus.MustRegister(ret.promTickerTicks)
	prometheus.MustRegister(ret.promPrices)
	prometheus.MustRegister(ret.promAggregatePrices)
	prometheus.MustRegister(ret.promProviderTick)
	prometheus.MustRegister(ret.promProviderCount)
	prometheus.MustRegister(ret.promSlinkyBuildInfo)

	return &ret
}

type noOpOracleMetrics struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &noOpOracleMetrics{}
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *noOpOracleMetrics) AddTick() {}

// AddTickerTick increments the number of ticks for a given ticker. Specifically, this
// is used to track the number of times a ticker was updated.
func (m *noOpOracleMetrics) AddTickerTick(_ string) {}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *noOpOracleMetrics) UpdatePrice(_, _ string, _ uint64, _ float64) {}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *noOpOracleMetrics) UpdateAggregatePrice(string, uint64, float64) {}

// AddProviderTick increments the number of ticks for a given provider. Specifically,
// this is used to track the number of times a provider included a price update that
// was used in the aggregation.
func (m *noOpOracleMetrics) AddProviderTick(_, _ string, _ bool) {}

// AddProviderCountForMarket increments the number of providers that were utilized
// to calculate the final price for a given market.
func (m *noOpOracleMetrics) AddProviderCountForMarket(string, int) {}

// SetSlinkyBuildInfo sets the build information for the Slinky binary.
func (m *noOpOracleMetrics) SetSlinkyBuildInfo() {}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *OracleMetricsImpl) AddTick() {
	m.promTicks.Add(1)
	m.statsdClient.Incr(TicksMetricName, []string{}, 1)
}

// AddTickerTick increments the number of ticks for a given ticker. Specifically, this
// is used to track the number of times a ticker was updated.
func (m *OracleMetricsImpl) AddTickerTick(ticker string) {
	m.promTickerTicks.With(prometheus.Labels{
		PairIDLabel: strings.ToLower(ticker),
	},
	).Add(1)

	m.statsdClient.Incr(TickerTicksMetricName, []string{strings.ToLower(ticker)}, 1)
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *OracleMetricsImpl) UpdatePrice(
	providerName, pairID string,
	decimals uint64,
	price float64,
) {
	m.promPrices.With(prometheus.Labels{
		ProviderLabel: strings.ToLower(providerName),
		PairIDLabel:   strings.ToLower(pairID),
		DecimalsLabel: fmt.Sprintf("%d", decimals),
	},
	).Set(price)

	metricName := strings.Join([]string{PricesMetricName, strings.ToLower(providerName), strings.ToLower(pairID)}, ".")
	m.statsdClient.Gauge(metricName, price, []string{fmt.Sprintf("%d", decimals)}, 1)
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *OracleMetricsImpl) UpdateAggregatePrice(
	pairID string,
	decimals uint64,
	price float64,
) {
	m.promAggregatePrices.With(prometheus.Labels{
		PairIDLabel:   strings.ToLower(pairID),
		DecimalsLabel: fmt.Sprintf("%d", decimals),
	},
	).Set(price)

	metricName := strings.Join([]string{AggregatePricesMetricName, strings.ToLower(PairIDLabel)}, ".")
	m.statsdClient.Gauge(metricName, price, []string{fmt.Sprintf("%d", decimals)}, 1)
}

// AddProviderTick increments the number of ticks for a given provider. Specifically,
// this is used to track the number of times a provider included a price update that
// was used in the aggregation.
func (m *OracleMetricsImpl) AddProviderTick(providerName, pairID string, success bool) {
	m.promProviderTick.With(prometheus.Labels{
		ProviderLabel: strings.ToLower(providerName),
		PairIDLabel:   strings.ToLower(pairID),
		SuccessLabel:  fmt.Sprintf("%t", success),
	},
	).Add(1)

	metricName := strings.Join([]string{ProviderTickMetricName, strings.ToLower(ProviderLabel), strings.ToLower(PairIDLabel)}, ".")
	m.statsdClient.Incr(metricName, []string{fmt.Sprintf("%t", success)}, 1)
}

// AddProviderCountForMarket increments the number of providers that were utilized
// to calculate the final price for a given market.
func (m *OracleMetricsImpl) AddProviderCountForMarket(market string, count int) {
	m.promProviderCount.With(prometheus.Labels{
		PairIDLabel: strings.ToLower(market),
	},
	).Set(float64(count))

	metricName := strings.Join([]string{ProviderCountMetricName, strings.ToLower(PairIDLabel)}, ".")
	m.statsdClient.Gauge(metricName, float64(count), []string{}, 1)
}

// SetSlinkyBuildInfo sets the build information for the Slinky binary. The version exported
// is determined by the build time version in accordance with the build pkg.
func (m *OracleMetricsImpl) SetSlinkyBuildInfo() {
	m.promSlinkyBuildInfo.With(prometheus.Labels{
		Version: build.Build,
	}).Set(1)

	metricName := strings.Join([]string{SlinkyBuildInfoMetricName, strings.ToLower(build.Build)}, ".")
	m.statsdClient.Gauge(metricName, float64(1), []string{}, 1)
}
