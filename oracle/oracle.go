package oracle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	servicetypes "github.com/skip-mev/slinky/service/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"golang.org/x/sync/errgroup"
)

var _ servicetypes.Oracle = (*Oracle)(nil)

// Oracle implements the core component responsible for fetching exchange rates
// for a given set of currency pairs and determining exchange rates.
type Oracle struct {
	// --------------------- General Config --------------------- //
	mtx    sync.RWMutex
	logger log.Logger
	closer *ssync.Closer

	// --------------------- Provider Config --------------------- //
	// Providers is the set of providers that the oracle will fetch prices from.
	// Each provider is responsible for fetching prices for a given set of
	// currency pairs (base, quote). The oracle will fetch prices from each
	// provider concurrently.
	providers []types.Provider

	// --------------------- Oracle Config --------------------- //
	// oracleTicker is the interval at which the oracle will fetch prices from
	// providers.
	oracleTicker time.Duration

	// lastPriceSync is the last time the oracle successfully updated its prices.
	lastPriceSync time.Time

	// status is the current status of the oracle (running or not).
	status atomic.Bool

	// priceAggregator maintains the state of prices for each provider and
	// computes the aggregate price for each currency pair.
	priceAggregator *types.PriceAggregator

	// metrics is the set of metrics that the oracle will expose.
	metrics metrics.Metrics
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
	logger log.Logger,
	oracleTicker time.Duration,
	providers []types.Provider,
	aggregateFn types.AggregateFn,
	m metrics.Metrics,
) *Oracle {
	if logger == nil {
		panic("logger cannot be nil")
	}

	if providers == nil {
		panic("price providers cannot be nil")
	}

	if m == nil {
		m = metrics.NewNopMetrics()
	}

	return &Oracle{
		logger:          logger,
		closer:          ssync.NewCloser(),
		oracleTicker:    oracleTicker,
		providers:       providers,
		priceAggregator: types.NewPriceAggregator(aggregateFn),
		metrics:         m,
	}
}

// NewOracleFromConfig returns a new oracle instance from the given OracleConfig.
func NewOracleFromConfig(logger log.Logger, cfg *config.Config) (*Oracle, error) {
	// construct providers from the given currency pairs
	providers, err := cfg.GetProviders(logger)
	if err != nil {
		return nil, err
	}

	// configure metrics for the oracle
	m := metrics.NewNopMetrics()
	if cfg.OracleMetrics.Enabled {
		m = metrics.NewMetrics()
	}

	return New(logger, cfg.UpdateInterval, providers, types.ComputeMedian(), m), nil
}

// NewDefaultOracle returns a new instance of an Oracle with a default aggregate
// function. It registers the given providers with the same set of currency pairs.
// The default aggregate function computes the median price across all providers.
func NewDefaultOracle(
	logger log.Logger,
	oracleTicker time.Duration,
	providers []types.Provider,
	currencyPairs []oracletypes.CurrencyPair,
) *Oracle {
	// validate the currency-pairs
	for _, cp := range currencyPairs {
		if err := cp.ValidateBasic(); err != nil {
			panic(fmt.Sprintf("invalid currency pair %s", cp))
		}
	}

	for _, provider := range providers {
		provider.SetPairs(currencyPairs...)
	}

	return New(logger, oracleTicker, providers, types.ComputeMedian(), nil)
}

// IsRunning returns true if the oracle is running.
func (o *Oracle) IsRunning() bool {
	return o.status.Load()
}

// Start starts the (blocking) oracle process. It will return when the context
// is cancelled or the oracle is stopped. The oracle will fetch prices from each
// provider concurrently every oracleTicker interval.
func (o *Oracle) Start(ctx context.Context) error {
	o.logger.Info("starting oracle")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	o.status.Store(true)
	defer o.status.Store(false)

	ticker := time.NewTicker(o.oracleTicker)
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
			o.tick(ctx)
		}
	}
}

// Stop stops the oracle process and waits for it to gracefully exit.
func (o *Oracle) Stop() {
	o.logger.Info("stopping oracle")

	o.closer.Close()
	<-o.closer.Done()
}

