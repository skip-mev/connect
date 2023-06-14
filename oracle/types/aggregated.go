package types

import (
	"strings"
	"sync"

	"golang.org/x/exp/maps"
)

type (
	// AggregatedProviderPrices defines a type alias for a map
	// of provider -> asset -> TickerPrice
	AggregatedProviderPrices map[string]map[string]TickerPrice
)

// PriceAggregator is a simple aggregator for provider prices.
// It is thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
type PriceAggregator struct {
	mtx sync.RWMutex

	providerPrices AggregatedProviderPrices
}

func NewPriceAggregator() *PriceAggregator {
	return &PriceAggregator{
		providerPrices: make(AggregatedProviderPrices),
	}
}

// SetTickerPrices updates the price aggregator with the latest ticker prices
// from the given provider.
func (p *PriceAggregator) SetPrices(provider Provider, prices map[string]TickerPrice) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	providerName := strings.ToLower(provider.Name())

	if len(prices) == 0 {
		p.providerPrices[providerName] = make(map[string]TickerPrice)
		return
	}

	p.providerPrices[providerName] = prices
}

// GetProviderPrices returns a copy of the aggregated provider prices.
func (p *PriceAggregator) GetProviderPrices() AggregatedProviderPrices {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(AggregatedProviderPrices)
	maps.Copy(cpy, p.providerPrices)

	return cpy
}
