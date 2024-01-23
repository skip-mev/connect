package metrics

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	TickerLabel     = "ticker"
	InclusionLabel  = "included"
	AppNamespace    = "app"
	ProviderLabel   = "provider"
	StatusLabel     = "status"
	ABCIMethodLabel = "abci_method"
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

type ABCIMethod int

const (
	PrepareProposal ABCIMethod = iota
	ProcessProposal
	ExtendVote
	VerifyVoteExtension
	FinalizeBlock
)

func (a ABCIMethod) String() string {
	switch a {
	case PrepareProposal:
		return "prepare_proposal"
	case ProcessProposal:
		return "process_proposal"
	case ExtendVote:
		return "extend_vote"
	case VerifyVoteExtension:
		return "verify_vote_extension"
	case FinalizeBlock:
		return "finalize_block"
	default:
		return "not_implemented"
	}
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

	// ObserveABCIMethodLatency reports the given latency (as a duration), for the given ABCIMethod, and updates the ABCIMethodLatency histogram w/ that value.
	ObserveABCIMethodLatency(method ABCIMethod, duration time.Duration)
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
		abciMethodLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: AppNamespace,
			Name:      "abci_method_latency",
			Help:      "The time it took for an ABCI method to execute slinky specific logic (in seconds)",
			Buckets:   []float64{.0001, .0004, .002, .009, .02, .1, .65, 2, 6, 25},
		}, []string{ABCIMethodLabel}),
	}

	// register the above metrics
	prometheus.MustRegister(m.oracleResponseLatency)
	prometheus.MustRegister(m.oracleResponseCounter)
	prometheus.MustRegister(m.voteIncludedInLastCommit)
	prometheus.MustRegister(m.tickerInclusionStatus)
	prometheus.MustRegister(m.abciMethodLatency)

	return m
}

type metricsImpl struct {
	oracleResponseLatency    prometheus.Histogram
	oracleResponseCounter    *prometheus.CounterVec
	voteIncludedInLastCommit *prometheus.CounterVec
	tickerInclusionStatus    *prometheus.CounterVec
	abciMethodLatency        *prometheus.HistogramVec
}

func (m *metricsImpl) ObserveABCIMethodLatency(method ABCIMethod, duration time.Duration) {
	m.abciMethodLatency.With(prometheus.Labels{
		ABCIMethodLabel: method.String(),
	}).Observe(duration.Seconds())
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

// NewServiceMetricsFromConfig returns a new Metrics implementation based on the config. The Metrics
// returned is safe to be used in the client, and in the Oracle used by the PreBlocker.
// If the metrics are not enabled, a nop implementation is returned.
func NewServiceMetricsFromConfig(cfg config.AppMetricsConfig) (Metrics, sdk.ConsAddress, error) {
	if !cfg.Enabled {
		return NewNopMetrics(), nil, nil
	}

	// ensure that the metrics are enabled
	if err := cfg.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	// get the cons address
	consAddress, err := cfg.ConsAddress()
	if err != nil {
		return nil, nil, err
	}

	// create the metrics
	metrics := NewMetrics()
	return metrics, consAddress, nil
}
