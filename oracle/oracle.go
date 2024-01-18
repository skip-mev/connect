package oracle

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	providertypes "github.com/skip-mev/slinky/providers/types"
	servicetypes "github.com/skip-mev/slinky/service/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ servicetypes.Oracle = (*Oracle)(nil)

// Oracle implements the core component responsible for fetching exchange rates
// for a given set of currency pairs and determining exchange rates.
type Oracle struct {
	// --------------------- General Config --------------------- //
	mtx    sync.RWMutex
	logger *zap.Logger
	closer *ssync.Closer

	// --------------------- Provider Config --------------------- //
	// Providers is the set of providers that the oracle will fetch prices from.
	// Each provider is responsible for fetching prices for a given set of
	// currency pairs (base, quote). The oracle will fetch prices from each
	// provider concurrently.
	providers []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]

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
	priceAggregator *aggregator.DataAggregator[string, map[oracletypes.CurrencyPair]*big.Int]

	// metrics is the set of metrics that the oracle will expose.
	metrics metrics.Metrics

	// cfg is the oracle config.
	cfg config.OracleConfig
}

// New returns a new instance of an Oracle. The oracle inputs providers that are
// responsible for fetching prices for a given set of currency pairs (base, quote). The oracle
// will fetch new prices concurrently every oracleTicker interval. In the case where
// the oracle fails to fetch prices from a given provider, it will continue to fetch prices
// from the remaining providers. The oracle currently assumes that each provider aggregates prices
// using TWAPs, TVWAPs, etc. When determining final prices, the oracle will utilize the aggregateFn
// to compute the final price for each currency pair. By default, the oracle will compute the median
// price across all providers.
func New(
	cfg config.OracleConfig,
	opts ...OracleOption,
) (*Oracle, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid oracle config: %w", err)
	}

	o := &Oracle{
		closer:  ssync.NewCloser(),
		cfg:     cfg,
		logger:  zap.NewNop(),
		metrics: metrics.NewNopMetrics(),
		priceAggregator: aggregator.NewDataAggregator[string, map[oracletypes.CurrencyPair]*big.Int](
			aggregator.WithAggregateFn(aggregator.ComputeMedian()),
		),
	}

	for _, opt := range opts {
		opt(o)
	}

	o.logger.Info("creating oracle", zap.Int("num_providers", len(o.providers)))

	return o, nil
}

// IsRunning returns true if the oracle is running.
func (o *Oracle) IsRunning() bool {
	return o.running.Load()
}

// Start starts the (blocking) oracle process. It will return when the context
// is cancelled or the oracle is stopped. The oracle will fetch prices from each
// provider concurrently every oracleTicker interval.
func (o *Oracle) Start(ctx context.Context) error {
	o.logger.Info("starting oracle")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	o.providerCh = make(chan error)
	go o.StartProviders(ctx)

	o.running.Store(true)
	defer o.running.Store(false)

	ticker := time.NewTicker(o.cfg.UpdateInterval)
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
func (o *Oracle) Stop() {
	o.logger.Info("stopping oracle")

	o.closer.Close()
	<-o.closer.Done()

	// Wait for the providers to exit.
	err := <-o.providerCh
	o.logger.Info("providers exited", zap.Error(err))
}

// tick executes a single oracle tick. It fetches prices from each provider's
// cache and computes the aggregated price for each currency pair.
func (o *Oracle) tick() {
	o.logger.Info("starting oracle tick")

	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("oracle tick panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	// Reset all of the provider prices before fetching new prices.
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
func (o *Oracle) fetchPrices(provider providertypes.Provider[oracletypes.CurrencyPair, *big.Int]) {
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("provider panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	o.logger.Info("retrieving prices", zap.String("provider", provider.Name()))

	// Fetch and set prices from the provider.
	prices := provider.GetData()
	if prices == nil {
		o.logger.Info("provider returned nil prices", zap.String("provider", provider.Name()))
		return
	}

	timeFilteredPrices := make(map[oracletypes.CurrencyPair]*big.Int)
	for pair, result := range prices {
		// If the price is older than the update interval, skip it.
		diff := time.Now().UTC().Sub(result.Timestamp)
		if diff > o.cfg.UpdateInterval {
			o.logger.Debug(
				"skipping price",
				zap.String("provider", provider.Name()),
				zap.String("pair", pair.String()),
				zap.Duration("diff", diff),
			)
			continue
		}

		o.logger.Debug(
			"adding price",
			zap.String("provider", provider.Name()),
			zap.String("pair", pair.String()),
			zap.String("price", result.Value.String()),
			zap.Duration("diff", diff),
		)
		timeFilteredPrices[pair] = result.Value
	}

	o.logger.Info("provider returned prices", zap.String("provider", provider.Name()), zap.Int("prices", len(prices)))
	o.priceAggregator.SetProviderData(provider.Name(), timeFilteredPrices)
}

// GetLastSyncTime returns the last time the oracle successfully updated prices.
func (o *Oracle) GetLastSyncTime() time.Time {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	return o.lastPriceSync
}

// setLastSyncTime sets the last time the oracle successfully updated prices.
func (o *Oracle) setLastSyncTime(t time.Time) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.lastPriceSync = t
}

// GetPrices returns the aggregate prices from the oracle.
func (o *Oracle) GetPrices() map[oracletypes.CurrencyPair]*big.Int {
	return o.priceAggregator.GetAggregatedData()
}
