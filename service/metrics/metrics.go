package metrics

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/skip-mev/slinky/oracle/metrics"
)

const (
	TickerLabel    = "ticker"
	InclusionLabel = "included"
	AppNamespace   = "app"
)

type Config struct {
	// Enabled indicates whether metrics should be enabled
	Enabled bool `mapstructure:"enabled" toml:"enabled"`

	// ValidatorConsAddress is the validator's consensus address
	ValidatorConsAddress string `mapstructure:"validator_cons_address" toml:"validator_cons_address"`
}

// ValidateBasic performs basic validation of the config
func (c Config) ValidateBasic() error {
	if c.Enabled {
		_, err := sdk.ConsAddressFromBech32(c.ValidatorConsAddress)
		return err
	}

	return nil
}

func (c Config) ConsAddress() (sdk.ConsAddress, error) {
	if c.Enabled {
		return sdk.ConsAddressFromBech32(c.ValidatorConsAddress)
	}

	return nil, nil
}

//go:generate mockery --name Metrics --filename mock_metrics.go
type Metrics interface {
	// ObserveOracleResponseLatency records the time it took for the oracle to respond (this is a histogram)
	ObserveOracleResponseLatency(duration time.Duration)

	// AddOracleResponse increments the number of oracle responses, this can represent a liveness counter. This metric is paginated by status.
	AddOracleResponse(status metrics.Status)

	// AddVoteIncludedInLastCommit increments the number of votes included in the last commit
	AddVoteIncludedInLastCommit(included bool)

	// AddTickerInclusionStatus increments the counter representing the number of times a ticker was included (or not included) in the last commit.
	AddTickerInclusionStatus(ticker string, included bool)
}

type nopMetricsImpl struct{}

// NewNopMetrics returns a Metrics implementation that does nothing.
func NewNopMetrics() Metrics {
	return &nopMetricsImpl{}
}

func (m *nopMetricsImpl) ObserveOracleResponseLatency(_ time.Duration) {}
func (m *nopMetricsImpl) AddOracleResponse(_ metrics.Status)           {}
func (m *nopMetricsImpl) AddVoteIncludedInLastCommit(_ bool)           {}
func (m *nopMetricsImpl) AddTickerInclusionStatus(_ string, _ bool)    {}

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
		}, []string{metrics.StatusLabel}),
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
	}

	// register the above metrics
	prometheus.MustRegister(m.oracleResponseLatency)
	prometheus.MustRegister(m.oracleResponseCounter)
	prometheus.MustRegister(m.voteIncludedInLastCommit)
	prometheus.MustRegister(m.tickerInclusionStatus)

	return m
}

type metricsImpl struct {
	oracleResponseLatency    prometheus.Histogram
	oracleResponseCounter    *prometheus.CounterVec
	voteIncludedInLastCommit *prometheus.CounterVec
	tickerInclusionStatus    *prometheus.CounterVec
}

func (m *metricsImpl) ObserveOracleResponseLatency(duration time.Duration) {
	m.oracleResponseLatency.Observe(float64(duration.Milliseconds()))
}

func (m *metricsImpl) AddOracleResponse(status metrics.Status) {
	m.oracleResponseCounter.With(prometheus.Labels{
		metrics.StatusLabel: status.String(),
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
func NewServiceMetricsFromConfig(cfg Config) (Metrics, sdk.ConsAddress, error) {
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
