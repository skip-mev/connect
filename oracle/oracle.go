package oracle

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/metrics"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ Oracle = (*OracleImpl)(nil)

// Oracle defines the expected interface for an oracle. It is consumed by the oracle server.
//
//go:generate mockery --name Oracle --filename mock_oracle.go
type Oracle interface {
	IsRunning() bool
	GetLastSyncTime() time.Time
	GetPrices() map[slinkytypes.CurrencyPair]*big.Int
	Start(ctx context.Context) error
	Stop()
}

// OracleImpl implements the core component responsible for fetching exchange rates
// for a given set of currency pairs and determining exchange rates.
type OracleImpl struct { //nolint
	// --------------------- General Config --------------------- //
	mtx    sync.RWMutex
	logger *zap.Logger
	closer *ssync.Closer

	// --------------------- Provider Config --------------------- //
	// Providers is the set of providers that the oracle will fetch prices from.
	// Each provider is responsible for fetching prices for a given set of
	// currency pairs (base, quote). The oracle will fetch prices from each
	// provider concurrently.
	providers []providertypes.Provider[slinkytypes.CurrencyPair, *big.Int]

	// providerCh is the channel that the oracle will use to signal whether all of the
	// providers are running or not.
	providerCh chan error

	// --------------------- Oracle Config --------------------- //
	// lastPriceSync is the last time the oracle successfully updated its prices.
	lastPriceSync time.Time

	// running is the current status of the main oracle process (running or not).
	running atomic.Bool

	// priceAggregator maintains the state of prices for each provider and
	// computes the aggregate price for each currency pair.
	priceAggregator *aggregator.DataAggregator[string, map[slinkytypes.CurrencyPair]*big.Int]

	// metrics is the set of metrics that the oracle will expose.
	metrics metrics.Metrics

	// updateInterval is the interval at which the oracle will fetch prices from
	// each provider.
	updateInterval time.Duration
}

// New returns a new instance of an Oracle. The oracle inputs providers that are
// responsible for fetching prices for a given set of currency pairs (base, quote). The oracle
// will fetch new prices concurrently every oracleTicker interval. In the case where
// the oracle fails to fetch prices from a given provider, it will continue to fetch prices
// from the remaining providers. The oracle currently assumes that each provider aggregates prices
// using TWAPs, TVWAPs, etc. When determining final prices, the oracle will utilize the aggregateFn
// to compute the final price for each currency pair. By default, the oracle will compute the median
// price across all providers.
func New(opts ...Option) (*OracleImpl, error) {
	o := &OracleImpl{
		closer:  ssync.NewCloser(),
		logger:  zap.NewNop(),
		metrics: metrics.NewNopMetrics(),
		priceAggregator: aggregator.NewDataAggregator[string, map[slinkytypes.CurrencyPair]*big.Int](
			aggregator.WithAggregateFn(aggregator.ComputeMedian()),
		),
		updateInterval: 1 * time.Second,
	}

	for _, opt := range opts {
		opt(o)
	}

	o.logger.Info("creating oracle", zap.Int("num_providers", len(o.providers)))

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
	o.logger.Info("starting oracle")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	o.providerCh = make(chan error)
	go o.StartProviders(ctx)

	o.running.Store(true)
	defer o.running.Store(false)

	ticker := time.NewTicker(o.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.Stop()
			o.logger.Info("oracle stopped via context")
			return ctx.Err()

		case <-o.closer.Done():
			o.logger.Info("oracle stopped via closer")
			return nil

		case <-ticker.C:
			o.tick()
		}
	}
}

// Stop stops the oracle process and waits for it to gracefully exit.
func (o *OracleImpl) Stop() {
	o.logger.Info("stopping oracle")

	o.closer.Close()
	<-o.closer.Done()

	// Wait for the providers to exit.
	err := <-o.providerCh
	o.logger.Info("providers exited", zap.Error(err))
}

// tick executes a single oracle tick. It fetches prices from each provider's
// cache and computes the aggregated price for each currency pair.
func (o *OracleImpl) tick() {
	o.logger.Info("starting oracle tick")

	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("oracle tick panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	// Reset the provider prices before fetching new prices.
	o.priceAggregator.ResetProviderData()

	// Retrieve the latest prices from each provider.
	for _, priceProvider := range o.providers {
		o.fetchPrices(priceProvider)
	}

	o.logger.Info("oracle fetched prices from providers")

	// Compute aggregated prices and update the oracle.
	o.priceAggregator.AggregateData()
	o.setLastSyncTime(time.Now().UTC())

	// update the last sync time
	o.metrics.AddTick()

	o.logger.Info("oracle updated prices", zap.Time("last_sync", o.GetLastSyncTime()), zap.Int("num_prices", len(o.GetPrices())))
}

// fetchPrices retrieves the latest prices from a given provider and updates the aggregator
// iff the price age is less than the update interval.
func (o *OracleImpl) fetchPrices(provider providertypes.Provider[slinkytypes.CurrencyPair, *big.Int]) {
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("provider panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	o.logger.Info("retrieving prices", zap.String("provider", provider.Name()), zap.String("data handler type", string(provider.Type())))

	// Fetch and set prices from the provider.
	prices := provider.GetData()
	if prices == nil {
		o.logger.Info("provider returned nil prices", zap.String("provider", provider.Name()), zap.String("data handler type", string(provider.Type())))
		return
	}

	timeFilteredPrices := make(map[slinkytypes.CurrencyPair]*big.Int)
	for pair, result := range prices {
		floatValue, _ := result.Value.Float64() // we ignore the accuracy in this conversion

		// update price metric
		o.metrics.UpdatePrice(provider.Name(), string(provider.Type()), pair.String(), floatValue)

		// If the price is older than the update interval, skip it.
		diff := time.Now().UTC().Sub(result.Timestamp)
		if diff > o.updateInterval {
			o.logger.Debug(
				"skipping price",
				zap.String("provider", provider.Name()),
				zap.String("data handler type", string(provider.Type())),
				zap.String("pair", pair.String()),
				zap.Duration("diff", diff),
			)
			continue
		}

		o.logger.Debug(
			"adding price",
			zap.String("provider", provider.Name()),
			zap.String("data handler type", string(provider.Type())),
			zap.String("pair", pair.String()),
			zap.String("price", result.Value.String()),
			zap.Duration("diff", diff),
		)
		timeFilteredPrices[pair] = result.Value
	}

	o.logger.Info("provider returned prices",
		zap.String("provider", provider.Name()),
		zap.String("data handler type", string(provider.Type())),
		zap.Int("prices", len(prices)),
	)
	o.priceAggregator.SetProviderData(provider.Name(), timeFilteredPrices)
}

// GetLastSyncTime returns the last time the oracle successfully updated prices.
func (o *OracleImpl) GetLastSyncTime() time.Time {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	return o.lastPriceSync
}

// setLastSyncTime sets the last time the oracle successfully updated prices.
func (o *OracleImpl) setLastSyncTime(t time.Time) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.lastPriceSync = t
}

// GetPrices returns the aggregate prices from the oracle.
func (o *OracleImpl) GetPrices() map[slinkytypes.CurrencyPair]*big.Int {
	prices := o.priceAggregator.GetAggregatedData()

	// set metrics in background
	go func() {
		for cp, price := range prices {
			o.metrics.UpdateAggregatePrice(strings.ToLower(cp.String()), float64(price.Int64()))
		}
	}()

	return prices
}
