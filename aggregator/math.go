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

			medianPrices[cp] = CalculateMedian(prices)
		}

		return medianPrices
	}
}

// CalculateMedian calculates the median price from a list of prices. Returns an
// average of the two middle prices if the number of prices is even.
func CalculateMedian(prices []*big.Int) *big.Int {
	// Sort the prices.
	sort.SliceStable(prices, func(i, j int) bool {
		return prices[i].Cmp(prices[j]) < 0
	})

	// Calculate the median price.
	middle := len(prices) / 2
	median := new(big.Int).Set(prices[middle])

	// If the number of prices is even, compute the average of the two middle prices.
	if len(prices)%2 == 0 {
		median = median.Add(median, prices[middle-1])
		median = median.Div(median, big.NewInt(2))
	}

	return median
}
