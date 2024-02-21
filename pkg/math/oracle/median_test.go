package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// acceptableDelta is the acceptable difference between the expected and actual price.
	// In this case, we use a delta of 1e-10. This means we will accept any price that is
	// within 1e-10 of the expected price.
	acceptableDelta = 1e-10

	usdtbtc = mmtypes.NewTicker("USDT", "BTC", 8, 1)

	logger, _ = zap.NewDevelopment()
	marketmap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			constants.BITCOIN_USD.String():      constants.BITCOIN_USD,
			constants.BITCOIN_USDC.String():     constants.BITCOIN_USDC,
			constants.BITCOIN_USDT.String():     constants.BITCOIN_USDT,
			constants.USDC_USD.String():         constants.USDC_USD,
			constants.USDT_USD.String():         constants.USDT_USD,
			constants.ETHEREUM_USD.String():     constants.ETHEREUM_USD,
			constants.ETHEREUM_USDT.String():    constants.ETHEREUM_USDT,
			constants.ETHEREUM_BITCOIN.String(): constants.ETHEREUM_BITCOIN,
			usdtbtc.String():                    usdtbtc,
		},
		Providers: map[string]mmtypes.Providers{
			constants.BITCOIN_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USD],
				},
			},
			constants.BITCOIN_USDC.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USDC],
				},
			},
			constants.BITCOIN_USDT.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USDT],
				},
			},
			constants.USDC_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.USDC_USD],
				},
			},
			constants.USDT_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.USDT_USD],
				},
			},
			constants.ETHEREUM_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_USD],
				},
			},
			constants.ETHEREUM_USDT.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_USDT],
				},
			},
			constants.ETHEREUM_BITCOIN.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_BITCOIN],
				},
			},
			usdtbtc.String(): {
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           "coinbase",
						OffChainTicker: "USDT/BTC",
					},
				},
			},
		},
		Paths: map[string]mmtypes.Paths{
			constants.BITCOIN_USD.String(): constants.BITCOIN_USD_PATHS,
		},
	}
)

func TestMedian(t *testing.T) {
	testCases := []struct {
		name              string
		pricesPerProvider types.AggregatedProviderPrices
		expected          types.TickerPrices
	}{
		{
			name: "no prices",
			pricesPerProvider: types.AggregatedProviderPrices{
				"coinbase": {},
			},
			expected: types.TickerPrices{},
		},
		{
			name: "single resolved price",
			pricesPerProvider: types.AggregatedProviderPrices{
				"coinbase": {
					constants.BITCOIN_USD: createPrice(40_000, 8),
				},
			},
			expected: types.TickerPrices{
				constants.BITCOIN_USD: createPrice(40_000, 8),
			},
		},
		{
			name: "must convert to get a single final price",
			pricesPerProvider: types.AggregatedProviderPrices{
				"coinbase": {
					constants.BITCOIN_USDT: createPrice(40_000, 8),
					constants.USDT_USD:     createPrice(1.1, 8),
				},
			},
			expected: types.TickerPrices{
				constants.BITCOIN_USD:  createPrice(44_000, 8),
				constants.BITCOIN_USDT: createPrice(40_000, 8),
				constants.USDT_USD:     createPrice(1.1, 8),
			},
		},
		{
			name: "calculates median price between two separate conversions",
			pricesPerProvider: types.AggregatedProviderPrices{
				"coinbase": {
					constants.BITCOIN_USD:  createPrice(40_000, 8),
					constants.BITCOIN_USDT: createPrice(40_000, 8),
					constants.USDT_USD:     createPrice(1.1, 8),
				},
			},
			expected: types.TickerPrices{
				constants.BITCOIN_USD:  createPrice(42_000, 8), // median average of 40_000 and 44_000
				constants.BITCOIN_USDT: createPrice(40_000, 8),
				constants.USDT_USD:     createPrice(1.1, 8),
			},
		},
		{
			name: "calculates median price between three separate conversions",
			pricesPerProvider: types.AggregatedProviderPrices{
				"coinbase": {
					constants.BITCOIN_USD:  createPrice(40_000, 8),
					constants.BITCOIN_USDT: createPrice(40_000, 8),
					constants.USDT_USD:     createPrice(1.1, 8),
					constants.BITCOIN_USDC: createPrice(40_000, 8),
					constants.USDC_USD:     createPrice(1.2, 8),
				},
			},
			expected: types.TickerPrices{
				constants.BITCOIN_USD:  createPrice(44_000, 8), // median of 40_000, 44_000, and 48_000
				constants.BITCOIN_USDT: createPrice(40_000, 8),
				constants.USDT_USD:     createPrice(1.1, 8),
				constants.BITCOIN_USDC: createPrice(40_000, 8),
				constants.USDC_USD:     createPrice(1.2, 8),
			},
		},
		{
			name: "calculates median price with a price of 0",
			pricesPerProvider: types.AggregatedProviderPrices{
				"coinbase": {
					constants.BITCOIN_USD: createPrice(0, 8),
				},
			},
			expected: types.TickerPrices{
				constants.BITCOIN_USD: createPrice(0, 8),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			median, err := oracle.NewMedianAggregator(logger, marketmap)
			require.NoError(t, err)

			aggFn := median.AggregateFn()
			prices := aggFn(tc.pricesPerProvider)
			require.Equal(t, len(tc.expected), len(prices))
			for cp, expectedPrice := range tc.expected {
				actualPrice, ok := prices[cp]
				require.True(t, ok)

				verifyPrice(t, expectedPrice, actualPrice)
			}
		})
	}
}

