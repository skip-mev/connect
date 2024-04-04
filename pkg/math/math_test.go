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

func TestAbs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		val      int
		expected int
	}{
		{
			name:     "positive",
			val:      1,
			expected: 1,
		},
		{
			name:     "negative",
			val:      -1,
			expected: 1,
		},
		{
			name:     "zero",
			val:      0,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := math.Abs(tc.val)
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
