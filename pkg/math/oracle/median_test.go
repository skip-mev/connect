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
		},
		AggregatedFeeds: map[string]config.AggregateFeedConfig{
			"BITCOIN/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				Conversions: [][]config.Conversion{
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
				},
			},
			"ETHEREUM/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
				Conversions: [][]config.Conversion{
					{
						{
							CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
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

// func TestMedian(t *testing.T) {
// 	testCases := []struct {
// 		name              string
// 		pricesPerProvider map[string]map[oracletypes.CurrencyPair]*big.Int
// 		expected          map[oracletypes.CurrencyPair]*big.Int
// 	}{
// 		{
// 			name: "no prices",
// 			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
// 				"coinbase": {},
// 			},
// 			expected: map[oracletypes.CurrencyPair]*big.Int{},
// 		},
// 		{
// 			name: "single resolved price",
// 			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
// 				"coinbase": {
// 					oracletypes.NewCurrencyPair("BITCOIN", "USD"): big.NewInt(100),
// 				},
// 			},
// 			expected: map[oracletypes.CurrencyPair]*big.Int{
// 				oracletypes.NewCurrencyPair("BITCOIN", "USD"): big.NewInt(100),
// 			},
// 		},
// 		{
// 			name: "must convert to get a single final price",
// 			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
// 				"coinbase": {
// 					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): big.NewInt(100),
// 					oracletypes.NewCurrencyPair("USDT", "USD"):     big.NewInt(2),
// 				},
// 			},
// 			expected: map[oracletypes.CurrencyPair]*big.Int{
// 				oracletypes.NewCurrencyPair("BITCOIN", "USD"): big.NewInt(200),
// 			},
// 		},
// 		{
// 			name: "calculates median price between two separate conversions",
// 			pricesPerProvider: map[string]map[oracletypes.CurrencyPair]*big.Int{
// 				"coinbase": {
// 					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  big.NewInt(100),
// 					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): big.NewInt(100),
// 					oracletypes.NewCurrencyPair("USDT", "USD"):     big.NewInt(2),
// 				},
// 			},
// 			expected: map[oracletypes.CurrencyPair]*big.Int{
// 				oracletypes.NewCurrencyPair("BITCOIN", "USD"): big.NewInt(150),
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			median, err := oracle.NewMedianAggregator(logger, cfg)
// 			require.NoError(t, err)

// 			aggFn := median.AggregateFn()
// 			prices := aggFn(tc.pricesPerProvider)
// 			require.Equal(t, tc.expected, prices)
// 		})
// 	}
// }

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
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): big.NewInt(100),
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
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): big.NewInt(100),
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
				oracletypes.NewCurrencyPair("BITCOIN", "USD"): big.NewInt(100),
			},
			expected:      big.NewInt(100),
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
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): big.NewInt(100),
				oracletypes.NewCurrencyPair("USDT", "USD"):     big.NewInt(2),
			},
			expected:      big.NewInt(200),
			expectedError: false,
		},
		{
			name:    "successful conversion from converted prices with an inverted feed",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{

				{
					CurrencyPair: oracletypes.NewCurrencyPair("USD", "BITCOIN"),
					Invert:       true,
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("USD", "BITCOIN"): big.NewInt(1_000_000),
			},
			expected:      big.NewInt(100),
			expectedError: false,
		},
		{
			name:    "successful conversion from converted prices with an inverted feed and a conversion",
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
				oracletypes.NewCurrencyPair("USDT", "BITCOIN"): big.NewInt(1_000_000),
				oracletypes.NewCurrencyPair("USDT", "USD"):     big.NewInt(2),
			},
			expected:      big.NewInt(200),
			expectedError: false,
		},
		{
			name:    "multiple non-inverted conversions",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
				{
					CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USDC"),
				},
				{
					CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
				},
			},
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"): big.NewInt(100),
				oracletypes.NewCurrencyPair("USDT", "USDC"):    big.NewInt(2),
				oracletypes.NewCurrencyPair("USDC", "USD"):     big.NewInt(2),
			},
			expected:      big.NewInt(400),
			expectedError: false,
		},
		{
			name:    "varying decimal points",
			outcome: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			operations: []config.Conversion{
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
			medians: map[oracletypes.CurrencyPair]*big.Int{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"):  createPrice(40_000, 8, false),
				oracletypes.NewCurrencyPair("USDT", "ETHEREUM"): createPrice(2_000, 18, true),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"):  createPrice(2_000, 18, false),
			},
			expected: createPrice(40_000, 8, false), // 40_000 * (1 / 2_000) * 2_000
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
			require.Equal(t, tc.expected, price)
		})
	}
}

func createPrice(price, decimals int64, invert bool) *big.Int {
	scaledOne := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	nonInvertedPrice := new(big.Int).Mul(big.NewInt(price), scaledOne)
	if !invert {
		return nonInvertedPrice
	}

	return new(big.Int).Div(scaledOne, nonInvertedPrice)
}
