package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/metrics"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math/oracle"
	pkgtypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	usdtusdCP = pkgtypes.NewCurrencyPair("USDT", "USD")
	btcusdCP  = pkgtypes.NewCurrencyPair("BTC", "USD")
	ethusdCP  = pkgtypes.NewCurrencyPair("ETH", "USD")
)

func TestAggregateData(t *testing.T) {
	testCases := []struct {
		name           string
		malleate       func(aggregator *oracle.IndexPriceAggregator)
		expectedPrices types.Prices
	}{
		{
			name:           "no data",
			malleate:       func(*oracle.IndexPriceAggregator) {},
			expectedPrices: types.Prices{},
		},
		{
			name: "coinbase direct feed for BTC/USD - fail since it does not have enough providers",
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrices: types.Prices{},
		},
		{
			name: "coinbase direct feed, coinbase adjusted feed, binance adjusted feed for BTC/USD - fail since index price does not exist",
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD":  big.NewFloat(70_000),
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"BTCUSDT": big.NewFloat(69_000),
				}
				aggregator.SetProviderPrices(binance.Name, prices)
			},
			expectedPrices: types.Prices{},
		},
		{
			name: "coinbase direct feed, coinbase adjusted feed, binance adjusted feed for BTC/USD with index prices - success",
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD":  big.NewFloat(70_000),
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"BTCUSDT": big.NewFloat(69_000),
				}
				aggregator.SetProviderPrices(binance.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.1),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrices: types.Prices{
				BTC_USD.String(): big.NewFloat(75_900), // median of 70_000, 75_900, 77_000
			},
		},
		{
			name: "coinbase USDT direct, coinbase USDC/USDT inverted, binance direct feeds for USDT/USD - success",
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"USDT-USD":  big.NewFloat(1.1),
					"USDC-USDT": big.NewFloat(1.1),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"USDTUSD": big.NewFloat(1.2),
				}
				aggregator.SetProviderPrices(binance.Name, prices)
			},
			expectedPrices: types.Prices{
				USDT_USD.String(): big.NewFloat(1.1), // median of 0.90909, 1.1, 1.2
			},
		},
		{
			name: "coinbase USDT direct, binance USDT/USD direct feeds for USDT/USD - success (average of two prices)",
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"USDT-USD": big.NewFloat(1.1),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"USDTUSD": big.NewFloat(1.2),
				}
				aggregator.SetProviderPrices(binance.Name, prices)
			},
			expectedPrices: types.Prices{
				USDT_USD.String(): big.NewFloat(1.15), // average of 1.1, 1.2
			},
		},
		{
			name: "coinbase USDT direct, kucoin BTC/USDT inverted, index BTC/USD direct feeds for USDT/USD - success",
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"USDT-USD": big.NewFloat(1.0),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(kucoin.Name, prices)

				indexPrices := types.Prices{
					BTC_USD.String(): big.NewFloat(77_000),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrices: types.Prices{
				USDT_USD.String(): big.NewFloat(1.05), // average of 1.1, 1.0
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewIndexPriceAggregator(logger, marketmap, metrics.NewNopMetrics())
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m)

			// Aggregate the data.
			m.AggregatePrices()

			// Ensure that the aggregated data is as expected.
			result := m.GetIndexPrices()
			require.Equal(t, len(tc.expectedPrices), len(result))
			for ticker, price := range result {
				expectedPrice, ok := tc.expectedPrices[ticker]
				require.True(t, ok)
				require.Equal(t, expectedPrice.SetPrec(36), price.SetPrec(36))
			}
		})
	}
}

