package aggregator_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/aggregator"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

var (
	btcusd = slinkytypes.NewCurrencyPair("btc", "usd")

	ethusd = slinkytypes.NewCurrencyPair("eth", "usd")

	usdtusd = slinkytypes.NewCurrencyPair("usdt", "usd")
)

func TestComputeMedian(t *testing.T) {
	testCases := []struct {
		name           string
		providerPrices aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]
		expectedPrices map[slinkytypes.CurrencyPair]*big.Int
	}{
		{
			"empty provider prices",
			aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{},
			map[slinkytypes.CurrencyPair]*big.Int{},
		},
		{
			"single provider price",
			aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				"provider1": {
					btcusd: big.NewInt(100),
					ethusd: big.NewInt(200),
				},
			},
			map[slinkytypes.CurrencyPair]*big.Int{
				btcusd: big.NewInt(100),
				ethusd: big.NewInt(200),
			},
		},
		{
			"multiple provider prices",
			aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				"provider1": {
					btcusd: big.NewInt(100),
					ethusd: big.NewInt(200),
				},
				"provider2": {
					btcusd: big.NewInt(200),
					ethusd: big.NewInt(300),
				},
			},
			map[slinkytypes.CurrencyPair]*big.Int{
				btcusd: big.NewInt(150),
				ethusd: big.NewInt(250),
			},
		},
		{
			"multiple provider prices with different assets",
			aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
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
			map[slinkytypes.CurrencyPair]*big.Int{
				btcusd: big.NewInt(150),
				ethusd: big.NewInt(250),
			},
		},
		{
			"odd number of provider prices",
			aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
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
			map[slinkytypes.CurrencyPair]*big.Int{
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
