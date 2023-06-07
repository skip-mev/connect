package types

import (
	"sync"

	"golang.org/x/exp/maps"
)

type (
	// AggregatedProviderPrices defines a type alias for a map
	// of provider -> asset -> TickerPrice
	AggregatedProviderPrices map[string]map[string]TickerPrice

	// AggregatedProviderCandles defines a type alias for a map
	// of provider -> asset -> []Candle
	AggregatedProviderCandles map[string]map[string][]Candle
)

// PriceAggregator is a simple aggregator for provider prices and candles.
// It is thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
type PriceAggregator struct {
	mtx sync.RWMutex

	providerPrices  AggregatedProviderPrices
	providerCandles AggregatedProviderCandles
}

func NewPriceAggregator() *PriceAggregator {
	return &PriceAggregator{
		providerPrices:  make(AggregatedProviderPrices),
		providerCandles: make(AggregatedProviderCandles),
	}
}

func (p *PriceAggregator) SetTickerPricesAndCandles(
	providerName string,
	prices map[string]TickerPrice,
	candles map[string][]Candle,
	pair CurrencyPair,
) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// set prices and candles for this provider if we haven't seen it before
	if _, ok := p.providerPrices[providerName]; !ok {
		p.providerPrices[providerName] = make(map[string]TickerPrice)
	}
	if _, ok := p.providerCandles[providerName]; !ok {
		p.providerCandles[providerName] = make(map[string][]Candle)
	}

	// set price for provider/base (e.g. Binance -> ATOM -> 11.98)
	tp, pricesOk := prices[pair.String()]
	if pricesOk {
		p.providerPrices[providerName][pair.Base] = tp
	}

	// set candle for provider/base (e.g. Binance -> ATOM-> [<11.98, 24000, 12:00UTC>])
	cp, candlesOk := candles[pair.String()]
	if candlesOk {
		p.providerCandles[providerName][pair.Base] = cp
	}

	// return true if we set at least one price or candle
	return pricesOk || candlesOk
}

func (p *PriceAggregator) GetProviderPrices() AggregatedProviderPrices {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(AggregatedProviderPrices)
	maps.Copy(cpy, p.providerPrices)

	return cpy
}

func (p *PriceAggregator) GetProviderCandles() AggregatedProviderCandles {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	cpy := make(AggregatedProviderCandles)
	maps.Copy(cpy, p.providerCandles)

	return cpy
}
