package math_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/pkg/math"
)

func TestMin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		vals     []int
		expected int
	}{
		{
			name:     "one value",
			vals:     []int{1},
			expected: 1,
		},
		{
			name:     "two values",
			vals:     []int{1, 2},
			expected: 1,
		},
		{
			name:     "three values",
			vals:     []int{1, 2, 3},
			expected: 1,
		},
		{
			name:     "five values, negative",
			vals:     []int{1, 2, 3, 4, -5},
			expected: -5,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := math.Min(tc.vals...)
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestFloat64StringToBigInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		base     uint64
		expected *big.Int
	}{
		{
			"zero",
			strconv.FormatFloat(0, 'g', -1, 64),
			6,
			big.NewInt(0),
		},
		{
			"one",
			strconv.FormatFloat(1, 'g', -1, 64),
			6,
			big.NewInt(1000000),
		},
		{
			"one point one",
			strconv.FormatFloat(1.1, 'g', -1, 64),
			6,
			big.NewInt(1100000),
		},
		{
			"many decimal points",
			strconv.FormatFloat(1.123456789, 'g', -1, 64),
			6,
			big.NewInt(1123456),
		},
		{
			"random big number with many decimal points",
			strconv.FormatFloat(123456789.123456789, 'g', -1, 64),
			6,
			big.NewInt(123456789123456),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := math.Float64StringToBigInt(tc.input, tc.base)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestFloat64ToBigInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    float64
		base     uint64
		expected *big.Int
	}{
		{
			"zero",
			0,
			6,
			big.NewInt(0),
		},
		{
			"one",
			1,
			6,
			big.NewInt(1000000),
		},
		{
			"one point one",
			1.1,
			6,
			big.NewInt(1100000),
		},
		{
			"many decimal points",
			1.123456789,
			6,
			big.NewInt(1123456),
		},
		{
			"random big number with many decimal points",
			123456789.123456789,
			6,
			big.NewInt(123456789123456),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := math.Float64ToBigInt(tc.input, tc.base)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestBigFloatToBigInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    *big.Float
		base     uint64
		expected *big.Int
	}{
		{
			"zero",
			big.NewFloat(0),
			6,
			big.NewInt(0),
		},
		{
			"one",
			big.NewFloat(1),
			6,
			big.NewInt(1000000),
		},
		{
			"one point one",
			big.NewFloat(1.1),

			6,
			big.NewInt(1100000),
		},
		{
			"many decimal points",
			big.NewFloat(1.123456789),
			6,
			big.NewInt(1123456),
		},
		{
			"random big number with many decimal points",
			big.NewFloat(123456789.123456789),
			6,
			big.NewInt(123456789123456),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := math.BigFloatToBigInt(tc.input, tc.base)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestFloat64StringToBigFloat(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  *big.Float
		err  bool
	}{
		{
			name: "zero",
			in:   "0",
			out:  big.NewFloat(0),
			err:  false,
		},
		{
			name: "1",
			in:   "1",
			out:  big.NewFloat(1),
			err:  false,
		},
		{
			name: "value that has more 0s than decimals",
			in:   "0.0000000000000001", // 1e-16
			out:  big.NewFloat(1e-16),
			err:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			out, err := math.Float64StringToBigFloat(tc.in)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.out.SetPrec(uint(20)), out.SetPrec(uint(20)))
			}
		})
	}
}

func TestScaleBigFloat(t *testing.T) {
	testCases := []struct {
		name     string
		in       *big.Float
		decimals uint64
		out      *big.Float
	}{
		{
			name:     "zero",
			in:       big.NewFloat(0),
			decimals: 6,
			out:      big.NewFloat(0),
		},
		{
			name:     "1",
			in:       big.NewFloat(1),
			decimals: 6,
			out:      big.NewFloat(1e6),
		},
		{
			name:     "value that has more 0s than decimals",
			in:       big.NewFloat(1e-16),
			decimals: 6,
			out:      big.NewFloat(1e-10),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := math.ScaleBigFloat(tc.in, tc.decimals)
			require.Equal(t, tc.out.SetPrec(uint(20)), out.SetPrec(uint(20)))
		})
	}
}

func TestCalculateMedian(t *testing.T) {
	testCases := []struct {
		name     string
		values   []*big.Float
		expected *big.Float
	}{
		{
			name:     "do nothing for nil slice",
			values:   nil,
			expected: nil,
		},
		{
			name: "calculate median for even number of values",
			values: []*big.Float{
				big.NewFloat(-2),
				big.NewFloat(0),
				big.NewFloat(10),
				big.NewFloat(100),
			},
			expected: big.NewFloat(5),
		},
		{
			name: "calculate median for odd number of values",
			values: []*big.Float{
				big.NewFloat(10),
				big.NewFloat(-2),
				big.NewFloat(100),
				big.NewFloat(0),
				big.NewFloat(0),
			},
			expected: big.NewFloat(0),
		},
		{
			"calculates median for even number of values with decimals",
			[]*big.Float{
				big.NewFloat(-2),
				big.NewFloat(0),
				big.NewFloat(0),
				big.NewFloat(1),
				big.NewFloat(10),
				big.NewFloat(100),
			},
			big.NewFloat(0.5),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, math.CalculateMedian(tc.values))
		})
	}
}

func TestSortBigInts(t *testing.T) {
	testCases := []struct {
		name     string
		values   []*big.Float
		expected []*big.Float
	}{
		{
			name: "do nothing for nil slice",
		},
		{
			name: "sort a slice",
			values: []*big.Float{
				big.NewFloat(10),
				big.NewFloat(-2),
				big.NewFloat(100),
				big.NewFloat(0),
				big.NewFloat(0),
			},
			expected: []*big.Float{
				big.NewFloat(-2),
				big.NewFloat(0),
				big.NewFloat(0),
				big.NewFloat(10),
				big.NewFloat(100),
			},
		},
		{
			name: "do nothing for same values",
			values: []*big.Float{
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(-2),
				big.NewFloat(100),
				big.NewFloat(0),
				big.NewFloat(0),
			},
			expected: []*big.Float{
				big.NewFloat(-2),
				big.NewFloat(0),
				big.NewFloat(0),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(10),
				big.NewFloat(100),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			math.SortBigFloats(tc.values)
			require.Equal(t, tc.expected, tc.values)
		})
	}
}
