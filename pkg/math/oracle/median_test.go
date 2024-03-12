package oracle_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestAggregateData(t *testing.T) {}

func TestCalculateConvertedPrices(t *testing.T) {
	testCases := []struct {
		name           string
		target         mmtypes.Ticker
		paths          mmtypes.Paths
		malleate       func(aggregator *types.PriceAggregator)
		expectedPrices types.TickerPrices
		expectedErr    bool
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap)
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.PriceAggregator)

			// Calculate the converted prices.
			prices, err := m.CalculateConvertedPrices(tc.target, tc.paths)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedPrices, prices)
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
			name:   "price does not exist for the provider with an operation that is exactly the target (BTC/USD)",
			target: constants.BITCOIN_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
			},
			malleate:      func(aggregator *types.PriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price exists for the provider with an operation that is exactly the target (BTC/USD)",
			target: constants.BITCOIN_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.BITCOIN_USD: createPrice(70_000, constants.BITCOIN_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(70_000, constants.BITCOIN_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price needs to be adjusted but the index price does not exist (BTC/USDT * USDT/USD = BTC/USD)",
			target: constants.BITCOIN_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: constants.USDT_USD.CurrencyPair,
					Provider:     oracle.IndexProviderPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.BITCOIN_USDT: createPrice(70_000, constants.BITCOIN_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price needs to be adjusted and the index price exists (BTC/USDT * USDT/USD = BTC/USD)",
			target: constants.BITCOIN_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: constants.USDT_USD.CurrencyPair,
					Provider:     oracle.IndexProviderPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.BITCOIN_USDT: createPrice(70_000, constants.BITCOIN_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					constants.USDT_USD: createPrice(1, constants.USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(70_000, constants.BITCOIN_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price needs to be inverted to determine the adjusted price (USDT/BTC * BTC/USD = USDT/USD)",
			target: constants.USDT_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
				{
					CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
					Provider:     oracle.IndexProviderPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.BITCOIN_USDT: createPrice(70_000, constants.BITCOIN_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					constants.BITCOIN_USD: createPrice(70_000, constants.BITCOIN_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1, constants.USDT_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price is adjusted using USDT/USDC pairings (USDC/USDT ^ -1 = USDT/USDC)",
			target: constants.USDT_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.USDC_USDT.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.USDC_USDT: createPrice(1.1, constants.USDC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(0.9090909090909090909090909091, constants.USDT_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price is adjust using eth pairings (ETH/USDT * USDT/USD = ETH/USD)",
			target: constants.ETHEREUM_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       false,
				},
				{
					CurrencyPair: constants.USDT_USD.CurrencyPair,
					Provider:     oracle.IndexProviderPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.ETHEREUM_USDT: createPrice(4_000, constants.ETHEREUM_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					constants.USDT_USD: createPrice(1.1, constants.USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(4_400, constants.ETHEREUM_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price for USDT/USD needs to be adjust by eth prices (USDT/ETH * ETH/USD = USDT/USD)",
			target: constants.USDT_USD,
			operations: []mmtypes.Operation{
				{
					CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
					Provider:     coinbase.Name,
					Invert:       true,
				},
				{
					CurrencyPair: constants.ETHEREUM_USD.CurrencyPair,
					Provider:     oracle.IndexProviderPrice,
					Invert:       false,
				},
			},
			malleate: func(aggregator *types.PriceAggregator) {
				prices := types.TickerPrices{
					constants.ETHEREUM_USDT: createPrice(4_100, constants.ETHEREUM_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					constants.ETHEREUM_USD: createPrice(4_000, constants.ETHEREUM_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.97560975, constants.USDT_USD.Decimals),
			expectedErr:   false,
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
