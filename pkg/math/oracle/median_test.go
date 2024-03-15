package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestAggregateData(t *testing.T) {
	testCases := []struct {
		name           string
		malleate       func(aggregator *types.PriceAggregator)
		expectedPrices types.TickerPrices
	}{
		{
			name:           "no data",
			malleate:       func(*types.PriceAggregator) {},
			expectedPrices: types.TickerPrices{},
		},
		{
			name: "coinbase direct feed for BTC/USD - fail since it does not have enough providers",
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(70_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrices: types.TickerPrices{},
		},
		{
			name: "coinbase direct feed, coinbase adjusted feed, binance adjusted feed for BTC/USD - fail since index price does not exist",
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker:  createPrice(70_000, BTC_USD.Ticker.Decimals),
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USDT.Ticker: createPrice(69_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: types.TickerPrices{},
		},
		{
			name: "coinbase direct feed, coinbase adjusted feed, binance adjusted feed for BTC/USD with index prices - success",
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker:  createPrice(70_000, BTC_USD.Ticker.Decimals),
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USDT.Ticker: createPrice(69_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: types.TickerPrices{
				BTC_USD.Ticker: createPrice(75_900, BTC_USD.Ticker.Decimals), // median of 70_000, 75_900, 77_000
			},
		},
		{
			name: "coinbase USDT direct, coinbase USDC/USDT inverted, binance direct feeds for USDT/USD - success",
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					USDT_USD.Ticker:  createPrice(1.1, USDT_USD.Ticker.Decimals),
					USDC_USDT.Ticker: createPrice(1.1, USDC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.2, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: types.TickerPrices{
				USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals), // median of 0.90909, 1, 1.2
			},
		},
		{
			name: "coinbase USDT direct, binance USDT/USD direct feeds for USDT/USD - success (average of two prices)",
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.2, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: types.TickerPrices{
				USDT_USD.Ticker: createPrice(1.15, USDT_USD.Ticker.Decimals), // average of 1.1, 1.2
			},
		},
		{
			name: "coinbase USDT direct, kucoin BTC/USDT inverted, index BTC/USD direct feeds for USDT/USD - success",
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.0, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(kucoin.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(77_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: types.TickerPrices{
				USDT_USD.Ticker: createPrice(1.05, USDT_USD.Ticker.Decimals), // average of 1.1, 1.0
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap)
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.PriceAggregator)

			// Aggregate the data.
			m.AggregatedData()

			// Ensure that the aggregated data is as expected.
			result := m.PriceAggregator.GetAggregatedData()
			require.Equal(t, len(tc.expectedPrices), len(result))
			for ticker, price := range result {
				verifyPrice(t, tc.expectedPrices[ticker], price)
			}
		})
	}
}

func TestCalculateConvertedPrices(t *testing.T) {
	testCases := []struct {
		name           string
		target         mmtypes.Ticker
		paths          mmtypes.Paths
		malleate       func(aggregator *types.PriceAggregator)
		expectedPrices []*big.Int
	}{
		{
			name:   "too many conversion operations",
			target: BTC_USD.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       false,
							},
							{
								CurrencyPair: USDT_USD.Ticker.CurrencyPair,
								Provider:     mmtypes.IndexPrice,
								Invert:       false,
							},
							{
								CurrencyPair: USDT_USD.Ticker.CurrencyPair,
								Provider:     mmtypes.IndexPrice,
								Invert:       false,
							},
						},
					},
				},
			},
			malleate:       func(*types.PriceAggregator) {},
			expectedPrices: make([]*big.Int, 0),
		},
		{
			name:           "no conversion paths",
			target:         BTC_USD.Ticker,
			paths:          mmtypes.Paths{},
			malleate:       func(*types.PriceAggregator) {},
			expectedPrices: make([]*big.Int, 0),
		},
		{
			name:   "no conversion operations in a path",
			target: BTC_USD.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{},
					},
				},
			},
			malleate:       func(*types.PriceAggregator) {},
			expectedPrices: make([]*big.Int, 0),
		},
		{
			name:   "single conversion path with a single direct conversion (BTC/USD)",
			target: BTC_USD.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USD.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       false,
							},
						},
					},
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(70_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrices: []*big.Int{createPrice(70_000, BTC_USD.Ticker.Decimals)},
		},
		{
			name:   "single conversion path with a single adjusted conversion (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       false,
							},
							{
								CurrencyPair: USDT_USD.Ticker.CurrencyPair,
								Provider:     mmtypes.IndexPrice,
								Invert:       false,
							},
						},
					},
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: []*big.Int{createPrice(77_000, BTC_USD.Ticker.Decimals)},
		},
		{
			name:   "single conversion path with a single adjusted conversion (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       true,
							},
							{
								CurrencyPair: BTC_USD.Ticker.CurrencyPair,
								Provider:     mmtypes.IndexPrice,
								Invert:       false,
							},
						},
					},
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(77_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: []*big.Int{createPrice(1.1, USDT_USD.Ticker.Decimals)},
		},
		{
			name:   "single conversion path with a single adjusted conversion (USDC/USDT ^ -1 = USDT/USDC)",
			target: USDT_USD.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: USDC_USDT.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       true,
							},
						},
					},
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					USDC_USDT.Ticker: createPrice(1.1, USDC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrices: []*big.Int{createPrice(0.9090909090909090909090909091, USDT_USD.Ticker.Decimals)},
		},
		{
			name:   "two conversion paths both with a single direct conversion (BTC/USD)",
			target: BTC_USDT.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USD.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USD.Ticker.CurrencyPair,
								Provider:     binance.Name,
								Invert:       false,
							},
						},
					},
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(70_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USD.Ticker: createPrice(69_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: []*big.Int{
				createPrice(70_000, BTC_USD.Ticker.Decimals),
				createPrice(69_000, BTC_USD.Ticker.Decimals),
			},
		},
		{
			name:   "two conversion paths both with a single adjusted conversion (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			paths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
								Provider:     coinbase.Name,
								Invert:       false,
							},
							{
								CurrencyPair: USDT_USD.Ticker.CurrencyPair,
								Provider:     mmtypes.IndexPrice,
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
								Provider:     binance.Name,
								Invert:       false,
							},
							{
								CurrencyPair: USDT_USD.Ticker.CurrencyPair,
								Provider:     mmtypes.IndexPrice,
								Invert:       false,
							},
						},
					},
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USDT.Ticker: createPrice(69_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: []*big.Int{
				createPrice(77_000, BTC_USD.Ticker.Decimals),
				createPrice(75_900, BTC_USD.Ticker.Decimals),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap)
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.PriceAggregator)

			// Calculate the converted prices.
			prices := m.CalculateConvertedPrices(tc.target, tc.paths)
			require.Len(t, prices, len(tc.expectedPrices))
			if len(tc.expectedPrices) == 0 {
				require.Empty(t, prices)
				return
			}

			// Ensure that the prices are as expected.
			for i, price := range prices {
				verifyPrice(t, tc.expectedPrices[i], price)
			}
		})
	}
}

