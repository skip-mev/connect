package oracle

import (
	"fmt"
	"math/big"
)

// ScaledDecimals is the standard number of decimal places each price will be converted to
// during the conversion process.
const ScaledDecimals = 36

// ScaleUpCurrencyPairPrice scales a price up to the standard number of decimals by performing the
// following operation:
// 1. price * 10^(ScaledDecimals - decimals)
// 2. Convert the result to a big.Int
//
// NOTE: This function should only be used on prices that have not already been scaled to the
// standard number of decimals. We scale the price to the standard number of decimals for ease
// of comparison.
func ScaleUpCurrencyPairPrice(decimals uint64, price *big.Int) (*big.Int, error) {
	if decimals > ScaledDecimals {
		return nil, fmt.Errorf("cannot scale up price with more decimals than the standard: max=%d, current=%d", ScaledDecimals, decimals)
	}

	diff := ScaledDecimals - decimals
	exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(diff)), nil)
	return new(big.Int).Mul(price, exp), nil
}

// ScaleDownCurrencyPairPrice scales a price down to the standard number of decimals by performing the
// following operation:
// 1. price / 10^(ScaledDecimals - decimals)
// 2. Convert the result to a big.Int
//
// NOTE: This function should only be used on prices that have already been scaled to the standard
// number of decimals. The output of this returns the price to its expected number of decimals.
func ScaleDownCurrencyPairPrice(decimals uint64, price *big.Int) (*big.Int, error) {
	if decimals > ScaledDecimals {
		return nil, fmt.Errorf("cannot scale down price with more decimals than the standard: max=%d, current=%d", ScaledDecimals, decimals)
	}

	diff := ScaledDecimals - decimals
	exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(diff)), nil)
	return new(big.Int).Div(price, exp), nil
}

// InvertCurrencyPairPrice inverts a price by performing the following operation:
// 1. 1 / price
// 2. Scale the result by the number of decimals
// 3. Convert the result to a big.Int
//
// NOTE: This function should only be used on prices that have already been scaled
// to the standard number of decimals.
func InvertCurrencyPairPrice(price *big.Int, decimals uint64) *big.Int {
	one := ScaledOne(decimals)

	// Convert the price to a big.Float so we can perform the division
	// and then convert the result back to a big.Int This operation is
	// the equivalent of 1 / price.
	ratio := new(big.Float).Quo(new(big.Float).SetInt(one), new(big.Float).SetInt(price))

	// Scale the ratio by the number of decimals.
	scaledRatio := new(big.Float).Mul(ratio, new(big.Float).SetInt(one))

	// Convert the scaled ratio back to a big.Int
	inverted, _ := scaledRatio.Int(nil)
	return inverted
}

// ScaledOne returns a big.Int that represents the number 1 scaled to the standard
// number of decimals.
func ScaledOne(decimals uint64) *big.Int {
	return big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
}
