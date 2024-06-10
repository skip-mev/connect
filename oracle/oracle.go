package oracle

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/skip-mev/slinky/oracle/orchestrator"

	"go.uber.org/zap"

	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
)

var _ Oracle = (*OracleImpl)(nil)

// Oracle defines the expected interface for an oracle. It is consumed by the oracle server.
//
//go:generate mockery --name Oracle --filename mock_oracle.go
type Oracle interface {
	IsRunning() bool
	GetLastSyncTime() time.Time
	GetPrices() types.Prices
	Start(ctx context.Context) error
	Stop()
}

// OracleImpl implements the core component responsible for fetching exchange rates
// for a given set of currency pairs and determining exchange rates.
type OracleImpl struct { //nolint
	// --------------------- General Config --------------------- //
	mtx    sync.RWMutex
	logger *zap.Logger

	// --------------------- Provider Config --------------------- //
	orc *orchestrator.ProviderOrchestrator

	// running is the current status of the main oracle process (running or not).
	running atomic.Bool

	// metrics is the set of metrics that the oracle will expose.
	metrics oraclemetrics.Metrics
}

// New returns a new instance of an Oracle. The oracle inputs providers that are
// responsible for fetching prices for a given set of currency pairs (base, quote). The oracle
// will fetch new prices concurrently every oracleTicker interval. In the case where
// the oracle fails to fetch prices from a given provider, it will continue to fetch prices
// from the remaining providers. The oracle currently assumes that each provider aggregates prices
// using TWAPs, TVWAPs, etc. When determining final prices, the oracle will utilize the aggregateFn
// to compute the final price for each currency pair. By default, the oracle will compute the median
// price across all providers.
func New(orc *orchestrator.ProviderOrchestrator, opts ...Option) (*OracleImpl, error) {
	o := &OracleImpl{
		logger:  zap.NewNop(),
		metrics: oraclemetrics.NewNopMetrics(),
		orc:     orc,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o, nil
}

// IsRunning returns true if the oracle is running.
func (o *OracleImpl) IsRunning() bool {
	return o.running.Load()
}

// Start starts the (blocking) oracle process. It will return when the context
// is cancelled or the oracle is stopped. The oracle will fetch prices from each
// provider concurrently every oracleTicker interval.
func (o *OracleImpl) Start(ctx context.Context) error {
	return o.orc.Start(ctx)
}

// Stop stops the oracle process and waits for it to gracefully exit.
func (o *OracleImpl) Stop() {
	o.logger.Info("stopping oracle")
	o.orc.Stop()
}

func (o *OracleImpl) GetLastSyncTime() time.Time {
	return o.orc.GetLastSyncTime()
}

func (o *OracleImpl) GetPrices() types.Prices {
	return o.orc.GetPrices()
}
