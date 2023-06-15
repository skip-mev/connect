package oracle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/oracle/types"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	"golang.org/x/sync/errgroup"
)

// Oracle implements the core component responsible for fetching exchange rates
// for a given set of currency pairs and determining exchange rates.
type Oracle struct {
	// --------------------- General Config --------------------- //
	mtx    sync.RWMutex
	logger log.Logger
	closer *ssync.Closer

	// --------------------- Provider Config --------------------- //
	// providerTimeout is the maximum amount of time to wait for a provider to
	// respond to a price request. If a provider fails to respond within this
	// timeout, the oracle will ignore the provider and continue to fetch prices
	// from the remaining providers.
	providerTimeout time.Duration

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
}

// New returns a new instance of an Oracle. The oracle inputs providers that are
// responsible for fetching prices for a given set of currency pairs (base, quote). The oracle
// will fetch new prices concurrently every oracleTicker interval. In the case where
// the oracle fails to fetch prices from a given provider, it will continue to fetch prices
// from the remaining providers. The oracle currently assumes that each provider aggregates prices
// using TWAPs, TVWAPs, etc. When determining the aggregated price for a
// given currency pair, the oracle will compute the median price across all providers.
func New(
	logger log.Logger,
	providerTimeout, oracleTicker time.Duration,
	providers []types.Provider,
	aggregateFn types.AggregateFn,
) *Oracle {
	if logger == nil {
		panic("logger cannot be nil")
	}

	if providers == nil {
		panic("price providers cannot be nil")
	}

	if providerTimeout > oracleTicker {
		panic("provider timeout cannot be greater than oracle ticker")
	}

	return &Oracle{
		logger:          logger,
		closer:          ssync.NewCloser(),
		providerTimeout: providerTimeout,
		oracleTicker:    oracleTicker,
		providers:       providers,
		priceAggregator: types.NewPriceAggregator(aggregateFn),
	}
}

// NewDefaultOracle returns a new instance of an Oracle with a default aggregate
// function. It registers the given providers with the same set of currency pairs.
// The default aggregate function computes the median price across all providers.
func NewDefaultOracle(
	logger log.Logger,
	providerTimeout, oracleTicker time.Duration,
	providers []types.Provider,
	currencyPairs []types.CurrencyPair,
) *Oracle {
	for _, provider := range providers {
		provider.SetPairs(currencyPairs...)
	}

	return New(logger, providerTimeout, oracleTicker, providers, types.ComputeMedian())
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

	o.status.Store(true)
	defer o.status.Store(false)

	ticker := time.NewTicker(o.oracleTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.Stop()
			return ctx.Err()

		case <-o.closer.Done():
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

	// Create a goroutine group to fetch prices from each provider concurrently. All
	// of the goroutines will be cancelled after the oracle timeout.
	groupCtx, cancel := context.WithDeadline(ctx, time.Now().Add(o.oracleTicker))
	g, _ := errgroup.WithContext(groupCtx)
	g.SetLimit(len(o.providers))

	// In the case where the oracle panics, the oracle will log the error, cancel all of the
	// the goroutines, will not update the prices and will attempt to fetch prices again on the next tick.
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
		g.Go(o.fetchPricesFn(priceProvider))
	}

	// By default, errorgroup will wait for all goroutines to finish before returning. In
	// the case where any one of the goroutines fails, the entire set of goroutines will
	// be cancelled and the oracle will not update the prices.
	if err := g.Wait(); err != nil {
		o.logger.Error("wait group failed with error", "err", err)
		return
	}

	// Compute aggregated prices and update the oracle.
	o.priceAggregator.UpdatePrices()
	o.SetLastSyncTime(time.Now().UTC())

	o.logger.Info("oracle updated prices")
}

// fetchPrices returns a closure that fetches prices from the given provider. This is meant
// to be used in a goroutine. It accepts the provider and price aggregator as inputs. In the
// case where the provider fails to fetch prices, we will log the error and not update the
// price aggregator. We gracefully handle panics by recovering and logging the error. If the
// function panics, the wait group will cancel all other goroutines and skip the update for the
// oracle.
func (o *Oracle) fetchPricesFn(provider types.Provider) func() error {
	return func() (err error) {
		o.logger.Info("fetching prices from provider", provider.Name())

		doneCh := make(chan bool, 1)
		errCh := make(chan error, 1)

		go func() {
			// Recover from any panics while fetching prices.
			defer func() {
				if r := recover(); r != nil {
					errCh <- fmt.Errorf("panic when fetching prices %v", r)
				}
			}()

			// Fetch and set prices from the provider.
			prices, err := provider.GetPrices()
			if err != nil {
				errCh <- err
				return
			}

			o.priceAggregator.SetProviderPrices(provider, prices)
			o.logger.Info(provider.Name(), "number of assets fetched", len(prices))

			doneCh <- true
		}()

		select {
		case <-doneCh:
			break

		case err := <-errCh:
			o.logger.Error("failed to fetch prices from provider", provider.Name(), err)
			o.priceAggregator.SetProviderPrices(provider, nil)
			break

		case <-time.After(o.providerTimeout):
			o.logger.Error("provider timed out", provider.Name())
			o.priceAggregator.SetProviderPrices(provider, nil)
			break
		}

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
func (o *Oracle) GetPrices() map[types.CurrencyPair]*uint256.Int {
	return o.priceAggregator.GetPrices()
}
