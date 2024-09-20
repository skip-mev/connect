package metrics

import (
	"fmt"
	"strings"
	"sync"

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
	// Version is a label for the Connect version.
	Version = "version"

	TicksMetricName            = "health_check_system_updates_total"
	TickerTicksMetricName      = "health_check_ticker_updates_total"
	PricesMetricName           = "provider_price"
	AggregatePricesMetricName  = "aggregated_price"
	ProviderTickMetricName     = "health_check_provider_updates_total"
	ProviderCountMetricName    = "health_check_market_providers"
	ConnectBuildInfoMetricName = "connect_build_info"
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
	AddTickerTick(pairID string)

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
	AddProviderCountForMarket(pairID string, count int)

	// SetConnectBuildInfo sets the build information for the Connect binary.
	SetConnectBuildInfo()

	// MissingPrices sets a list of missing prices for the given aggregation tick.
	MissingPrices(pairIDs []string)

	// GetMissingPrices gets the current list of missing prices.
	GetMissingPrices() []string
}

// OracleMetricsImpl is a Metrics implementation that does nothing.
type OracleMetricsImpl struct {
	promTicks             prometheus.Counter
	promTickerTicks       *prometheus.CounterVec
	promPrices            *prometheus.GaugeVec
	promAggregatePrices   *prometheus.GaugeVec
	promProviderTick      *prometheus.CounterVec
	promProviderCount     *prometheus.GaugeVec
	promConnectBuildInfo  *prometheus.GaugeVec
	statsdClient          statsd.ClientInterface
	nodeIdentifier        string
	missingPricesInternal []string
	missingPricesMtx      sync.Mutex
}

// NewMetricsFromConfig returns an oracle Metrics implementation based on the provided
// config.
func NewMetricsFromConfig(config config.MetricsConfig, nodeClient NodeClient) Metrics {
	if config.Enabled {
		var err error

		var statsdClient statsd.ClientInterface = &statsd.NoOpClient{}
		identifier := ""
		if !config.Telemetry.Disabled && nodeClient != nil {
			// Group these metrics into a statsd namespace
			identifier, err = nodeClient.DeriveNodeIdentifier()
			if err == nil { // only publish statsd data when connected to a node
				c, err := statsd.New(config.Telemetry.PushAddress, func(c *statsd.Options) error {
					// Prepends all messages with connect.sidecar.
					c.Namespace = "connect.sidecar."
					c.Tags = []string{identifier}
					return nil
				})
				if err == nil {
					statsdClient = c
				}
			}
		}
		return NewMetrics(statsdClient, identifier)
	}
	return NewNopMetrics()
}

// NewMetrics returns a Metrics implementation that exposes metrics to Prometheus.
func NewMetrics(statsdClient statsd.ClientInterface, nodeIdentifier string) Metrics {
	ret := OracleMetricsImpl{
		missingPricesInternal: make([]string, 0),
	}

	ret.statsdClient = statsdClient
	ret.nodeIdentifier = nodeIdentifier

	ret.promTicks = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: OracleSubsystem,
		Name:      TicksMetricName,
		Help:      "Number of ticks with a successful oracle update.",
	})
	ret.promTickerTicks = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: OracleSubsystem,
		Name:      TickerTicksMetricName,
		Help:      "Number of ticks with a successful ticker update.",
	}, []string{PairIDLabel})

	ret.promPrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      PricesMetricName,
		Help:      "Price gauge for a given currency pair on a provider",
	}, []string{ProviderLabel, PairIDLabel, DecimalsLabel})
	ret.promAggregatePrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      AggregatePricesMetricName,
		Help:      "Aggregate price for a given currency pair",
	}, []string{PairIDLabel, DecimalsLabel})
	ret.promProviderTick = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: OracleSubsystem,
		Name:      ProviderTickMetricName,
		Help:      "Number of ticks with a successful provider update.",
	}, []string{ProviderLabel, PairIDLabel, SuccessLabel})
	ret.promProviderCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      ProviderCountMetricName,
		Help:      "Number of providers that were utilized to calculate the final price for a given market.",
	}, []string{PairIDLabel})
	ret.promConnectBuildInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: OracleSubsystem,
		Name:      ConnectBuildInfoMetricName,
		Help:      "Information about the connect build",
	}, []string{Version})

	prometheus.MustRegister(ret.promTicks)
	prometheus.MustRegister(ret.promTickerTicks)
	prometheus.MustRegister(ret.promPrices)
	prometheus.MustRegister(ret.promAggregatePrices)
	prometheus.MustRegister(ret.promProviderTick)
	prometheus.MustRegister(ret.promProviderCount)
	prometheus.MustRegister(ret.promConnectBuildInfo)

	return &ret
}

type noOpOracleMetrics struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &noOpOracleMetrics{}
}

func (m *noOpOracleMetrics) MissingPrices(_ []string) {}

func (m *noOpOracleMetrics) GetMissingPrices() []string { return []string{} }

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

// SetConnectBuildInfo sets the build information for the Connect binary.
func (m *noOpOracleMetrics) SetConnectBuildInfo() {}

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

	metricName := strings.Join([]string{TickerTicksMetricName, m.nodeIdentifier, strings.ToLower(ticker)}, ".")
	m.statsdClient.Incr(metricName, []string{}, 1)
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

	metricName := strings.Join([]string{PricesMetricName, m.nodeIdentifier, strings.ToLower(providerName), strings.ToLower(pairID)}, ".")
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

	metricName := strings.Join([]string{AggregatePricesMetricName, m.nodeIdentifier, strings.ToLower(pairID)}, ".")
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

	metricName := strings.Join([]string{ProviderTickMetricName, m.nodeIdentifier, strings.ToLower(providerName), strings.ToLower(pairID)}, ".")
	m.statsdClient.Incr(metricName, []string{fmt.Sprintf("%t", success)}, 1)
}

// AddProviderCountForMarket increments the number of providers that were utilized
// to calculate the final price for a given market.
func (m *OracleMetricsImpl) AddProviderCountForMarket(market string, count int) {
	m.promProviderCount.With(prometheus.Labels{
		PairIDLabel: strings.ToLower(market),
	},
	).Set(float64(count))

	metricName := strings.Join([]string{ProviderCountMetricName, m.nodeIdentifier, strings.ToLower(market)}, ".")
	m.statsdClient.Gauge(metricName, float64(count), []string{}, 1)
}

// MissingPrices updates the list of missing prices for the given tick.
func (m *OracleMetricsImpl) MissingPrices(pairIDs []string) {
	m.missingPricesMtx.Lock()
	defer m.missingPricesMtx.Unlock()

	m.missingPricesInternal = pairIDs
}

// GetMissingPrices gets the internal missing prices array.
func (m *OracleMetricsImpl) GetMissingPrices() []string {
	m.missingPricesMtx.Lock()
	defer m.missingPricesMtx.Unlock()

	return m.missingPricesInternal
}

// SetConnectBuildInfo sets the build information for the Connect binary. The version exported
// is determined by the build time version in accordance with the build pkg.
func (m *OracleMetricsImpl) SetConnectBuildInfo() {
	m.promConnectBuildInfo.With(prometheus.Labels{
		Version: build.Build,
	}).Set(1)

	encodedBuild := strings.ToLower(strings.ReplaceAll(build.Build, ".", "_"))
	metricName := strings.Join([]string{ConnectBuildInfoMetricName, m.nodeIdentifier, encodedBuild}, ".")
	m.statsdClient.Gauge(metricName, float64(1), []string{}, 1)
}
