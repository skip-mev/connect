package uniswap

import (
	"math/big"

	"github.com/skip-mev/slinky/pkg/math/oracle"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// ConvertSquareRootX96Price converts the slot 0 sqrtPriceX96 value to a price. Note that this
// price is not scaled to the token decimals.
func ConvertSquareRootX96Price(
	sqrtPriceX96 *big.Int,
) *big.Float {
	// Convert the original sqrtPriceX96 to a big float to retain precision when dividing.
	sqrtPriceX96Float := new(big.Float).SetInt(sqrtPriceX96)

	// x96Float is the fixed-point precision for Uniswap V3 prices. This is equal to 2^96.
	x96Float := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil),
	)

	// Divide the sqrtPriceX96 by the fixed-point precision to get the price.
	sqrtPriceFloat := new(big.Float).Quo(sqrtPriceX96Float, x96Float)

	// Square the price to get the final result. We multiply the prices here instead of converting
	// to big.Int to retain precision.
	return new(big.Float).Mul(sqrtPriceFloat, sqrtPriceFloat)
}

// ScalePrice scales the price to the respective token decimals.
func ScalePrice(
	ticker mmtypes.Ticker,
	cfg PoolConfig,
	price *big.Float,
) *big.Float {
	// Determine the scaling factor for the price.
	diff := cfg.BaseDecimals - cfg.QuoteDecimals
	one := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(diff)), nil),
	)

	// Adjust the price based on the difference between the token decimals.
	scaleAdjustedPrice := new(big.Float).Quo(price, one)
	if !cfg.Invert {
		return scaleAdjustedPrice
	}

	// Invert the price if necessary.
	return InvertCurrencyPairPrice(scaleAdjustedPrice, ticker.Decimals)
}

// InvertCurrencyPairPrice inverts a price by performing the following operation:
// 1. 1 / price
// 2. Scale the result by the number of decimals
// 3. Convert the result to a big.Int
//
// NOTE: This function should only be used on prices that have already been scaled
// to the standard number of decimals.
func InvertCurrencyPairPrice(price *big.Float, decimals uint64) *big.Float {
	one := new(big.Float).SetInt(oracle.ScaledOne(decimals))

	// Convert the price to a big.Float so we can perform the division
	// and then convert the result back to a big.Int This operation is
	// the equivalent of 1 / price.
	ratio := new(big.Float).Quo(one, price)

	// Scale the ratio by the number of decimals.
	return ratio
}
