package voteweighted_test

import (
	"math/big"
	"testing"

	"github.com/skip-mev/slinky/pkg/math/voteweighted"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
)

func TestThresholdWeightCalc(t *testing.T) {
	testCases := []struct {
		name           string
		currentPrice   *big.Int
		proposedPrice  *big.Int
		ppm            *big.Int
		priceInfo      voteweighted.PriceInfo
		expectedWeight math.Int
	}{
		{
			"all prices equal 2 ticks",
			big.NewInt(10),
			big.NewInt(20),
			big.NewInt(500_000),
			voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{math.NewInt(2), big.NewInt(20)},
					{math.NewInt(2), big.NewInt(20)},
					{math.NewInt(2), big.NewInt(20)},
				},
			},
			math.NewInt(6),
		},
		{
			"all prices equal 1 tick",
			big.NewInt(10),
			big.NewInt(15),
			big.NewInt(500_000),
			voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{math.NewInt(2), big.NewInt(15)},
					{math.NewInt(2), big.NewInt(15)},
					{math.NewInt(2), big.NewInt(15)},
				},
			},
			math.NewInt(6),
		},
		{
			"all prices equal 10 ticks",
			big.NewInt(10),
			big.NewInt(60),
			big.NewInt(500_000),
			voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{math.NewInt(2), big.NewInt(60)},
					{math.NewInt(2), big.NewInt(60)},
					{math.NewInt(2), big.NewInt(60)},
				},
			},
			math.NewInt(6),
		},
		{
			"prices spread 10 ticks",
			big.NewInt(10),
			big.NewInt(60),
			big.NewInt(500_000),
			voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{math.NewInt(1), big.NewInt(90)},
					{math.NewInt(3), big.NewInt(30)},
					{math.NewInt(5), big.NewInt(60)},
				},
			},
			math.NewInt(5),
		},
		{
			"prices spread 0 ticks",
			big.NewInt(10),
			big.NewInt(14),
			big.NewInt(500_000),
			voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{math.NewInt(1), big.NewInt(20)},
					{math.NewInt(3), big.NewInt(9)},
					{math.NewInt(5), big.NewInt(15)},
				},
			},
			math.NewInt(12),
		},
		{
			"prices spread 0 ticks",
			big.NewInt(10),
			big.NewInt(14),
			big.NewInt(500_000),
			voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{math.NewInt(1), big.NewInt(20)},
					{math.NewInt(3), big.NewInt(9)},
					{math.NewInt(5), big.NewInt(15)},
				},
			},
			math.NewInt(12),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.expectedWeight.Int64(),
				voteweighted.ThresholdWeightCalc(
					tc.currentPrice,
					tc.proposedPrice,
					tc.ppm,
					tc.priceInfo,
				).Int64(),
			)
		})
	}
}