func TestCalculateConvertedPrices(t *testing.T) {
	testCases := []struct {
		name          string
		outcome       mmtypes.Ticker
		operations    mmtypes.Path
		medians       types.TickerPrices
		expected      *big.Int
		expectedError bool
	}{
		{
			name:       "no operations",
			outcome:    constants.BITCOIN_USD,
			operations: mmtypes.Path{},
			medians: types.TickerPrices{
				constants.BITCOIN_USD: createPrice(40_000, 8),
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name:    "not enough median prices",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.BITCOIN_USDT: createPrice(40_000, 8),
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name:    "successful conversion directly from a median price",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.BITCOIN_USD: createPrice(40_000, oracle.ScaledDecimals),
			},
			expected:      createPrice(40_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from converted prices",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.BITCOIN_USDT: createPrice(40_000, oracle.ScaledDecimals),
				constants.USDT_USD:     createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from converted prices with an inverted price",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: usdtbtc.CurrencyPair,
						Invert:       true,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				usdtbtc:            createPrice(0.000025, oracle.ScaledDecimals),
				constants.USDT_USD: createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from with reasonably small numbers",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.BITCOIN_USDT: createPrice(0.0000000000004, oracle.ScaledDecimals), // 4e-13
				constants.USDT_USD:     createPrice(0.0000000000012, oracle.ScaledDecimals), // 1.2e-12
			},
			expected:      createPrice(0.00000000000000000000000048, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from with reasonably large numbers",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.BITCOIN_USDT: createPrice(40_000_000_000_000_000, oracle.ScaledDecimals), // 4e16 + scaled to 40 decimals
				constants.USDT_USD:     createPrice(1_200_000, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000_000_000_000_000_000_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion with 3 conversion operations and an inversion",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.ETHEREUM_BITCOIN.CurrencyPair,
						Invert:       true,
					},
					{
						CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.ETHEREUM_BITCOIN: createPrice(5, oracle.ScaledDecimals-2),
				constants.ETHEREUM_USDT:    createPrice(2000, oracle.ScaledDecimals),
				constants.USDT_USD:         createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "path contains a price of 0 at the start",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.ETHEREUM_BITCOIN.CurrencyPair,
						Invert:       true,
					},
					{
						CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.ETHEREUM_BITCOIN: createPrice(0, oracle.ScaledDecimals),
				constants.ETHEREUM_USDT:    createPrice(2000, oracle.ScaledDecimals),
				constants.USDT_USD:         createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      big.NewInt(0),
			expectedError: false,
		},
		{
			name:    "path contains a price of 0 in the middle",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.ETHEREUM_BITCOIN.CurrencyPair,
						Invert:       true,
					},
					{
						CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
					},
					{
						CurrencyPair: constants.USDT_USD.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.ETHEREUM_BITCOIN: createPrice(20, oracle.ScaledDecimals),
				constants.ETHEREUM_USDT:    createPrice(0, oracle.ScaledDecimals),
				constants.USDT_USD:         createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      big.NewInt(0),
			expectedError: false,
		},
		{
			name:    "conversion path is invalid",
			outcome: constants.BITCOIN_USD,
			operations: mmtypes.Path{
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: constants.ETHEREUM_BITCOIN.CurrencyPair,
						Invert:       true,
					},
					{
						CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
					},
				},
			},
			medians: types.TickerPrices{
				constants.ETHEREUM_BITCOIN: createPrice(20, oracle.ScaledDecimals),
				constants.ETHEREUM_USDT:    createPrice(2000, oracle.ScaledDecimals),
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aggregator, err := oracle.NewMedianAggregator(logger, marketmap)
			require.NoError(t, err)

			price, err := aggregator.CalculateConvertedPrice(tc.outcome, tc.operations, tc.medians)
			if tc.expectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			verifyPrice(t, tc.expected, price)
		})
	}
}

// verifyPrice verifies that the expected price matches the actual price within an acceptable delta.
func verifyPrice(t *testing.T, expected, actual *big.Int) {
	t.Helper()

	zero := big.NewInt(0)
	if expected.Cmp(zero) == 0 {
		require.Equal(t, zero, actual)
		return
	}

	var diff *big.Float
	if expected.Cmp(actual) > 0 {
		diff = new(big.Float).Sub(new(big.Float).SetInt(expected), new(big.Float).SetInt(actual))
	} else {
		diff = new(big.Float).Sub(new(big.Float).SetInt(actual), new(big.Float).SetInt(expected))
	}

	scaledDiff := new(big.Float).Quo(diff, new(big.Float).SetInt(expected))
	delta, _ := scaledDiff.Float64()
	t.Logf("expected price: %s; actual price: %s; diff %s", expected.String(), actual.String(), diff.String())
	t.Logf("acceptable delta: %.25f; actual delta: %.25f", acceptableDelta, delta)

	switch {
	case delta == 0:
		// If the difference between the expected and actual price is 0, the prices match.
		// No need for a delta comparison.
		return
	case delta <= acceptableDelta:
		// If the difference between the expected and actual price is within the acceptable delta,
		// the prices match.
		return
	default:
		// If the difference between the expected and actual price is greater than the acceptable delta,
		// the prices do not match.
		require.Fail(t, "expected price does not match the actual price; delta is too large")
	}
}

// createPrice creates a price with the given number of decimals.
func createPrice(price float64, decimals uint64) *big.Int {
	// Convert the price to a big float so we can perform the multiplication.
	floatPrice := big.NewFloat(price)

	// Scale the price and convert it to a big.Int.
	one := oracle.ScaledOne(decimals)
	scaledPrice := new(big.Float).Mul(floatPrice, new(big.Float).SetInt(one))
	intPrice, _ := scaledPrice.Int(nil)
	return intPrice
}
