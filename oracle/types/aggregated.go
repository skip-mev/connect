package types

import (
	"sort"
	"sync"

	"github.com/holiman/uint256"
	"golang.org/x/exp/maps"
)

type (
	// AggregatedProviderPrices defines a type alias for a map
	// of provider -> asset -> QuotePrice
	AggregatedProviderPrices map[string]map[CurrencyPair]QuotePrice

	// AggregateFn is the function used to aggregate prices from each provider. Providers
	// should be responsible for aggregating prices using TWAPs, TVWAPs, etc. The oracle
	// will then compute the canonical price for a given currency pair by computing the
	// median price across all providers.
	AggregateFn func(providers AggregatedProviderPrices) map[CurrencyPair]*uint256.Int
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
	prices map[CurrencyPair]*uint256.Int
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
		prices:         make(map[CurrencyPair]*uint256.Int),
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

// SetQuotePrices updates the price aggregator with the latest ticker prices
// from the given provider.
func (p *PriceAggregator) SetProviderPrices(provider Provider, prices map[CurrencyPair]QuotePrice) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	providerName := provider.Name()

	if len(prices) == 0 {
		p.providerPrices[providerName] = make(map[CurrencyPair]QuotePrice)
		return
	}

	p.providerPrices[providerName] = prices
}

// ResetProviderPrices resets the price aggregator for all providers.
func (p *PriceAggregator) ResetProviderPrices() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.providerPrices = make(AggregatedProviderPrices)
}

// UpdatePrices updates the current set of prices by using the aggregate function.
func (p *PriceAggregator) UpdatePrices() {
	providerPrices := p.GetProviderPrices()

	prices := p.aggregateFn(providerPrices)
	p.SetPrices(prices)
}

// GetPrices returns the aggregated prices based on the provided currency pairs.
func (p *PriceAggregator) GetPrices() map[CurrencyPair]*uint256.Int {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(map[CurrencyPair]*uint256.Int)
	maps.Copy(cpy, p.prices)

	return cpy
}

// SetPrices sets the current set of prices.
func (p *PriceAggregator) SetPrices(prices map[CurrencyPair]*uint256.Int) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.prices = prices
}

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
func ComputeMedian() AggregateFn {
	return func(providers AggregatedProviderPrices) map[CurrencyPair]*uint256.Int {
		pricesByAsset := make(map[CurrencyPair][]QuotePrice)
		for _, providerPrices := range providers {
			for currencyPair, ticker := range providerPrices {
				// Only include prices that are not nil
				if ticker.Price == nil {
					continue
				}

				// Initialize the asset array if it doesn't exist
				if _, ok := pricesByAsset[currencyPair]; !ok {
					pricesByAsset[currencyPair] = make([]QuotePrice, 0)
				}

				pricesByAsset[currencyPair] = append(pricesByAsset[currencyPair], ticker)
			}
		}

		medianPrices := make(map[CurrencyPair]*uint256.Int)

		// Iterate through all assets and compute the median price
		for currencyPair, prices := range pricesByAsset {
			if len(prices) == 0 {
				continue
			}

			sort.SliceStable(prices, func(i, j int) bool {
				return prices[i].Price.Lt(prices[j].Price)
			})

			middle := len(prices) / 2

			// If the number of prices is even, compute the average of the two middle prices.
			numPrices := len(prices)
			if numPrices%2 == 0 {
				medianPrice := new(uint256.Int).Add(prices[middle-1].Price, prices[middle].Price)
				medianPrice = medianPrice.Div(medianPrice, new(uint256.Int).SetUint64(2))

				medianPrices[currencyPair] = medianPrice
			} else {
				medianPrices[currencyPair] = prices[middle].Price
			}
		}

		return medianPrices
	}
}
