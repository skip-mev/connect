package median

import (
	"math/big"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
)

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
//
// NOTE: This function should only be used for testing purposes.
func ComputeMedian() types.PriceAggregationFn {
	return func(providers types.AggregatedProviderPrices) types.AggregatorPrices {
		// Aggregate prices across all providers for each asset.
		pricesByAsset := make(map[string][]*big.Float)
		for _, providerPrices := range providers {
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
		medianPrices := make(types.AggregatorPrices, len(pricesByAsset))
		for cp, prices := range pricesByAsset {
			if len(prices) == 0 {
				continue
			}

			medianPrices[cp] = math.CalculateMedian(prices)
		}

		return medianPrices
	}
}
