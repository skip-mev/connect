package voteweighted

import (
	"cosmossdk.io/math"
	"math/big"
)

var (
	Zero       = big.NewInt(0)
	Two        = big.NewInt(2)
	OneMillion = big.NewInt(1_000_000)
)

// CalculateRelativeTicks computes the ticks of newPrime relative to the tick size of new
func CalculateRelativeTicks(old *big.Int, new *big.Int, newPrime *big.Int, ppm *big.Int) *big.Int {
	/*
		1_000_000 * abs(new-newPrime)
		_____________________________
				(old*ppm)
	*/

	// abs(new-newPrime)
	delta := new(big.Int).Abs(new(big.Int).Sub(new, newPrime))
	// 1_000_000 * delta
	oneMillionDelta := new(big.Int).Mul(OneMillion, delta)
	// old * ppm
	oldPpm := new(big.Int).Mul(old, ppm)
	// returns (delta * 1_000_000) / (old * ppm)
	return new(big.Int).Quo(oneMillionDelta, oldPpm)
}

// CalculateTicks computes the ticks of new relative to old
func CalculateTicks(old *big.Int, new *big.Int, ppm *big.Int) *big.Int {
	/*
		1_000_000 * abs(old-new)
		________________________
			   (old*ppm)
	*/

	// abs(old-new)
	delta := new(big.Int).Abs(new(big.Int).Sub(old, new))
	// 1_000_000 * delta
	oneMillionDelta := new(big.Int).Mul(OneMillion, delta)
	// old * ppm
	oldPpm := new(big.Int).Mul(old, ppm)
	// returns (delta * 1_000_000) / (old * ppm)
	return new(big.Int).Quo(oneMillionDelta, oldPpm)
}

// ThresholdWeightCalc does some magic. Do not question the magic. Do not refactor the magic.
// Definitely don't break the magic.
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
			// If relativeTicks <= sqrt(priceTicks)
		} else if relativeTicks.Cmp(new(big.Int).Sqrt(priceTicks)) <= 0 {
			contributedWeight = validatorPrice.VoteWeight
		}
		totalWeight = totalWeight.Add(contributedWeight)
	}
	return totalWeight

}
