package metrics

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/connect/v2/oracle/config"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// ObserveOracleResponseLatency records the time it took for the oracle to respond (this is a histogram)
	ObserveOracleResponseLatency(duration time.Duration)

	// AddOracleResponse increments the number of oracle responses, this can represent a liveness counter. This metric is paginated by status.
	AddOracleResponse(status Labeller)

	// ObserveABCIMethodLatency reports the given latency (as a duration), for the given ABCIMethod, and updates the ABCIMethodLatency histogram w/ that value.
	ObserveABCIMethodLatency(method ABCIMethod, duration time.Duration)

	// AddABCIRequest updates a counter corresponding to the given ABCI method and status.
	AddABCIRequest(method ABCIMethod, status Labeller)

	// ObserveMessageSize updates a histogram per Connect message type with the size of that message
	ObserveMessageSize(msg MessageType, size int)

	// ObservePriceForTicker updates a gauge with the price for the given ticker, this is updated each time a price is written to state
	ObservePriceForTicker(ticker connecttypes.CurrencyPair, price float64)

	// AddValidatorPriceForTicker updates a gauge per validator with the price they observed for a given ticker, this is updated when prices
	// to be written to state are aggregated
	AddValidatorPriceForTicker(validator string, ticker connecttypes.CurrencyPair, price float64)

	// AddValidatorReportForTicker updates a counter per validator + status. This counter represents the number of times a validator
	// for a ticker with a price, w/o a price, or w/ an absent.
	AddValidatorReportForTicker(validator string, ticker connecttypes.CurrencyPair, status ReportStatus)
}

type nopMetricsImpl struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &nopMetricsImpl{}
}

func (m *nopMetricsImpl) ObserveOracleResponseLatency(_ time.Duration)                 {}
func (m *nopMetricsImpl) AddOracleResponse(_ Labeller)                                 {}
func (m *nopMetricsImpl) ObserveABCIMethodLatency(_ ABCIMethod, _ time.Duration)       {}
func (m *nopMetricsImpl) AddABCIRequest(_ ABCIMethod, _ Labeller)                      {}
func (m *nopMetricsImpl) ObserveMessageSize(_ MessageType, _ int)                      {}
func (m *nopMetricsImpl) ObservePriceForTicker(_ connecttypes.CurrencyPair, _ float64) {}
func (m *nopMetricsImpl) AddValidatorReportForTicker(_ string, _ connecttypes.CurrencyPair, _ ReportStatus) {
}

func (m *nopMetricsImpl) AddValidatorPriceForTicker(_ string, _ connecttypes.CurrencyPair, _ float64) {
}

func NewMetrics(chainID string) Metrics {
	m := &metricsImpl{
		oracleResponseLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "oracle_response_latency",
			Help:      "The time it took for the oracle to respond",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}, []string{ChainIDLabel}),
		oracleResponseCounter: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: AppNamespace,
			Name:      "oracle_responses",
			Help:      "The number of oracle responses",
		}, []string{StatusLabel, ChainIDLabel}),
		abciMethodLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "abci_method_latency",
			Help:      "The time it took for an ABCI method to execute Connect specific logic (in seconds)",
			Buckets:   []float64{.0001, .0004, .002, .009, .02, .1, .65, 2, 6, 25},
		}, []string{ABCIMethodLabel, ChainIDLabel}),
		abciRequests: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: AppNamespace,
			Name:      "abci_requests",
			Help:      "The number of requests made to the ABCI server",
		}, []string{ABCIMethodLabel, StatusLabel, ChainIDLabel}),
		messageSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "message_size",
			Help:      "The size of the message in bytes",
			Buckets:   []float64{100, 500, 1000, 2000, 3000, 4000, 5000, 10000},
		}, []string{ChainIDLabel, MessageTypeLabel}),
		prices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: AppNamespace,
			Name:      "prices",
			Help:      "The price of the ticker that is written to state",
		}, []string{ChainIDLabel, TickerLabel}),
		reportsPerValidator: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: AppNamespace,
			Name:      "reports_per_validator",
			Help:      "The price reported for a specific validator and ticker",
		}, []string{ChainIDLabel, ValidatorLabel, TickerLabel}),
		reportStatusPerValidator: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: AppNamespace,
			Name:      "report_status_per_validator",
			Help:      "The status of the report for a specific validator and ticker",
		}, []string{ChainIDLabel, ValidatorLabel, TickerLabel, StatusLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.oracleResponseLatency)
	prometheus.MustRegister(m.oracleResponseCounter)
	prometheus.MustRegister(m.abciMethodLatency)
	prometheus.MustRegister(m.abciRequests)
	prometheus.MustRegister(m.messageSize)
	prometheus.MustRegister(m.prices)
	prometheus.MustRegister(m.reportsPerValidator)
	prometheus.MustRegister(m.reportStatusPerValidator)

	m.chainID = chainID

	return m
}

