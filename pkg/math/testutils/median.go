package testutils

import (
	"math/big"
	"sync"

	"github.com/skip-mev/connect/v2/oracle"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var _ oracle.PriceAggregator = &MedianAggregator{}

type MedianAggregator struct {
	mtx sync.Mutex

	providerPrices map[string]types.Prices
	finalPrices    types.Prices
}

// NewMedianAggregator returns a new Median aggregator.
func NewMedianAggregator() *MedianAggregator {
	return &MedianAggregator{
		providerPrices: make(map[string]types.Prices),
		finalPrices:    make(types.Prices),
	}
}

// SetProviderPrices updates the data aggregator with the given provider and data.
func (m *MedianAggregator) SetProviderPrices(provider string, data types.Prices) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if data == nil {
		data = make(types.Prices)
	}

	m.providerPrices[provider] = data
}

func (m *MedianAggregator) UpdateMarketMap(_ mmtypes.MarketMap) {}

// AggregatePrices inputs the aggregated prices from all providers and computes
// the median price for each asset.
func (m *MedianAggregator) AggregatePrices() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// Aggregate prices across all providers for each asset.
	pricesByAsset := make(map[string][]*big.Float)
	for _, providerPrices := range m.providerPrices {
		for cp, price := range providerPrices {
			// Only include prices that are not nil
			if price == nil {
				continue
			}

			// Initialize the asset array if it doesn't exist
			if _, ok := pricesByAsset[cp]; !ok {
				pricesByAsset[cp] = make([]*big.Float, 0)
			}

			pricesByAsset[cp] = append(pricesByAsset[cp], price)
		}
	}

	// Iterate through all assets and compute the median price
	medianPrices := make(types.Prices)
	for cp, prices := range pricesByAsset {
		if len(prices) == 0 {
			continue
		}

		medianPrices[cp] = math.CalculateMedian(prices)
	}
	m.finalPrices = medianPrices
}

// GetPrices returns the aggregated data the aggregator has.
func (m *MedianAggregator) GetPrices() types.Prices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.finalPrices
}

// Reset resets the data aggregator for all providers.
func (m *MedianAggregator) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.providerPrices = make(map[string]types.Prices)
}
