package types

import (
	"sort"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/maps"
)

type (
	// AggregatedProviderPrices defines a type alias for a map
	// of provider -> asset -> TickerPrice
	AggregatedProviderPrices map[string]map[CurrencyPair]TickerPrice

	// AggregateFn is the function used to aggregate prices from each provider. Providers
	// should be responsible for aggregating prices using TWAPs, TVWAPs, etc. The oracle
	// will then compute the canonical price for a given currency pair by computing the
	// median price across all providers.
	AggregateFn func(providers AggregatedProviderPrices) map[CurrencyPair]sdk.Dec
)

// PriceAggregator is a simple aggregator for provider prices.
// It is thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
type PriceAggregator struct {
	mtx sync.RWMutex

	aggregateFn AggregateFn

	// providerPrices is a map of provider -> asset -> TickerPrice
	providerPrices AggregatedProviderPrices

	// prices is the current set of prices aggregated across the providers.
	prices map[CurrencyPair]sdk.Dec
}

func NewPriceAggregator(aggregateFn AggregateFn) *PriceAggregator {
	if aggregateFn == nil {
		panic("Aggregate function cannot be nil")
	}

	return &PriceAggregator{
		providerPrices: make(AggregatedProviderPrices),
		aggregateFn:    aggregateFn,
		prices:         make(map[CurrencyPair]sdk.Dec),
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

// SetTickerPrices updates the price aggregator with the latest ticker prices
// from the given provider.
func (p *PriceAggregator) SetProviderPrices(provider Provider, prices map[CurrencyPair]TickerPrice) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	providerName := provider.Name()

	if len(prices) == 0 {
		p.providerPrices[providerName] = make(map[CurrencyPair]TickerPrice)
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
func (p *PriceAggregator) GetPrices() map[CurrencyPair]sdk.Dec {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(map[CurrencyPair]sdk.Dec)
	maps.Copy(cpy, p.prices)

	return cpy
}

// SetPrices sets the current set of prices.
func (p *PriceAggregator) SetPrices(prices map[CurrencyPair]sdk.Dec) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.prices = prices
}

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
func ComputeMedian() AggregateFn {
	return func(providers AggregatedProviderPrices) map[CurrencyPair]sdk.Dec {
		pricesByAsset := make(map[CurrencyPair][]TickerPrice)
		for _, providerPrices := range providers {
			for currencyPair, price := range providerPrices {
				// Initialize the asset array if it doesn't exist
				if _, ok := pricesByAsset[currencyPair]; !ok {
					pricesByAsset[currencyPair] = make([]TickerPrice, 0)
				}

				pricesByAsset[currencyPair] = append(pricesByAsset[currencyPair], price)
			}
		}

		medianPrices := make(map[CurrencyPair]sdk.Dec)

		// Iterate through all assets and compute the median price
		for currencyPair, prices := range pricesByAsset {
			sort.SliceStable(prices, func(i, j int) bool {
				return prices[i].Price.LT(prices[j].Price)
			})

			medianPrice := prices[len(prices)/2].Price
			medianPrices[currencyPair] = medianPrice
		}

		return medianPrices
	}
}
