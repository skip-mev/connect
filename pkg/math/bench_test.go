package math_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/skip-mev/connect/v2/pkg/math"
)

func BenchmarkFloat64StringToBigInt(b *testing.B) {
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
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, _ = math.Float64StringToBigInt(tc.input, tc.base)
			}
		})
	}
}

func BenchmarkFloat64ToBigInt(b *testing.B) {
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
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = math.Float64ToBigInt(tc.input, tc.base)
			}
		})
	}
}

func BenchmarkBigFloatToBigInt(b *testing.B) {
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
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = math.BigFloatToBigInt(tc.input, tc.base)
			}
		})
	}
}