func TestCalculateAdjustedPrice(t *testing.T) {
	testCases := []struct {
		name          string
		target        mmtypes.Ticker
		operations    []mmtypes.Operation
		malleate      func(aggregator *types.PriceAggregator)
		expectedPrice *big.Int
		expectedErr   bool
	}{
		{
			name:          "nil operations",
			target:        BTC_USDT.Ticker,
			operations:    nil,
			malleate:      func(*types.PriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:          "empty operations",
			target:        BTC_USDT.Ticker,
			operations:    []mmtypes.Operation{},
			malleate:      func(*types.PriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "too many operations",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate:      func(*types.PriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price does not exist for the provider with an operation that is exactly the target (BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
			},
			malleate:      func(*types.PriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price exists for the provider with an operation that is exactly the target (BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(70_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(70_000, BTC_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price needs to be adjusted but the index price does not exist (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price needs to be adjusted and the index price exists (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(70_000, BTC_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price needs to be inverted to determine the adjusted price (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(70_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1, USDT_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price is adjusted using USDT/USDC pairings (USDC/USDT ^ -1 = USDT/USDC)",
			target: USDT_USD.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: USDC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					USDC_USDT.Ticker: createPrice(1.1, USDC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(0.9090909090909090909090909091, USDT_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price is adjust using eth pairings (ETH/USDT * USDT/USD = ETH/USD)",
			target: ETH_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: ETH_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					ETH_USDT.Ticker: createPrice(4_000, ETH_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(4_400, ETH_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price for USDT/USD needs to be adjust by eth prices (USDT/ETH * ETH/USD = USDT/USD)",
			target: USDT_USD.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: ETH_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
				{
					CurrencyPair: ETH_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					ETH_USDT.Ticker: createPrice(4_100, ETH_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					ETH_USD.Ticker: createPrice(4_000, ETH_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.97560975, USDT_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price for PEPE/USDT needs to be adjusted by USDT/USD (different decimals) (PEPE/USDT * USDT/USD = PEPE/USD)",
			target: PEPE_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: PEPE_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					PEPE_USDT.Ticker: createPrice(0.00000831846, PEPE_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.000009150306, PEPE_USDT.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a direct conversion with a sufficiently small number (BTC/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(0.0000001, BTC_USD.Ticker.Decimals), // 0.0000001 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(0.0000001, BTC_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with a sufficiently small number (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(0.0000001, BTC_USDT.Ticker.Decimals), // 0.0000001 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.0000001, BTC_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with inverting with a sufficiently small number (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(0.00001, BTC_USDT.Ticker.Decimals), // 0.00001 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(0.00002, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(2, USDT_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a direct conversion with a sufficiently large number (BTC/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(1_000_000_000, BTC_USD.Ticker.Decimals), // 1,000,000,000 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(1_000_000_000, BTC_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with a sufficiently large number (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(1_000_000_000, BTC_USDT.Ticker.Decimals), // 1,000,000,000 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1_100_000_000, BTC_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with inverting with a sufficiently large number (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
				{
					CurrencyPair: BTC_USD.Ticker.CurrencyPair,
					Provider:     mmtypes.IndexPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(1_000_000_000, BTC_USDT.Ticker.Decimals), // 1,000,000,000 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD.Ticker: createPrice(1_100_000_000, BTC_USD.Ticker.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1.1, USDT_USD.Ticker.Decimals),
			expectedErr:   false,
		},
		{
			name:   "second provider is not the index price",
			target: BTC_USDT.Ticker,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: BTC_USDT.Ticker.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: USDT_USD.Ticker.CurrencyPair,
					Provider:     binance.Name,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT.Ticker: createPrice(70_000, BTC_USDT.Ticker.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					USDT_USD.Ticker: createPrice(1.1, USDT_USD.Ticker.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrice: nil,
			expectedErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap)
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.PriceAggregator)

			// Calculate the adjusted price.
			price, err := m.CalculateAdjustedPrice(tc.target, tc.operations)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			verifyPrice(t, tc.expectedPrice, price)
		})
	}
}
