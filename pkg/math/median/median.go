package median

import (
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func ComputeMedianWithContext(_ sdk.Context) types.PriceAggregationFn {
	return ComputeMedian()
}

// ComputeMedian inputs the aggregated prices from all providers and computes
// the median price for each asset.
func ComputeMedian() types.PriceAggregationFn {
	return func(providers types.AggregatedProviderPrices) types.TickerPrices {
		// Aggregate prices across all providers for each asset.
		pricesByAsset := make(map[mmtypes.Ticker][]*big.Int)
		for _, providerPrices := range providers {
			for ticker, price := range providerPrices {
				// Only include prices that are not nil
				if price == nil {
					continue
				}

				// Initialize the asset array if it doesn't exist
				if _, ok := pricesByAsset[ticker]; !ok {
					pricesByAsset[ticker] = make([]*big.Int, 0)
				}

				pricesByAsset[ticker] = append(pricesByAsset[ticker], price)
			}
		}

		medianPrices := make(types.TickerPrices, len(pricesByAsset))
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

// SortBigInts is a stable slices sort for an array of big.Ints.
func SortBigInts(values []*big.Int) {
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
}

// CalculateMedian calculates the median from a list of big.Ints. Returns an
// average if the number of values is even.
func CalculateMedian(values []*big.Int) *big.Int {
	SortBigInts(values)

	middleIndex := len(values) / 2

	// Calculate the median.
	numValues := len(values)
	var median *big.Int
	if numValues%2 == 0 { // even
		median = new(big.Int).Add(values[middleIndex-1], values[middleIndex])
		median = median.Div(median, new(big.Int).SetUint64(2))
	} else { // odd
		median = values[middleIndex]
	}

	return median
}
