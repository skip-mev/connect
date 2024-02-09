package oracle_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	// acceptableDelta is the acceptable difference between the expected and actual price.
	// In this case, we use a delta of 1e-10. This means we will accept any price that is
	// within 1e-10 of the expected price.
	acceptableDelta = 1e-10

	logger, _ = zap.NewDevelopment()
	cfg       = config.AggregateMarketConfig{
		Feeds: map[string]config.FeedConfig{
			"BITCOIN/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDT": {
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"USDT/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
			},
			"ETHEREUM/USDT": {
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"BITCOIN/ETHEREUM": {
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"),
			},
			"USDT/ETHEREUM": {
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "ETHEREUM"),
			},
			"ETHEREUM/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
		},
		AggregatedFeeds: map[string]config.AggregateFeedConfig{
			"BITCOIN/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				Conversions: []config.Conversions{
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						},
					},
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						},
					},
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"),
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						},
					},
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("USDT", "ETHEREUM"),
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
						},
					},
				},
			},
			"ETHEREUM/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
				Conversions: []config.Conversions{
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						},
					},
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("USDT", "ETHEREUM"),
							Invert:       true,
						},
						{
							CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						},
					},
				},
			},
		},
	}
)

func TestMedian(t *testing.T) {
	testCases := []struct {
		name              string
		pricesPerProvider map[string]map[oracletypes.CurrencyPair]*big.Int
		expected          map[oracletypes.CurrencyPair]*big.Int
	}{
		{
			name: "no prices",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "single resolved price",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(40_000, 8),
				},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(40_000, 8),
			},
		},
		{
			name: "must convert to get a single final price",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): createPrice(40_000, 8),
					oracletypes.NewCurrencyPair("USDT", "USD"):     createPrice(1.1, 8),
				},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(44_000, 8),
			},
		},
		{
			name: "calculates median price between two separate conversions",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  createPrice(40_000, 8),
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): createPrice(40_000, 8),
					oracletypes.NewCurrencyPair("USDT", "USD"):     createPrice(1.1, 8),
				},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(42_000, 8), // median average of 40_000 and 44_000
			},
		},
		{
			name: "calculates median price between three separate conversions",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):      createPrice(40_000, 8),
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"):     createPrice(40_000, 8),
					oracletypes.NewCurrencyPair("USDT", "USD"):         createPrice(1.1, 8),
					oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"): createPrice(22, 18),
					oracletypes.NewCurrencyPair("ETHEREUM", "USDT"):    createPrice(2000, 8),
				},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"):  createPrice(44_000, 8), // median average of 40_000, 44_000, and 48_400
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"): createPrice(2_200, 8),
			},
		},
		{
			name: "calculates median price with an inverted price",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {
					oracletypes.NewCurrencyPair("USDT", "ETHEREUM"): createPrice(0.0005, 18),
					oracletypes.NewCurrencyPair("USDT", "USD"): createPrice(1.1,
						8),
				},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"): createPrice(2_200, 8),
			},
		},
		{
			name: "calculates median price with a price of 0",
			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
				"coinbase": {
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(0, 8),
				},
			},
			expected: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(0, 8),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			median, err := oracle.NewMedianAggregator(logger, cfg)
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
		outcome       oracletypes.CurrencyPair
		operations    []config.Conversion
		medians       map[oracletypes.CurrencyPair]*big.Int
		expected      *big.Int
		expectedError bool
	}{
		{
			name:       "no operations",
			outcome:    oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(40_000, 8),
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name:    "not enough median prices",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): createPrice(40_000, 8),
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name:    "successful conversion directly from a median price",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): createPrice(40_000, oracle.ScaledDecimals),
			},
			expected:      createPrice(40_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from converted prices",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): createPrice(40_000, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("USDT", "USD"):     createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from converted prices with an inverted price",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN"),
					Invert:       true,
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("USDT", "BITCOIN"): createPrice(0.000025, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("USDT", "USD"):     createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from with reasonably small numbers",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): createPrice(0.0000000000004, oracle.ScaledDecimals), // 4e-13
				oracletypes.NewCurrencyPair("USDT", "USD"):     createPrice(0.0000000000012, oracle.ScaledDecimals), // 1.2e-12
			},
			expected:      createPrice(0.00000000000000000000000048, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion from with reasonably large numbers",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): createPrice(40_000_000_000_000_000, oracle.ScaledDecimals), // 4e16 + scaled to 40 decimals
				oracletypes.NewCurrencyPair("USDT", "USD"):     createPrice(1_200_000, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000_000_000_000_000_000_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "successful conversion with 3 conversion operations",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"): createPrice(20, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT"):    createPrice(2000, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("USDT", "USD"):         createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      createPrice(48_000, oracle.ScaledDecimals),
			expectedError: false,
		},
		{
			name:    "path contains a price of 0 at the start",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"): createPrice(0, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT"):    createPrice(2000, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("USDT", "USD"):         createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      big.NewInt(0),
			expectedError: false,
		},
		{
			name:    "path contains a price of 0 in the middle",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"): createPrice(20, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT"):    createPrice(0, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("USDT", "USD"):         createPrice(1.2, oracle.ScaledDecimals),
			},
			expected:      big.NewInt(0),
			expectedError: false,
		},
		{
			name:    "conversion path is invalid",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "ETHEREUM"): createPrice(20, oracle.ScaledDecimals),
				oracletypes.NewCurrencyPair("ETHEREUM", "USDT"):    createPrice(2000, oracle.ScaledDecimals),
			},
			expected:      nil,
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aggregator, err := oracle.NewMedianAggregator(logger, cfg)
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
func createPrice(price float64, decimals int64) *big.Int {
	// Convert the price to a big float so we can perform the multiplication.
	floatPrice := big.NewFloat(price)

	// Scale the price and convert it to a big.Int.
	one := oracle.ScaledOne(decimals)
	scaledPrice := new(big.Float).Mul(floatPrice, new(big.Float).SetInt(one))
	intPrice, _ := scaledPrice.Int(nil)
	return intPrice
}
