package math_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/pkg/math"
)

func TestMin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		vals []int
		min  int
	}{
		{
			name: "one value",
			vals: []int{1},
			min:  1,
		},
		{
			name: "two values",
			vals: []int{1, 2},
			min:  1,
		},
		{
			name: "three values",
			vals: []int{1, 2, 3},
			min:  1,
		},
		{
			name: "five values, negative",
			vals: []int{1, 2, 3, 4, -5},
			min:  -5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min := math.Min(tc.vals...)
			if min != tc.min {
				t.Errorf("expected %d, got %d", tc.min, min)
			}
		})
	}
}

func TestFloat64ToBigInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    float64
		decimals int64
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
			result := math.Float64ToBigInt(tc.input, tc.decimals)
			require.Equal(t, tc.expected, result)
		})
	}
}
