package median_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/median"
)

func TestComputeMedian(t *testing.T) {
	testCases := []struct {
		name           string
		providerPrices types.AggregatedProviderPrices
		expectedPrices map[string]*big.Float
	}{
		{
			"empty provider prices",
			types.AggregatedProviderPrices{},
			map[string]*big.Float{},
		},
		{
			"single provider price",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD.String():  big.NewFloat(100),
					constants.ETHEREUM_USD.String(): big.NewFloat(200),
				},
			},
			map[string]*big.Float{
				constants.BITCOIN_USD.String():  big.NewFloat(100),
				constants.ETHEREUM_USD.String(): big.NewFloat(200),
			},
		},
		{
			"multiple provider prices",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD.String():  big.NewFloat(100),
					constants.ETHEREUM_USD.String(): big.NewFloat(200),
				},
				"provider2": {
					constants.BITCOIN_USD.String():  big.NewFloat(200),
					constants.ETHEREUM_USD.String(): big.NewFloat(300),
				},
			},
			map[string]*big.Float{
				constants.BITCOIN_USD.String():  big.NewFloat(150),
				constants.ETHEREUM_USD.String(): big.NewFloat(250),
			},
		},
		{
			"multiple provider prices with different assets",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD.String():  big.NewFloat(100),
					constants.ETHEREUM_USD.String(): big.NewFloat(200),
				},
				"provider2": {
					constants.BITCOIN_USD.String():  big.NewFloat(200),
					constants.ETHEREUM_USD.String(): big.NewFloat(300),
					constants.USDT_USD.String():     nil, // should be ignored
				},
			},
			map[string]*big.Float{
				constants.BITCOIN_USD.String():  big.NewFloat(150),
				constants.ETHEREUM_USD.String(): big.NewFloat(250),
			},
		},
		{
			"odd number of provider prices",
			types.AggregatedProviderPrices{
				"provider1": {
					constants.BITCOIN_USD.String():  big.NewFloat(100),
					constants.ETHEREUM_USD.String(): big.NewFloat(200),
				},
				"provider2": {
					constants.BITCOIN_USD.String():  big.NewFloat(200),
					constants.ETHEREUM_USD.String(): big.NewFloat(300),
				},
				"provider3": {
					constants.BITCOIN_USD.String():  big.NewFloat(300),
					constants.ETHEREUM_USD.String(): big.NewFloat(400),
				},
			},
			map[string]*big.Float{
				constants.BITCOIN_USD.String():  big.NewFloat(200),
				constants.ETHEREUM_USD.String(): big.NewFloat(300),
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
