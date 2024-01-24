package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
)

//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// ObserveOracleResponseLatency records the time it took for the oracle to respond (this is a histogram)
	ObserveOracleResponseLatency(duration time.Duration)

	// AddOracleResponse increments the number of oracle responses, this can represent a liveness counter. This metric is paginated by status.
	AddOracleResponse(status Status)

	// AddVoteIncludedInLastCommit increments the number of votes included in the last commit
	AddVoteIncludedInLastCommit(included bool)

	// AddTickerInclusionStatus increments the counter representing the number of times a ticker was included (or not included) in the last commit.
	AddTickerInclusionStatus(ticker string, included bool)

	// ObserveABCIMethodLatency reports the given latency (as a duration), for the given ABCIMethod, and updates the ABCIMethodLatency histogram w/ that value.
	ObserveABCIMethodLatency(method ABCIMethod, duration time.Duration)

	// AddABCIRequest updates a counter corresponding to the given ABCI method and status.
	AddABCIRequest(method ABCIMethod, status Labeller)
}

type nopMetricsImpl struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &nopMetricsImpl{}
}

func (m *nopMetricsImpl) ObserveOracleResponseLatency(_ time.Duration)           {}
func (m *nopMetricsImpl) AddOracleResponse(_ Status)                             {}
func (m *nopMetricsImpl) AddVoteIncludedInLastCommit(_ bool)                     {}
func (m *nopMetricsImpl) AddTickerInclusionStatus(_ string, _ bool)              {}
func (m *nopMetricsImpl) ObserveABCIMethodLatency(_ ABCIMethod, _ time.Duration) {}
func (m *nopMetricsImpl) AddABCIRequest(_ ABCIMethod, _ Labeller)                {}

func NewMetrics(chainID string) Metrics {
	m := &metricsImpl{
		oracleResponseLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "oracle_response_latency",
			Help:      "The time it took for the oracle to respond",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}, []string{ChainIDLabel}),
		oracleResponseCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "oracle_responses",
			Help:      "The number of oracle responses",
		}, []string{StatusLabel, ChainIDLabel}),
		voteIncludedInLastCommit: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "vote_included_in_last_commit",
			Help:      "The number of times this validator's vote was included in the last commit",
		}, []string{InclusionLabel, ChainIDLabel}),
		tickerInclusionStatus: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "ticker_inclusion_status",
			Help:      "The number of times a ticker was included (or not included) in this validator's vote",
		}, []string{TickerLabel, InclusionLabel, ChainIDLabel}),
		abciMethodLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "abci_method_latency",
			Help:      "The time it took for an ABCI method to execute slinky specific logic (in seconds)",
			Buckets:   []float64{.0001, .0004, .002, .009, .02, .1, .65, 2, 6, 25},
		}, []string{ABCIMethodLabel, ChainIDLabel}),
		abciRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "abci_requests",
			Help:      "The number of requests made to the ABCI server",
		}, []string{ABCIMethodLabel, StatusLabel, ChainIDLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.oracleResponseLatency)
	prometheus.MustRegister(m.oracleResponseCounter)
	prometheus.MustRegister(m.voteIncludedInLastCommit)
	prometheus.MustRegister(m.tickerInclusionStatus)
	prometheus.MustRegister(m.abciMethodLatency)

	m.chainID = chainID

	return m
}

type metricsImpl struct {
	oracleResponseLatency    *prometheus.HistogramVec
	oracleResponseCounter    *prometheus.CounterVec
	voteIncludedInLastCommit *prometheus.CounterVec
	tickerInclusionStatus    *prometheus.CounterVec
	abciMethodLatency        *prometheus.HistogramVec
	abciRequests			 *prometheus.CounterVec
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

func (m *metricsImpl) AddOracleResponse(status Status) {
	m.oracleResponseCounter.With(prometheus.Labels{
		StatusLabel:  status.String(),
		ChainIDLabel: m.chainID,
	}).Inc()
}

func (m *metricsImpl) AddVoteIncludedInLastCommit(included bool) {
	m.voteIncludedInLastCommit.With(prometheus.Labels{
		InclusionLabel: strconv.FormatBool(included),
		ChainIDLabel:   m.chainID,
	}).Inc()
}

func (m *metricsImpl) AddTickerInclusionStatus(ticker string, included bool) {
	m.tickerInclusionStatus.With(prometheus.Labels{
		TickerLabel:    ticker,
		InclusionLabel: strconv.FormatBool(included),
		ChainIDLabel:   m.chainID,
	}).Inc()
}

func (m *metricsImpl) AddABCIRequest(method ABCIMethod, status Labeller) {
	m.abciRequests.With(prometheus.Labels{
		ABCIMethodLabel: method.String(),
		StatusLabel:     status.Label(),
		ChainIDLabel:    m.chainID,
	}).Inc()
}

// NewMetricsFromConfig returns a new Metrics implementation based on the config. The Metrics
// returned is safe to be used in the client, and in the Oracle used by the PreBlocker.
// If the metrics are not enabled, a nop implementation is returned.
func NewMetricsFromConfig(cfg config.AppConfig, chainID string) (Metrics, error) {
	if !cfg.Enabled {
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
