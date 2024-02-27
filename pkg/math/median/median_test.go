package median_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/median"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestComputeMedian(t *testing.T) {
	testCases := []struct {
		name           string
		providerPrices types.AggregatedProviderPrices
		expectedPrices map[mmtypes.Ticker]*big.Int
	}{
		{
			"empty provider prices",
			types.AggregatedProviderPrices{},
			map[mmtypes.Ticker]*big.Int{},
		},
		{
			"single provider price",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD:  big.NewInt(100),
					constants.ETHEREUM_USD: big.NewInt(200),
				},
			},
			map[mmtypes.Ticker]*big.Int{
				constants.BITCOIN_USD:  big.NewInt(100),
				constants.ETHEREUM_USD: big.NewInt(200),
			},
		},
		{
			"multiple provider prices",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD:  big.NewInt(100),
					constants.ETHEREUM_USD: big.NewInt(200),
				},
				"provider2": {
					constants.BITCOIN_USD:  big.NewInt(200),
					constants.ETHEREUM_USD: big.NewInt(300),
				},
			},
			map[mmtypes.Ticker]*big.Int{
				constants.BITCOIN_USD:  big.NewInt(150),
				constants.ETHEREUM_USD: big.NewInt(250),
			},
		},
		{
			"multiple provider prices with different assets",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD:  big.NewInt(100),
					constants.ETHEREUM_USD: big.NewInt(200),
				},
				"provider2": {
					constants.BITCOIN_USD:  big.NewInt(200),
					constants.ETHEREUM_USD: big.NewInt(300),
					constants.USDT_USD:     nil, // should be ignored
				},
			},
			map[mmtypes.Ticker]*big.Int{
				constants.BITCOIN_USD:  big.NewInt(150),
				constants.ETHEREUM_USD: big.NewInt(250),
			},
		},
		{
			"odd number of provider prices",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD:  big.NewInt(100),
					constants.ETHEREUM_USD: big.NewInt(200),
				},
				"provider2": {
					constants.BITCOIN_USD:  big.NewInt(200),
					constants.ETHEREUM_USD: big.NewInt(300),
				},
				"provider3": {
					constants.BITCOIN_USD:  big.NewInt(300),
					constants.ETHEREUM_USD: big.NewInt(400),
				},
			},
			map[mmtypes.Ticker]*big.Int{
				constants.BITCOIN_USD:  big.NewInt(200),
				constants.ETHEREUM_USD: big.NewInt(300),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			medianFn := median.ComputeMedian()
			prices := medianFn(tc.providerPrices)
			require.Equal(t, len(tc.expectedPrices), len(prices))

			for asset, expectedPrice := range tc.expectedPrices {
				price, ok := prices[asset]
				require.True(t, ok)
				require.Equal(t, expectedPrice, price)
			}
		})
	}
}

func TestSortBigInts(t *testing.T) {
	testCases := []struct {
		name     string
		values   []*big.Int
		expected []*big.Int
	}{
		{
			name: "do nothing for nil slice",
		},
		{
			name: "sort a slice",
			values: []*big.Int{
				big.NewInt(10),
				big.NewInt(-2),
				big.NewInt(100),
				big.NewInt(0),
				big.NewInt(0),
			},
			expected: []*big.Int{
				big.NewInt(-2),
				big.NewInt(0),
				big.NewInt(0),
				big.NewInt(10),
				big.NewInt(100),
			},
		},
		{
			name: "do nothing for same values",
			values: []*big.Int{
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(-2),
				big.NewInt(100),
				big.NewInt(0),
				big.NewInt(0),
			},
			expected: []*big.Int{
				big.NewInt(-2),
				big.NewInt(0),
				big.NewInt(0),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(10),
				big.NewInt(100),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			median.SortBigInts(tc.values)
			require.Equal(t, tc.expected, tc.values)
		})
	}
}
