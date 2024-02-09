package aggregator

import (
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func ComputeMedianWithContext(_ sdk.Context) AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
	return ComputeMedian()
}

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
func ComputeMedian() AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int] {
	return func(providers AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]) map[slinkytypes.CurrencyPair]*big.Int {
		// Aggregate prices across all providers for each asset.
		pricesByAsset := make(map[slinkytypes.CurrencyPair][]*big.Int)
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

		medianPrices := make(map[slinkytypes.CurrencyPair]*big.Int)

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

// CalculateMedian calculates the median from a list of big.Ints. Returns an
// average if the number of values is even.
func CalculateMedian(values []*big.Int) *big.Int {
	// Sort the values.
	sort.SliceStable(values, func(i, j int) bool {
		switch values[i].Cmp(values[j]) {
		case -1:
			return true
		case 1:
			return false
		default:
			return true
		}
	})

	middle := len(values) / 2

	// Calculate the median.
	numValues := len(values)
	var median *big.Int
	if numValues%2 == 0 {
		median = new(big.Int).Add(values[middle-1], values[middle])
		median = median.Div(median, new(big.Int).SetUint64(2))
	} else {
		median = values[middle]
	}

	return median
}
