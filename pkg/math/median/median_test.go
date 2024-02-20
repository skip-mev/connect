package median_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/pkg/math/median"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestComputeMedian(t *testing.T) {
	testCases := []struct {
		name           string
		providerPrices aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]
		expectedPrices map[mmtypes.Ticker]*big.Int
	}{
		{
			"empty provider prices",
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{},
			map[mmtypes.Ticker]*big.Int{},
		},
		{
			"single provider price",
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
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
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
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
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
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
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
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

			if len(prices) != len(tc.expectedPrices) {
				t.Fatalf("expected %d prices, got %d", len(tc.expectedPrices), len(prices))
			}

			for asset, expectedPrice := range tc.expectedPrices {
				price, ok := prices[asset]
				if !ok {
					t.Fatalf("expected price for asset %s", asset)
				}

				if price.Cmp(expectedPrice) != 0 {
					t.Fatalf("expected price %s, got %s", expectedPrice, price)
				}
			}
		})
	}
}
