package aggregator_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/aggregator"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	btcusd  = mmtypes.NewTicker("BITCOIN", "USD", 8, 1)
	ethusd  = mmtypes.NewTicker("ETHEREUM", "USD", 8, 1)
	usdtusd = mmtypes.NewTicker("USDT", "USD", 8, 1)
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
					btcusd: big.NewInt(100),
					ethusd: big.NewInt(200),
				},
			},
			map[mmtypes.Ticker]*big.Int{
				btcusd: big.NewInt(100),
				ethusd: big.NewInt(200),
			},
		},
		{
			"multiple provider prices",
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
				"provider1": {
					btcusd: big.NewInt(100),
					ethusd: big.NewInt(200),
				},
				"provider2": {
					btcusd: big.NewInt(200),
					ethusd: big.NewInt(300),
				},
			},
			map[mmtypes.Ticker]*big.Int{
				btcusd: big.NewInt(150),
				ethusd: big.NewInt(250),
			},
		},
		{
			"multiple provider prices with different assets",
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
				"provider1": {
					btcusd: big.NewInt(100),
					ethusd: big.NewInt(200),
				},
				"provider2": {
					btcusd:  big.NewInt(200),
					ethusd:  big.NewInt(300),
					usdtusd: nil, // should be ignored
				},
			},
			map[mmtypes.Ticker]*big.Int{
				btcusd: big.NewInt(150),
				ethusd: big.NewInt(250),
			},
		},
		{
			"odd number of provider prices",
			aggregator.AggregatedProviderData[string, map[mmtypes.Ticker]*big.Int]{
				"provider1": {
					btcusd: big.NewInt(100),
					ethusd: big.NewInt(200),
				},
				"provider2": {
					btcusd: big.NewInt(200),
					ethusd: big.NewInt(300),
				},
				"provider3": {
					btcusd: big.NewInt(300),
					ethusd: big.NewInt(400),
				},
			},
			map[mmtypes.Ticker]*big.Int{
				btcusd: big.NewInt(200),
				ethusd: big.NewInt(300),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			medianFn := aggregator.ComputeMedian()
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
