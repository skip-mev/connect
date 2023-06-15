package providers_test

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/providers"
	"github.com/stretchr/testify/require"
)

func TestFloat64ToUint256(t *testing.T) {
	testCases := []struct {
		name     string
		input    float64
		base     int
		expected *uint256.Int
	}{
		{
			"zero",
			0,
			6,
			uint256.NewInt(0),
		},
		{
			"one",
			1,
			6,
			uint256.NewInt(1000000),
		},
		{
			"one point one",
			1.1,
			6,
			uint256.NewInt(1100000),
		},
		{
			"many decimal points",
			1.123456789,
			6,
			uint256.NewInt(1123456),
		},
		{
			"random big number with many decimal points",
			123456789.123456789,
			6,
			uint256.NewInt(123456789123456),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := providers.Float64ToUint256(tc.input, tc.base)
			require.Equal(t, tc.expected, result)
		})
	}
}
