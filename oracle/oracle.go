package oracle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ Oracle = (*OracleImpl)(nil)

// Oracle defines the expected interface for an oracle. It is consumed by the oracle server.
//
//go:generate mockery --name Oracle --filename mock_oracle.go
type Oracle interface {
	IsRunning() bool
	GetLastSyncTime() time.Time
	GetPrices() types.Prices
	GetMarketMap() *marketmaptypes.MarketMap
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
	providers []*types.PriceProvider

	// --------------------- Oracle Config --------------------- //
	// lastPriceSync is the last time the oracle successfully updated its prices.
	lastPriceSync time.Time

	// running is the current status of the main oracle process (running or not).
	running atomic.Bool

	// priceAggregator maintains the state of prices for each provider and
	// computes the aggregate price for each currency pair.
	priceAggregator PriceAggregator

	// marketMapGetter gets the latest market map. It is a method implemented on ProviderOrchestrator.
	marketMapGetter func() marketmaptypes.MarketMap

	// metrics is the set of metrics that the oracle will expose.
	metrics oraclemetrics.Metrics

	// updateInterval is the interval at which the oracle will fetch prices from
	// each provider.
	updateInterval time.Duration

	// maxCacheAge is the longest amount of time a price will stay in our cache
	maxCacheAge time.Duration
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
		closer:         ssync.NewCloser(),
		logger:         zap.NewNop(),
		metrics:        oraclemetrics.NewNopMetrics(),
		updateInterval: 1 * time.Second,
		maxCacheAge:    time.Minute, // default max cache age is 1 minute
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

	o.running.Store(true)
	defer o.running.Store(false)

	ticker := time.NewTicker(o.updateInterval)
	defer ticker.Stop()

	// set the slinky build info on startup
	o.metrics.SetSlinkyBuildInfo()

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
}

// tick executes a single oracle tick. It fetches prices from each provider's
// cache and computes the aggregated price for each currency pair.
func (o *OracleImpl) tick() {
	o.logger.Debug("starting oracle tick")

	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("oracle tick panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	// Reset the provider prices before fetching new prices.
	o.priceAggregator.Reset()

	// Retrieve the latest prices from each provider.
	for _, priceProvider := range o.providers {
		o.fetchPrices(priceProvider)
	}

	o.logger.Debug("oracle fetched prices from providers")

	// Compute aggregated prices and update the oracle.
	o.priceAggregator.AggregatePrices()
	o.setLastSyncTime(time.Now().UTC())

	// update the last sync time
	o.metrics.AddTick()

	o.logger.Info("oracle updated prices", zap.Time("last_sync", o.GetLastSyncTime()), zap.Int("num_prices", len(o.GetPrices())))
}

// fetchPrices retrieves the latest prices from a given provider and updates the aggregator
// iff the price age is less than the update interval.
func (o *OracleImpl) fetchPrices(provider *types.PriceProvider) {
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("provider panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	if !provider.IsRunning() {
		o.logger.Debug(
			"provider is not running",
			zap.String("provider", provider.Name()),
		)

		return
	}

	o.logger.Debug(
		"retrieving prices",
		zap.String("provider", provider.Name()),
		zap.String("data handler type",
			string(provider.Type())),
	)

	// Fetch and set prices from the provider.
	prices := provider.GetData()
	if prices == nil {
		o.logger.Debug(
			"provider returned nil prices",
			zap.String("provider", provider.Name()),
			zap.String("data handler type", string(provider.Type())),
		)

		return
	}

	timeFilteredPrices := make(types.Prices)
	for pair, result := range prices {
		// If the price is older than the maxCacheAge, skip it.
		diff := time.Now().UTC().Sub(result.Timestamp)
		if diff > o.maxCacheAge {
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
		timeFilteredPrices[pair.GetOffChainTicker()] = result.Value
	}

	o.logger.Debug("provider returned prices",
		zap.String("provider", provider.Name()),
		zap.String("data handler type", string(provider.Type())),
		zap.Int("prices", len(prices)),
	)
	o.priceAggregator.SetProviderPrices(provider.Name(), timeFilteredPrices)
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
func (o *OracleImpl) GetPrices() types.Prices {
	prices := o.priceAggregator.GetPrices()
	return prices
}

// GetMarketMap returns the current market map configuration from the ProviderOrchestrator.
func (o *OracleImpl) GetMarketMap() *marketmaptypes.MarketMap {
	if o.marketMapGetter != nil {
		mm := o.marketMapGetter()
		return &mm
	}
	return &marketmaptypes.MarketMap{}
}