func TestCalculateConvertedPrices(t *testing.T) {
	testCases := []struct {
		name           string
		target         mmtypes.Ticker
		cfgs           []mmtypes.ProviderConfig
		malleate       func(aggregator *oracle.IndexPriceAggregator)
		expectedPrices []*big.Float
	}{
		{
			name:           "no conversion cfgs",
			target:         BTC_USD,
			cfgs:           []mmtypes.ProviderConfig{},
			malleate:       func(*oracle.IndexPriceAggregator) {},
			expectedPrices: make([]*big.Float, 0),
		},
		{
			name:   "single conversion path with a single direct conversion (BTC/USD)",
			target: BTC_USD,
			cfgs: []mmtypes.ProviderConfig{
				{
					Name:           coinbase.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrices: []*big.Float{big.NewFloat(70_000)},
		},
		{
			name:   "single conversion path with a single adjusted conversion (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			cfgs: []mmtypes.ProviderConfig{
				{
					Name:           coinbase.Name,
					OffChainTicker: "BTC-USDT",
					NormalizeByPair: &pkgtypes.CurrencyPair{
						Base:  "USDT",
						Quote: "USD",
					},
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.1),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrices: []*big.Float{big.NewFloat(77_000)},
		},
		{
			name:   "single conversion path with a single adjusted conversion (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD,
			cfgs: []mmtypes.ProviderConfig{
				{
					Name:           kucoin.Name,
					OffChainTicker: "BTC-USDT",
					Invert:         true,
					NormalizeByPair: &pkgtypes.CurrencyPair{
						Base:  "BTC",
						Quote: "USD",
					},
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(kucoin.Name, prices)

				indexPrices := types.Prices{
					btcusdCP.String(): big.NewFloat(77_000),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrices: []*big.Float{big.NewFloat(1.1)},
		},
		{
			name:   "single conversion path with a single adjusted conversion (USDC/USDT ^ -1 = USDT/USDC)",
			target: USDT_USD,
			cfgs: []mmtypes.ProviderConfig{
				{
					Name:           coinbase.Name,
					OffChainTicker: "USDC-USDT",
					Invert:         true,
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"USDC-USDT": big.NewFloat(1.1),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrices: []*big.Float{big.NewFloat(0.9090909090909090909090909091)},
		},
		{
			name:   "two conversion cfgs both with a single direct conversion (BTC/USD)",
			target: BTC_USD,
			cfgs: []mmtypes.ProviderConfig{
				{
					Name:           coinbase.Name,
					OffChainTicker: "BTC-USD",
				},
				{
					Name:           binance.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"BTC-USD": big.NewFloat(69_000),
				}
				aggregator.SetProviderPrices(binance.Name, prices)
			},
			expectedPrices: []*big.Float{
				big.NewFloat(70_000),
				big.NewFloat(69_000),
			},
		},
		{
			name:   "two conversion cfgs both with a single adjusted conversion (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			cfgs: []mmtypes.ProviderConfig{
				{
					Name:           coinbase.Name,
					OffChainTicker: "BTC-USDT",
					NormalizeByPair: &pkgtypes.CurrencyPair{
						Base:  "USDT",
						Quote: "USD",
					},
				},
				{
					Name:           binance.Name,
					OffChainTicker: "BTC-USDT",
					NormalizeByPair: &pkgtypes.CurrencyPair{
						Base:  "USDT",
						Quote: "USD",
					},
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				prices = types.Prices{
					"BTC-USDT": big.NewFloat(69_000),
				}
				aggregator.SetProviderPrices(binance.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.1),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrices: []*big.Float{
				big.NewFloat(77_000),
				big.NewFloat(75_900),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewIndexPriceAggregator(logger, marketmap, metrics.NewNopMetrics())
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m)

			// Calculate the converted prices.
			market := mmtypes.Market{
				Ticker:          tc.target,
				ProviderConfigs: tc.cfgs,
			}
			prices := m.CalculateConvertedPrices(market)
			require.Len(t, prices, len(tc.expectedPrices))
			if len(tc.expectedPrices) == 0 {
				require.Empty(t, prices)
				return
			}

			// Ensure that the prices are as expected.
			for i, price := range prices {
				require.Equal(t, tc.expectedPrices[i].SetPrec(36), price.SetPrec(36))
			}
		})
	}
}

func TestCalculateAdjustedPrice(t *testing.T) {
	testCases := []struct {
		name          string
		target        mmtypes.Ticker
		cfg           mmtypes.ProviderConfig
		malleate      func(aggregator *oracle.IndexPriceAggregator)
		expectedPrice *big.Float
		expectedErr   bool
	}{
		{
			name:          "empty cfg",
			target:        BTC_USD,
			cfg:           mmtypes.ProviderConfig{},
			malleate:      func(*oracle.IndexPriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price does not exist for the provider with an operation that is exactly the target (BTC/USD)",
			target: BTC_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "BTC-USD",
			},
			malleate:      func(*oracle.IndexPriceAggregator) {},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price exists for the provider with an operation that is exactly the target (BTC/USD)",
			target: BTC_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "BTC-USD",
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrice: big.NewFloat(70_000),
			expectedErr:   false,
		},
		{
			name:   "price needs to be adjusted but the index price does not exist (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "BTC-USDT",
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "USDT",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC_USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrice: nil,
			expectedErr:   true,
		},
		{
			name:   "price needs to be adjusted and the index price exists (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "BTC-USDT",
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "USDT",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.0),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrice: big.NewFloat(70_000),
			expectedErr:   false,
		},
		{
			name:   "price needs to be inverted to determine the adjusted price (USDT/BTC * BTC/USD = USDT/USD)",
			target: USDT_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           kucoin.Name,
				OffChainTicker: "BTC-USDT",
				Invert:         true,
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "BTC",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USDT": big.NewFloat(70_000),
				}
				aggregator.SetProviderPrices(kucoin.Name, prices)

				indexPrices := types.Prices{
					btcusdCP.String(): big.NewFloat(77_000),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrice: big.NewFloat(1.1),
			expectedErr:   false,
		},
		{
			name:   "price is adjusted using USDT/USDC pairings (USDC/USDT ^ -1 = USDT/USDC)",
			target: USDT_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "USDT-USDC",
				Invert:         true,
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"USDT-USDC": big.NewFloat(1.1),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrice: big.NewFloat(0.9090909090909090909090909091),
			expectedErr:   false,
		},
		{
			name:   "price is adjust using eth pairings (ETH/USDT * USDT/USD = ETH/USD)",
			target: ETH_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "ETH-USDT",
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "USDT",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"ETH-USDT": big.NewFloat(4_000),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.1),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrice: big.NewFloat(4_400),
			expectedErr:   false,
		},
		{
			name:   "price for USDT/USD needs to be adjust by eth prices (USDT/ETH * ETH/USD = USDT/USD)",
			target: USDT_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           kucoin.Name,
				OffChainTicker: "ETH-USDT",
				Invert:         true,
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "ETH",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"ETH-USDT": big.NewFloat(4_100),
				}
				aggregator.SetProviderPrices(kucoin.Name, prices)

				indexPrices := types.Prices{
					ethusdCP.String(): big.NewFloat(4_000),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrice: big.NewFloat(0.9756097561),
			expectedErr:   false,
		},
		{
			name:   "price for PEPE/USDT needs to be adjusted by USDT/USD (different decimals) (PEPE/USDT * USDT/USD = PEPE/USD)",
			target: PEPE_USD,
			cfg: mmtypes.ProviderConfig{
				OffChainTicker: "PEPEUSDT",
				Name:           binance.Name,
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "USDT",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"PEPEUSDT": big.NewFloat(0.00000831846),
				}
				aggregator.SetProviderPrices(binance.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.1),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrice: big.NewFloat(0.000009150306),
			expectedErr:   false,
		},
		{
			name:   "can make a direct conversion with a sufficiently small number (BTC/USD = BTC/USD)",
			target: BTC_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "BTC-USD",
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USD": big.NewFloat(0.1e-18),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)
			},
			expectedPrice: big.NewFloat(0.1e-18),
			expectedErr:   false,
		},
		{
			name:   "can make a adjusted conversion with a sufficiently small number (BTC/USDT * USDT/USD = BTC/USD)",
			target: BTC_USD,
			cfg: mmtypes.ProviderConfig{
				Name:           coinbase.Name,
				OffChainTicker: "BTC-USDT",
				NormalizeByPair: &pkgtypes.CurrencyPair{
					Base:  "USDT",
					Quote: "USD",
				},
			},
			malleate: func(aggregator *oracle.IndexPriceAggregator) {
				prices := types.Prices{
					"BTC-USDT": big.NewFloat(0.1e-18),
				}
				aggregator.SetProviderPrices(coinbase.Name, prices)

				indexPrices := types.Prices{
					usdtusdCP.String(): big.NewFloat(1.0),
				}
				aggregator.SetIndexPrices(indexPrices)
			},
			expectedPrice: big.NewFloat(0.1e-18),
			expectedErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := oracle.NewIndexPriceAggregator(logger, marketmap, metrics.NewNopMetrics())
			require.NoError(t, err)

			// Update the price aggregator with relevant data.
			tc.malleate(m)

			// Calculate the adjusted price.
			price, err := m.CalculateAdjustedPrice(tc.cfg)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedPrice.SetPrec(uint(36)), price.SetPrec(uint(36)))
		})
	}
}
