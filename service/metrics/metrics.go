package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	TickerLabel    = "ticker"
	InclusionLabel = "included"
	AppNamespace   = "app"
	ProviderLabel  = "provider"
	StatusLabel    = "status"
)

type Status int

const (
	StatusFailure Status = iota
	StatusSuccess
)

func (s Status) String() string {
	switch s {
	case StatusFailure:
		return "failure"
	case StatusSuccess:
		return "success"
	default:
		return "unknown"
	}
}

func StatusFromError(err error) Status {
	if err == nil {
		return StatusSuccess
	}
	return StatusFailure
}

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

	// ObserveProcessProposalTime records the time it took for the oracle-specific parts of process proposal
	ObserveProcessProposalTime(duration time.Duration)

	// ObservePrepareProposalTime records the time it took for the oracle-specific parts of prepare proposal
	ObservePrepareProposalTime(duration time.Duration)
}

type nopMetricsImpl struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &nopMetricsImpl{}
}

func (m *nopMetricsImpl) ObserveOracleResponseLatency(_ time.Duration) {}
func (m *nopMetricsImpl) AddOracleResponse(_ Status)                   {}
func (m *nopMetricsImpl) AddVoteIncludedInLastCommit(_ bool)           {}
func (m *nopMetricsImpl) AddTickerInclusionStatus(_ string, _ bool)    {}
func (m *nopMetricsImpl) ObservePrepareProposalTime(_ time.Duration)   {}
func (m *nopMetricsImpl) ObserveProcessProposalTime(_ time.Duration)   {}

func NewMetrics() Metrics {
	m := &metricsImpl{
		oracleResponseLatency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "oracle_response_latency",
			Help:      "The time it took for the oracle to respond",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}),
		oracleResponseCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "oracle_responses",
			Help:      "The number of oracle responses",
		}, []string{StatusLabel}),
		voteIncludedInLastCommit: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "vote_included_in_last_commit",
			Help:      "The number of times this validator's vote was included in the last commit",
		}, []string{InclusionLabel}),
		tickerInclusionStatus: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: AppNamespace,
			Name:      "ticker_inclusion_status",
			Help:      "The number of times a ticker was included (or not included) in this validator's vote",
		}, []string{TickerLabel, InclusionLabel}),
		prepareProposalTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "oracle_prepare_proposal_time",
			Help:      "The time it took for the oracle-specific parts of prepare proposal",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}),
		processProposalTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "oracle_process_proposal_time",
			Help:      "The time it took for the oracle-specific parts of process proposal",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		}),
	}

	// register the above metrics
	prometheus.MustRegister(m.oracleResponseLatency)
	prometheus.MustRegister(m.oracleResponseCounter)
	prometheus.MustRegister(m.voteIncludedInLastCommit)
	prometheus.MustRegister(m.tickerInclusionStatus)
	prometheus.MustRegister(m.prepareProposalTime)
	prometheus.MustRegister(m.processProposalTime)

	return m
}

type metricsImpl struct {
	oracleResponseLatency    prometheus.Histogram
	oracleResponseCounter    *prometheus.CounterVec
	voteIncludedInLastCommit *prometheus.CounterVec
	tickerInclusionStatus    *prometheus.CounterVec
	prepareProposalTime      prometheus.Histogram
	processProposalTime      prometheus.Histogram
}

func (m *metricsImpl) ObserveProcessProposalTime(duration time.Duration) {
	m.processProposalTime.Observe(float64(duration.Milliseconds()))
}

func (m *metricsImpl) ObservePrepareProposalTime(duration time.Duration) {
	m.prepareProposalTime.Observe(float64(duration.Milliseconds()))
}

func (m *metricsImpl) ObserveOracleResponseLatency(duration time.Duration) {
	m.oracleResponseLatency.Observe(float64(duration.Milliseconds()))
}

func (m *metricsImpl) AddOracleResponse(status Status) {
	m.oracleResponseCounter.With(prometheus.Labels{
		StatusLabel: status.String(),
	}).Inc()
}

func (m *metricsImpl) AddVoteIncludedInLastCommit(included bool) {
	m.voteIncludedInLastCommit.With(prometheus.Labels{
		InclusionLabel: strconv.FormatBool(included),
	}).Inc()
}

func (m *metricsImpl) AddTickerInclusionStatus(ticker string, included bool) {
	m.tickerInclusionStatus.With(prometheus.Labels{
		TickerLabel:    ticker,
		InclusionLabel: strconv.FormatBool(included),
	}).Inc()
}

// NewMetricsFromConfig returns a new Metrics implementation based on the config. The Metrics
// returned is safe to be used in the client, and in the Oracle used by the PreBlocker.
// If the metrics are not enabled, a nop implementation is returned.
func NewMetricsFromConfig(cfg config.AppConfig) Metrics {
	if !cfg.MetricsEnabled {
		return NewNopMetrics()
	}

	return NewMetrics()
}
