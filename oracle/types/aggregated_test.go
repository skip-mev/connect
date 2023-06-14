package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

var (
	btcusd = types.CurrencyPair{
		Base:  "btc",
		Quote: "usd",
	}

	ethusd = types.CurrencyPair{
		Base:  "eth",
		Quote: "usd",
	}

	usdtusd = types.CurrencyPair{
		Base:  "usdt",
		Quote: "usd",
	}
)

func TestComputeMedian(t *testing.T) {
	testCases := []struct {
		name           string
		providerPrices types.AggregatedProviderPrices
		expectedPrices map[types.CurrencyPair]sdk.Dec
	}{
		{
			"empty provider prices",
			types.AggregatedProviderPrices{},
			map[types.CurrencyPair]sdk.Dec{},
		},
		{
			"single provider price",
			types.AggregatedProviderPrices{
				"provider1": {
					btcusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(100)),
					},
					ethusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(200)),
					},
				},
			},
			map[types.CurrencyPair]sdk.Dec{
				btcusd: sdk.NewDecFromInt(sdk.NewInt(100)),
				ethusd: sdk.NewDecFromInt(sdk.NewInt(200)),
			},
		},
		{
			"multiple provider prices",
			types.AggregatedProviderPrices{
				"provider1": {
					btcusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(100)),
					},
					ethusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(200)),
					},
				},
				"provider2": {
					btcusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(200)),
					},
					ethusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(300)),
					},
				},
			},
			map[types.CurrencyPair]sdk.Dec{
				btcusd: sdk.NewDecFromInt(sdk.NewInt(200)),
				ethusd: sdk.NewDecFromInt(sdk.NewInt(300)),
			},
		},
		{
			"multiple provider prices with different assets",
			types.AggregatedProviderPrices{
				"provider1": {
					btcusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(100)),
					},
					ethusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(200)),
					},
				},
				"provider2": {
					btcusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(200)),
					},
					ethusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(300)),
					},
					usdtusd: types.TickerPrice{
						Price: sdk.NewDecFromInt(sdk.NewInt(400)),
					},
				},
			},
			map[types.CurrencyPair]sdk.Dec{
				btcusd:  sdk.NewDecFromInt(sdk.NewInt(200)),
				ethusd:  sdk.NewDecFromInt(sdk.NewInt(300)),
				usdtusd: sdk.NewDecFromInt(sdk.NewInt(400)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			medianFn := types.ComputeMedian()
			prices := medianFn(tc.providerPrices)

			if len(prices) != len(tc.expectedPrices) {
				t.Fatalf("expected %d prices, got %d", len(tc.expectedPrices), len(prices))
			}

			for asset, expectedPrice := range tc.expectedPrices {
				price, ok := prices[asset]
				if !ok {
					t.Fatalf("expected price for asset %s", asset)
				}

				if !price.Equal(expectedPrice) {
					t.Fatalf("expected price %s, got %s", expectedPrice, price)
				}
			}
		})
	}
}