type metricsImpl struct {
	oracleResponseLatency    *prometheus.HistogramVec
	oracleResponseCounter    *prometheus.GaugeVec
	reportsPerValidator      *prometheus.GaugeVec
	reportStatusPerValidator *prometheus.GaugeVec
	abciMethodLatency        *prometheus.HistogramVec
	abciRequests             *prometheus.GaugeVec
	messageSize              *prometheus.HistogramVec
	prices                   *prometheus.GaugeVec
	chainID                  string
}

func (m *metricsImpl) ObserveABCIMethodLatency(method ABCIMethod, duration time.Duration) {
	m.abciMethodLatency.With(prometheus.Labels{
		ABCIMethodLabel: method.String(),
		ChainIDLabel:    m.chainID,
	}).Observe(duration.Seconds())
}

func (m *metricsImpl) ObserveOracleResponseLatency(duration time.Duration) {
	m.oracleResponseLatency.With(prometheus.Labels{
		ChainIDLabel: m.chainID,
	}).Observe(float64(duration.Milliseconds()))
}

func (m *metricsImpl) AddOracleResponse(status Labeller) {
	m.oracleResponseCounter.With(prometheus.Labels{
		StatusLabel:  status.Label(),
		ChainIDLabel: m.chainID,
	}).Inc()
}

func (m *metricsImpl) AddABCIRequest(method ABCIMethod, status Labeller) {
	m.abciRequests.With(prometheus.Labels{
		ABCIMethodLabel: method.String(),
		StatusLabel:     status.Label(),
		ChainIDLabel:    m.chainID,
	}).Inc()
}

func (m *metricsImpl) ObserveMessageSize(messageType MessageType, size int) {
	m.messageSize.With(prometheus.Labels{
		ChainIDLabel:     m.chainID,
		MessageTypeLabel: messageType.String(),
	}).Observe(float64(size))
}

func (m *metricsImpl) ObservePriceForTicker(ticker connecttypes.CurrencyPair, price float64) {
	m.prices.With(prometheus.Labels{
		ChainIDLabel: m.chainID,
		TickerLabel:  strings.ToLower(ticker.String()),
	}).Set(price)
}

func (m *metricsImpl) AddValidatorPriceForTicker(validator string, ticker connecttypes.CurrencyPair, price float64) {
	m.reportsPerValidator.With(prometheus.Labels{
		ChainIDLabel:   m.chainID,
		TickerLabel:    strings.ToLower(ticker.String()),
		ValidatorLabel: validator,
	}).Set(price)
}

func (m *metricsImpl) AddValidatorReportForTicker(validator string, ticker connecttypes.CurrencyPair, rs ReportStatus) {
	m.reportStatusPerValidator.With(prometheus.Labels{
		ChainIDLabel:   m.chainID,
		ValidatorLabel: validator,
		TickerLabel:    strings.ToLower(ticker.String()),
		StatusLabel:    rs.String(),
	}).Inc()
}

// NewMetricsFromConfig returns a new Metrics implementation based on the config. The Metrics
// returned is safe to be used in the client, and in the Oracle used by the PreBlocker.
// If the metrics are not enabled, a nop implementation is returned.
func NewMetricsFromConfig(cfg config.AppConfig, chainID string) (Metrics, error) {
	if !cfg.MetricsEnabled {
		return NewNopMetrics(), nil
	}

	// ensure that the metrics are enabled
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	// create the metrics
	metrics := NewMetrics(chainID)
	return metrics, nil
}
