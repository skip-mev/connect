package metrics

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
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
	promMetrics  *PromMetrics
	statsdClient *statsd.Client
}

type PromMetrics struct {
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
		var telemetryPushAddress string
		if !config.Telemetry.Disabled {
			telemetryPushAddress = config.Telemetry.PushAddress
		}
		return NewMetrics(telemetryPushAddress)
	}
	return NewNopMetrics()
}

func NewNopMetrics() Metrics {
	return &OracleMetricsImpl{}
}

// NewMetrics returns a Metrics implementation that exposes metrics to Prometheus.
func NewMetrics(telemetryPushAddress string) Metrics {
	var statsdClient *statsd.Client
	if telemetryPushAddress != "" {
		statsdClient, _ = statsd.New(telemetryPushAddress)
	}
	promMetrics := &PromMetrics{
		ticks: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      TicksMetricName,
			Help:      "Number of ticks with a successful oracle update.",
		}),
		tickerTicks: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      TickerTicksMetricName,
			Help:      "Number of ticks with a successful ticker update.",
		}, []string{PairIDLabel}),
		prices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      PricesMetricName,
			Help:      "Price gauge for a given currency pair on a provider",
		}, []string{ProviderLabel, PairIDLabel, DecimalsLabel}),
		aggregatePrices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      AggregatePricesMetricName,
			Help:      "Aggregate price for a given currency pair",
		}, []string{PairIDLabel, DecimalsLabel}),
		providerTick: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: OracleSubsystem,
			Name:      ProviderTickMetricName,
			Help:      "Number of ticks with a successful provider update.",
		}, []string{ProviderLabel, PairIDLabel, SuccessLabel}),
		providerCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ProviderCountMetricName,
			Name:      ProviderTickMetricName,
			Help:      "Number of providers that were utilized to calculate the final price for a given market.",
		}, []string{PairIDLabel}),
		slinkyBuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: OracleSubsystem,
			Name:      SlinkyBuildInfoMetricName,
			Help:      "Information about the slinky build",
		}, []string{Version}),
	}

	prometheus.MustRegister(promMetrics.ticks)
	prometheus.MustRegister(promMetrics.tickerTicks)
	prometheus.MustRegister(promMetrics.prices)
	prometheus.MustRegister(promMetrics.aggregatePrices)
	prometheus.MustRegister(promMetrics.providerTick)
	prometheus.MustRegister(promMetrics.providerCount)
	prometheus.MustRegister(promMetrics.slinkyBuildInfo)

	return &OracleMetricsImpl{
		promMetrics,
		statsdClient,
	}
}

// AddTick increments the total number of ticks that have been processed by the oracle.
func (m *OracleMetricsImpl) AddTick() {
	if m.promMetrics != nil {
		m.promMetrics.ticks.Add(1)
	}
	if m.statsdClient != nil {
		m.statsdClient.Incr(TicksMetricName, []string{}, 1)
	}
}

// AddTickerTick increments the number of ticks for a given ticker. Specifically, this
// is used to track the number of times a ticker was updated.
func (m *OracleMetricsImpl) AddTickerTick(ticker string) {
	if m.promMetrics != nil {
		m.promMetrics.tickerTicks.With(prometheus.Labels{
			PairIDLabel: strings.ToLower(ticker),
		},
		).Add(1)
	}

	if m.statsdClient != nil {
		m.statsdClient.Incr(TickerTicksMetricName, []string{strings.ToLower(ticker)}, 1)
	}
}

// UpdatePrice price updates the price for the given pairID for the provider.
func (m *OracleMetricsImpl) UpdatePrice(
	providerName, pairID string,
	decimals uint64,
	price float64,
) {
	if m.promMetrics != nil {
		m.promMetrics.prices.With(prometheus.Labels{
			ProviderLabel: strings.ToLower(providerName),
			PairIDLabel:   strings.ToLower(pairID),
			DecimalsLabel: fmt.Sprintf("%d", decimals),
		},
		).Set(price)
	}

	if m.statsdClient != nil {
		metricName := strings.Join([]string{PricesMetricName, strings.ToLower(providerName), strings.ToLower(pairID)}, ".")
		m.statsdClient.Gauge(metricName, price, []string{fmt.Sprintf("%d", decimals)}, 1)
	}
}

// UpdateAggregatePrice updates the aggregated price for the given pairID.
func (m *OracleMetricsImpl) UpdateAggregatePrice(
	pairID string,
	decimals uint64,
	price float64,
) {
	if m.promMetrics != nil {
		m.promMetrics.aggregatePrices.With(prometheus.Labels{
			PairIDLabel:   strings.ToLower(pairID),
			DecimalsLabel: fmt.Sprintf("%d", decimals),
		},
		).Set(price)
	}

	if m.statsdClient != nil {
		metricName := strings.Join([]string{AggregatePricesMetricName, strings.ToLower(PairIDLabel)}, ".")
		m.statsdClient.Gauge(metricName, price, []string{fmt.Sprintf("%d", decimals)}, 1)
	}
}

// AddProviderTick increments the number of ticks for a given provider. Specifically,
// this is used to track the number of times a provider included a price update that
// was used in the aggregation.
func (m *OracleMetricsImpl) AddProviderTick(providerName, pairID string, success bool) {
	if m.promMetrics != nil {
		m.promMetrics.providerTick.With(prometheus.Labels{
			ProviderLabel: strings.ToLower(providerName),
			PairIDLabel:   strings.ToLower(pairID),
			SuccessLabel:  fmt.Sprintf("%t", success),
		},
		).Add(1)
	}

	if m.statsdClient != nil {
		metricName := strings.Join([]string{ProviderTickMetricName, strings.ToLower(ProviderLabel), strings.ToLower(PairIDLabel)}, ".")
		m.statsdClient.Incr(metricName, []string{fmt.Sprintf("%t", success)}, 1)
	}
}

// AddProviderCountForMarket increments the number of providers that were utilized
// to calculate the final price for a given market.
func (m *OracleMetricsImpl) AddProviderCountForMarket(market string, count int) {
	if m.promMetrics != nil {
		m.promMetrics.providerCount.With(prometheus.Labels{
			PairIDLabel: strings.ToLower(market),
		},
		).Set(float64(count))
	}

	if m.statsdClient != nil {
		metricName := strings.Join([]string{ProviderCountMetricName, strings.ToLower(PairIDLabel)}, ".")
		m.statsdClient.Gauge(metricName, float64(count), []string{}, 1)
	}
}

// SetSlinkyBuildInfo sets the build information for the Slinky binary. The version exported
// is determined by the build time version in accordance with the build pkg.
func (m *OracleMetricsImpl) SetSlinkyBuildInfo() {
	if m.promMetrics != nil {
		m.promMetrics.slinkyBuildInfo.With(prometheus.Labels{
			Version: build.Build,
		}).Set(1)
	}

	if m.statsdClient != nil {
		metricName := strings.Join([]string{SlinkyBuildInfoMetricName, strings.ToLower(build.Build)}, ".")
		m.statsdClient.Gauge(metricName, float64(1), []string{}, 1)
	}
}
