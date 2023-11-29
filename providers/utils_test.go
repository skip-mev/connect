package providers_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers"
)

func TestFloat64ToBigInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    float64
		base     int
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
			result := providers.Float64ToBigInt(tc.input, tc.base)
			require.Equal(t, tc.expected, result)
		})
	}
}
