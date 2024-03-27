package oracle_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/providers/websockets/okx"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

func TestAggregateData(t *testing.T) {
	testCases := []struct {
		name           string
		malleate       func(aggregator types.PriceAggregator)
		expectedPrices types.TickerPrices
	}{
		{
			name:           "no data",
			malleate:       func(types.PriceAggregator) {},
			expectedPrices: types.TickerPrices{},
		},
		{
			name: "coinbase direct feed, coinbase adjusted feed, no feed for usdt/usd - fail to report BTC/USD",
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(70_000, BTC_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USD: createPrice(69_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)

				prices = types.TickerPrices{
					BTC_USD: createPrice(71_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(okx.Name, prices)
			},
			expectedPrices: types.TickerPrices{},
		},
		{
			name: "coinbase direct feed, coinbase adjusted feed, binance adjusted feed for BTC/USD with index prices - success",
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(70_000, BTC_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USD: createPrice(69_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)

				prices = types.TickerPrices{
					BTC_USD: createPrice(71_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(okx.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: types.TickerPrices{
				BTC_USD: createPrice(78_500, BTC_USD.Decimals), // median of 70_000, 78,100, 79_500
			},
		},
		{
			name: "coinbase USDT direct, coinbase USDC/USDT inverted, binance direct feeds for USDT/USD - success",
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					USDT_USD:  createPrice(1.1, USDT_USD.Decimals),
					USDC_USDT: createPrice(1.1, USDC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					USDT_USD: createPrice(1.2, USDT_USD.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: types.TickerPrices{
				USDT_USD: createPrice(1.1, USDT_USD.Decimals), // median of 0.90909, 1, 1.2
			},
		},
		{
			name: "coinbase USDT direct, binance USDT/USD direct feeds for USDT/USD - success (average of two prices)",
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					USDT_USD: createPrice(1.2, USDT_USD.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: types.TickerPrices{
				USDT_USD: createPrice(1.15, USDT_USD.Decimals), // average of 1.1, 1.2
			},
		},
		{
			name: "coinbase USDT direct, kucoin BTC/USDT inverted, index BTC/USD direct feeds for USDT/USD - success",
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					USDT_USD: createPrice(1.0, USDT_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(kucoin.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD: createPrice(77_000, BTC_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: types.TickerPrices{
				USDT_USD: createPrice(1.05, USDT_USD.Decimals), // average of 1.1, 1.0
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap, metrics.NewNopMetrics())
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.DataAggregator)

			// Aggregate the data.
			m.AggregateData()

			// Ensure that the aggregated data is as expected.
			result := m.DataAggregator.GetAggregatedData()
			require.Equal(t, len(tc.expectedPrices), len(result))
			for ticker, price := range result {
				verifyPrice(t, tc.expectedPrices[ticker], price)
			}
		})
	}
}

func TestCalculateConvertedPrices(t *testing.T) {
	testCases := []struct {
		name            string
		target          mmtypes.Market
		providerConfigs []mmtypes.ProviderConfig
		malleate        func(aggregator types.PriceAggregator)
		expectedPrices  []*big.Int
	}{
		{
			name: "empty provider config",
			target: mmtypes.Market{
				Ticker:          BTC_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{},
			},
			malleate:       func(types.PriceAggregator) {},
			expectedPrices: make([]*big.Int, 0),
		},
		{
			name: "single conversion path with a single direct conversion (BTC/USD)",
			target: mmtypes.Market{
				Ticker: BTC_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  "btc-usd",
						NormalizeByPair: nil,
						Invert:          false,
						Metadata_JSON:   "",
					},
				},
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(70_000, BTC_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrices: []*big.Int{createPrice(70_000, BTC_USD.Decimals)},
		},
		{
			name: "single conversion with a single adjusted conversion (BTC/USDT * USDT/USD = BTC/USD)",
			target: mmtypes.Market{
				Ticker: BTC_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  "btc-usdt",
						NormalizeByPair: &USDT_USD.CurrencyPair,
						Invert:          false,
						Metadata_JSON:   "",
					},
				},
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: []*big.Int{createPrice(77_000, BTC_USD.Decimals)},
		},
		{
			name: "single inverted conversion path with a single adjusted conversion (USDT/BTC * BTC/USD = USDT/USD)",
			target: mmtypes.Market{
				Ticker: USDT_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  "btc-usdt",
						NormalizeByPair: &BTC_USD.CurrencyPair,
						Invert:          true,
						Metadata_JSON:   "",
					},
				},
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD: createPrice(77_000, BTC_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: []*big.Int{createPrice(1.1, USDT_USD.Decimals)},
		},
		{
			name: "single conversion path with a single adjusted conversion (USDC/USDT ^ -1 = USDT/USDC)",
			target: mmtypes.Market{
				Ticker: USDT_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  "btc-usdt",
						NormalizeByPair: &USDT_USD.CurrencyPair,
						Invert:          false,
						Metadata_JSON:   "",
					},
				},
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					USDC_USDT: createPrice(1.1, USDC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrices: []*big.Int{createPrice(0.9090909090909090909090909091, USDT_USD.Decimals)},
		},
		{
			name: "two provider configs both with a single direct conversion (BTC/USD)",
			target: mmtypes.Market{
				Ticker: BTC_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  "btc-usd",
						NormalizeByPair: nil,
						Invert:          false,
						Metadata_JSON:   "",
					},
					{
						Name:            binance.Name,
						OffChainTicker:  "btc-usd",
						NormalizeByPair: nil,
						Invert:          false,
						Metadata_JSON:   "",
					},
				},
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(70_000, BTC_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USD: createPrice(69_000, BTC_USD.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)
			},
			expectedPrices: []*big.Int{
				createPrice(70_000, BTC_USD.Decimals),
				createPrice(69_000, BTC_USD.Decimals),
			},
		},
		{
			name: "two provider configs both with a single adjusted conversion (BTC/USDT * USDT/USD = BTC/USD)",
			target: mmtypes.Market{
				Ticker: BTC_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  "btc-usdt",
						NormalizeByPair: &USDT_USD.CurrencyPair,
						Invert:          false,
						Metadata_JSON:   "",
					},
					{
						Name:            binance.Name,
						OffChainTicker:  "btc-usd",
						NormalizeByPair: nil,
						Invert:          false,
						Metadata_JSON:   "",
					},
				},
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				prices = types.TickerPrices{
					BTC_USDT: createPrice(69_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(binance.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrices: []*big.Int{
				createPrice(77_000, BTC_USD.Decimals),
				createPrice(75_900, BTC_USD.Decimals),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap, metrics.NewNopMetrics())
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.DataAggregator)

			// Calculate the converted prices.
			prices := m.CalculateConvertedPrices(tc.target)
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
		name           string
		target         mmtypes.Ticker
		providerConfig mmtypes.ProviderConfig
		malleate       func(aggregator types.PriceAggregator)
		expectedPrice  *big.Int
		expectedErr    bool
	}{
		{
			name:           "empty providerConfig",
			target:         BTC_USD,
			providerConfig: mmtypes.ProviderConfig{},
			malleate:       func(types.PriceAggregator) {},
			expectedPrice:  nil,
			expectedErr:    true,
		},
		{
			name:   "price does not exist for the provider with an operation that is exactly the target (BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name: coinbase.Name,
			},
			malleate:      func(types.PriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price exists for the provider with an operation that is exactly the target (BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name: coinbase.Name,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(70_000, BTC_USD.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(70_000, BTC_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price needs to be adjusted but the index price does not exist (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price needs to be adjusted and the index price exists (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(70_000, BTC_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price needs to be inverted to determine the adjusted price (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "btc_usdt",
				Invert:          true,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(70_000, BTC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD: createPrice(70_000, BTC_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1, USDT_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price is adjusted using USDT/USDC pairings (USDC/USDT ^ -1 = USDT/USDC)",
			target: USDT_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "usdc-usdt",
				Invert:         true,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					USDC_USDT: createPrice(1.1, USDC_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(0.9090909090909090909090909091, USDT_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price is adjust using eth pairings (ETH/USDT * USDT/USD = ETH/USD)",
			target: ETH_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "eth-usdt",
				Invert:          false,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					ETH_USDT: createPrice(4_000, ETH_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(4_400, ETH_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price for USDT/USD needs to be adjust by eth prices (USDT/ETH * ETH/USD = USDT/USD)",
			target: USDT_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "eth-usdt",
				Invert:          false,
				NormalizeByPair: &ETH_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					ETH_USDT: createPrice(4_100, ETH_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					ETH_USD: createPrice(4_000, ETH_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.97560975, USDT_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "price for PEPE/USDT needs to be adjusted by USDT/USD (different decimals) (PEPE/USDT * USDT/USD = PEPE/USD)",
			target: PEPE_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "pepe-usdt",
				Invert:          false,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					PEPE_USDT: createPrice(0.00000831846, PEPE_USDT.Decimals),
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.000009150306, PEPE_USDT.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a direct conversion with a sufficiently small number (BTC/USD = BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "btc-usd",
				Invert:         false,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(0.0000001, BTC_USD.Decimals), // 0.0000001 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(0.0000001, BTC_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with a sufficiently small number (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "btc-usdt",
				Invert:          false,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(0.0000001, BTC_USDT.Decimals), // 0.0000001 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(0.0000001, BTC_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with inverting with a sufficiently small number (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "btc-usdt",
				Invert:          true,
				NormalizeByPair: &BTC_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(0.00001, BTC_USDT.Decimals), // 0.00001 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD: createPrice(0.00002, BTC_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(2, USDT_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a direct conversion with a sufficiently large number (BTC/USD = BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "btc-usd",
				Invert:         false,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USD: createPrice(1_000_000_000, BTC_USD.Decimals), // 1,000,000,000 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)
			},
			expectedPrice: createPrice(1_000_000_000, BTC_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with a sufficiently large number (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "btc-usdt",
				Invert:          false,
				NormalizeByPair: &USDT_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(1_000_000_000, BTC_USDT.Decimals), // 1,000,000,000 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					USDT_USD: createPrice(1.1, USDT_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1_100_000_000, BTC_USD.Decimals),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with inverting with a sufficiently large number (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD,
			providerConfig: mmtypes.ProviderConfig{
				Name:            coinbase.Name,
				OffChainTicker:  "btc-usdt",
				Invert:          true,
				NormalizeByPair: &BTC_USD.CurrencyPair,
			},
			malleate: func(aggregator types.PriceAggregator) {
				prices := types.TickerPrices{
					BTC_USDT: createPrice(1_000_000_000, BTC_USDT.Decimals), // 1,000,000,000 BTC
				}
				aggregator.SetProviderData(coinbase.Name, prices)

				indexPrices := types.TickerPrices{
					BTC_USD: createPrice(1_100_000_000, BTC_USD.Decimals),
				}
				aggregator.SetAggregatedData(indexPrices)
			},
			expectedPrice: createPrice(1.1, USDT_USD.Decimals),
			expectedErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewMedianAggregator(logger, marketmap, metrics.NewNopMetrics())
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m.DataAggregator)

			// Calculate the adjusted price.
			price, err := m.CalculateAdjustedPrice(tc.target, tc.providerConfig)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			verifyPrice(t, tc.expectedPrice, price)
		})
	}
}
