package voteweighted

import (
	"math/big"

	"cosmossdk.io/math"
)

var (
	Zero       = big.NewInt(0)
	Two        = big.NewInt(2)
	OneMillion = big.NewInt(1_000_000)
)

// SignsDiffer determines whether newPrice and newPrime move in different directions relative to old.
func SignsDiffer(old *big.Int, newPrice *big.Int, newPrime *big.Int) bool {
	newSign := (&big.Int{}).Sub(old, newPrice).Sign()
	newPrimeSign := (&big.Int{}).Sub(old, newPrime).Sign()
	return newSign != newPrimeSign
}

// CalculateRelativeTicks computes the ticks of newPrime relative to the tick size of newPrice.
func CalculateRelativeTicks(old *big.Int, newPrice *big.Int, newPrime *big.Int, ppm *big.Int) *big.Int {
	/*
		1_000_000 * abs(newPrice-newPrime)
		_____________________________
				(old*ppm)
	*/

	// abs(new-newPrime)
	delta := (&big.Int{}).Abs((&big.Int{}).Sub(newPrice, newPrime))
	// 1_000_000 * delta
	oneMillionDelta := (&big.Int{}).Mul(OneMillion, delta)
	// old * ppm
	oldPpm := (&big.Int{}).Mul(old, ppm)
	// returns (delta * 1_000_000) / (old * ppm)
	return (&big.Int{}).Quo(oneMillionDelta, oldPpm)
}

// CalculateTicks computes the ticks of newPrice relative to old.
func CalculateTicks(old *big.Int, newPrice *big.Int, ppm *big.Int) *big.Int {
	/*
		1_000_000 * abs(old-newPrice)
		________________________
			   (old*ppm)
	*/

	// abs(old-newPrice)
	delta := (&big.Int{}).Abs((&big.Int{}).Sub(old, newPrice))
	// 1_000_000 * delta
	oneMillionDelta := (&big.Int{}).Mul(OneMillion, delta)
	// old * ppm
	oldPpm := (&big.Int{}).Mul(old, ppm)
	// returns (delta * 1_000_000) / (old * ppm)
	return (&big.Int{}).Quo(oneMillionDelta, oldPpm)
}

// ThresholdWeightCalc computes the amount of weight taken into account for the given price update.
// It follows an algorithm which requires increasing vote price correlation as the percent change in
// the price increases.
func ThresholdWeightCalc(
	currentPrice *big.Int,
	proposedPrice *big.Int,
	ppm *big.Int,
	priceInfo PriceInfo,
) math.Int {
	totalWeight := math.NewInt(0)
	proposedTicks := CalculateTicks(currentPrice, proposedPrice, ppm)
	for _, validatorPrice := range priceInfo.Prices {
		priceTicks := CalculateTicks(currentPrice, validatorPrice.Price, ppm)
		contributedWeight := math.NewInt(0)
		relativeTicks := CalculateRelativeTicks(currentPrice, proposedPrice, validatorPrice.Price, ppm)
		// If priceTicks == 0
		if priceTicks.Cmp(Zero) == 0 {
			// If proposedTicks < 2
			if proposedTicks.Cmp(Two) == -1 {
				contributedWeight = validatorPrice.VoteWeight.Mul(math.NewInt(2 - proposedTicks.Int64()))
			}
			// If relativeTicks <= sqrt(priceTicks) && price direction is the same
		} else if relativeTicks.Cmp(new(big.Int).Sqrt(priceTicks)) <= 0 &&
			!SignsDiffer(currentPrice, proposedPrice, validatorPrice.Price) {
			contributedWeight = validatorPrice.VoteWeight
		}
		totalWeight = totalWeight.Add(contributedWeight)
	}
	return totalWeight
}
