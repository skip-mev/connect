package aggregator

import (
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func ComputeMedianWithContext(_ sdk.Context) AggregateFn[string, map[types.CurrencyPair]*big.Int] {
	return ComputeMedian()
}

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
func ComputeMedian() AggregateFn[string, map[types.CurrencyPair]*big.Int] {
	return func(providers AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]) map[types.CurrencyPair]*big.Int {
		// Aggregate prices across all providers for each asset.
		pricesByAsset := make(map[types.CurrencyPair][]*big.Int)
		for _, providerPrices := range providers {
			for cp, price := range providerPrices {
				// Only include prices that are not nil
				if price == nil {
					continue
				}

				// Initialize the asset array if it doesn't exist
				if _, ok := pricesByAsset[cp]; !ok {
					pricesByAsset[cp] = make([]*big.Int, 0)
				}

				pricesByAsset[cp] = append(pricesByAsset[cp], price)
			}
		}

		medianPrices := make(map[types.CurrencyPair]*big.Int)

		// Iterate through all assets and compute the median price
		for cp, prices := range pricesByAsset {
			if len(prices) == 0 {
				continue
			}

			sort.SliceStable(prices, func(i, j int) bool {
				switch prices[i].Cmp(prices[j]) {
				case -1:
					return true
				case 1:
					return false
				default:
					return true
				}
			})

			middle := len(prices) / 2

			// If the number of prices is even, compute the average of the two middle prices.
			numPrices := len(prices)
			if numPrices%2 == 0 {
				medianPrice := new(big.Int).Add(prices[middle-1], prices[middle])
				medianPrice = medianPrice.Div(medianPrice, new(big.Int).SetUint64(2))

				medianPrices[cp] = medianPrice
			} else {
				medianPrices[cp] = prices[middle]
			}
		}

		return medianPrices
	}
}
