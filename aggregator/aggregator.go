package aggregator

import (
	"math/big"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/maps"

	"github.com/skip-mev/slinky/x/oracle/types"
)

type (
	// AggregatedProviderPrices defines a type alias for a map
	// of provider -> asset -> QuotePrice
	AggregatedProviderPrices map[string]map[types.CurrencyPair]QuotePrice

	// AggregateFn is the function used to aggregate prices from each provider. Providers
	// should be responsible for aggregating prices using TWAPs, TVWAPs, etc. The oracle
	// will then compute the canonical price for a given currency pair by computing the
	// median price across all providers.
	AggregateFn func(providers AggregatedProviderPrices) map[types.CurrencyPair]*big.Int

	// AggregateFnFromContext is a function that is used to parametrize an aggregateFn by an sdk.Context. This is used
	// to allow the aggregateFn to access the latest state of an application. I.e computing a stake weighted median based
	// on the latest validator set.
	AggregateFnFromContext func(ctx sdk.Context) AggregateFn
)

// PriceAggregator is a simple aggregator for provider prices.
// It is thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
type PriceAggregator struct {
	mtx sync.RWMutex

	// aggregateFn is the function used to aggregate prices from each provider.
	aggregateFn AggregateFn

	// providerPrices is a map of provider -> asset -> QuotePrice
	providerPrices AggregatedProviderPrices

	// prices is the current set of prices aggregated across the providers.
	prices map[types.CurrencyPair]*big.Int
}

// NewPriceAggregator returns a PriceAggregator. The PriceAggregator
// is responsible for aggregating prices from each provider and computing
// the final oracle price for each asset. The PriceAggregator also tracks
// the current set of prices from each provider. The PricesAggregator is
// thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
func NewPriceAggregator(aggregateFn AggregateFn) *PriceAggregator {
	if aggregateFn == nil {
		panic("Aggregate function cannot be nil")
	}

	return &PriceAggregator{
		providerPrices: make(AggregatedProviderPrices),
		aggregateFn:    aggregateFn,
		prices:         make(map[types.CurrencyPair]*big.Int),
	}
}

// GetProviderPrices returns a copy of the aggregated provider prices.
func (p *PriceAggregator) GetProviderPrices() AggregatedProviderPrices {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(AggregatedProviderPrices)
	maps.Copy(cpy, p.providerPrices)

	return cpy
}

// GetPricesByProvider returns the prices for a given provider.
func (p *PriceAggregator) GetPricesByProvider(provider string) map[types.CurrencyPair]QuotePrice {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(map[types.CurrencyPair]QuotePrice)
	maps.Copy(cpy, p.providerPrices[provider])

	return cpy
}

// SetQuotePrices updates the price aggregator with the latest ticker prices
// from the given provider.
func (p *PriceAggregator) SetProviderPrices(provider string, prices map[types.CurrencyPair]QuotePrice) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if len(prices) == 0 {
		p.providerPrices[provider] = make(map[types.CurrencyPair]QuotePrice)
		return
	}

	p.providerPrices[provider] = prices
}

// ResetProviderPrices resets the price aggregator for all providers.
func (p *PriceAggregator) ResetProviderPrices() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.providerPrices = make(AggregatedProviderPrices)
}

// SetAggregationFn sets the aggregate function used to aggregate prices from each provider.
func (p *PriceAggregator) SetAggregationFn(fn AggregateFn) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.aggregateFn = fn
}

// UpdatePrices updates the current set of prices by using the aggregate function.
func (p *PriceAggregator) UpdatePrices() {
	providerPrices := p.GetProviderPrices()

	// Ensure nil prices are not set
	if prices := p.aggregateFn(providerPrices); prices != nil {
		p.SetPrices(prices)
		return
	}

	p.SetPrices(make(map[types.CurrencyPair]*big.Int))
}

// GetPrices returns the aggregated prices based on the provided currency pairs.
func (p *PriceAggregator) GetPrices() map[types.CurrencyPair]*big.Int {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(map[types.CurrencyPair]*big.Int)
	maps.Copy(cpy, p.prices)

	return cpy
}

// SetPrices sets the current set of prices.
func (p *PriceAggregator) SetPrices(prices map[types.CurrencyPair]*big.Int) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.prices = prices
}