// tick executes a single oracle tick. It fetches prices from each provider
// concurrently and computes the aggregated price for each currency pair. The
// oracle then sets the aggregated prices. In the case where any one of the provider
// fails to provide a set of prices, the oracle will continue to aggregate prices
// from the remaining providers.
func (o *Oracle) tick(ctx context.Context) {
	o.logger.Info("starting oracle tick")

	// For safety, we first check if the context is cancelled before
	// fetching prices from each provider. Since go routines are non-preemptive
	// and are chosen at random when multiple are ready, we want to ensure that
	// we do not fetch prices from any provider if the context is cancelled.
	if ctx.Err() != nil {
		o.logger.Info("oracle tick skipped", "err", ctx.Err())
		return
	}

	// Create a goroutine group to fetch prices from each provider concurrently. All
	// of the goroutines will be cancelled after the oracle timeout.
	groupCtx, cancel := context.WithDeadline(ctx, time.Now().Add(o.oracleTicker))
	defer cancel()

	g, _ := errgroup.WithContext(groupCtx)
	g.SetLimit(len(o.providers))

	// In the case where the oracle panics, the oracle will log the error, cancel
	// all of the the goroutines, will not update the prices and will attempt to
	// fetch prices again on the next tick.
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("oracle tick panicked", "err", r)

			cancel()
			if err := g.Wait(); err != nil {
				o.logger.Error("wait group failed with error", "err", err)
			}

			o.logger.Info("oracle tick finished after recovering from panic")
		}
	}()

	// Reset all of the provider prices before fetching new prices.
	o.priceAggregator.ResetProviderPrices()

	// Fetch prices from each provider concurrently. Each provider is responsible
	// for fetching prices for the given set of (base, quote) currency pairs.
	for _, priceProvider := range o.providers {
		g.Go(o.fetchPricesFn(groupCtx, priceProvider))
	}

	// By default, errorgroup will wait for all goroutines to finish before returning. In
	// the case where any one of the goroutines fails, the entire set of goroutines will
	// be cancelled and the oracle will not update the prices.
	//
	// NOTE: Although each fetch routine catches and handles errors/panics, we expect this
	// to only happen in the case where the oracle is shutting down.
	if err := g.Wait(); err != nil {
		o.logger.Error("wait group failed with error", "err", err)
		return
	}

	o.logger.Info("oracle fetched prices from providers")

	// Compute aggregated prices and update the oracle.
	o.priceAggregator.UpdatePrices()
	o.SetLastSyncTime(time.Now().UTC())

	// update the last sync time
	o.metrics.AddTick()

	o.logger.Info("oracle updated prices")
}

// fetchPrices returns a closure that fetches prices from the given provider. This is meant
// to be used in a goroutine. It accepts the ctx and provider as inputs. In the case where the
// main go group gets cancelled, this will trigger this go routine to short-circuit return.
// In the case where the provider fails to fetch prices, we will log the error and not update
// the price aggregator. We gracefully handle panics by recovering and logging the error.
func (o *Oracle) fetchPricesFn(ctx context.Context, provider types.Provider) func() error {
	return func() error {
		o.logger.Info("fetching prices", "provider", provider.Name())

		pricesCh := make(chan map[oracletypes.CurrencyPair]types.QuotePrice)
		errCh := make(chan error)

		// start timer
		start := time.Now()

		go func() {
			// Recover from any panics while fetching prices.
			defer func() {
				if r := recover(); r != nil {
					errCh <- fmt.Errorf("panic while fetching prices: %v", r)
				}

				close(pricesCh)
				close(errCh)
			}()

			// Fetch and set prices from the provider.
			//
			// Note: Each provider MUST return a set of prices within the oracle timeout i.e.
			// the context deadline. In the case where the provider fails to return a set of
			// prices within the deadline, the provider prices will be ignored. If the context
			// gets cancelled before the provider returns a set of prices, the provider must
			// short-circuit return.
			prices, err := provider.GetPrices(ctx)
			if err != nil {
				errCh <- err
				return
			}

			// Switch on the context error. In the case where the context is cancelled or
			// the deadline is exceeded, we will log the error and not update the price
			// aggregator.
			if err := ctx.Err(); err != nil {
				errCh <- err
				return
			}

			pricesCh <- prices
		}()

		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
			o.logger.Info("context ended before receiving response", "provider", provider.Name(), "err", err)
			o.priceAggregator.SetProviderPrices(provider.Name(), nil)
		case prices := <-pricesCh:
			o.logger.Info("fetching prices finished", "provider", provider.Name(), "num_prices", len(prices))
			o.priceAggregator.SetProviderPrices(provider.Name(), prices)
		case err = <-errCh:
			o.logger.Error("fetching prices failed", "provider", provider.Name(), "err", err)
			o.priceAggregator.SetProviderPrices(provider.Name(), nil)
		}

		// update the provider status
		o.metrics.AddProviderResponse(provider.Name(), metrics.StatusFromError(err))
		// response time
		o.metrics.ObserveProviderResponseLatency(provider.Name(), time.Since(start))

		return nil
	}
}

// SetLastSyncTime sets the last time the oracle successfully updated prices.
func (o *Oracle) SetLastSyncTime(t time.Time) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.lastPriceSync = t
}

// GetLastSyncTime returns the last time the oracle successfully updated prices.
func (o *Oracle) GetLastSyncTime() time.Time {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	return o.lastPriceSync
}

// GetPrices returns the aggregate prices from the oracle.
func (o *Oracle) GetPrices() map[oracletypes.CurrencyPair]*uint256.Int {
	return o.priceAggregator.GetPrices()
}
