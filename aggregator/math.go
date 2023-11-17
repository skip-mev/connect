package aggregator

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func ComputeMedianWithContext(_ sdk.Context) AggregateFn {
	return ComputeMedian()
}

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
func ComputeMedian() AggregateFn {
	return func(providers AggregatedProviderPrices) map[types.CurrencyPair]*uint256.Int {
		pricesByAsset := make(map[types.CurrencyPair][]QuotePrice)
		for _, providerPrices := range providers {
			for cp, ticker := range providerPrices {
				// Only include prices that are not nil
				if ticker.Price == nil {
					continue
				}

				// Initialize the asset array if it doesn't exist
				if _, ok := pricesByAsset[cp]; !ok {
					pricesByAsset[cp] = make([]QuotePrice, 0)
				}

				pricesByAsset[cp] = append(pricesByAsset[cp], ticker)
			}
		}

		medianPrices := make(map[types.CurrencyPair]*uint256.Int)

		// Iterate through all assets and compute the median price
		for cp, prices := range pricesByAsset {
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

				medianPrices[cp] = medianPrice
			} else {
				medianPrices[cp] = prices[middle].Price
			}
		}

		return medianPrices
	}
}
