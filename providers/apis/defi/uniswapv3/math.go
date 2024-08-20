package uniswapv3

import (
	"math/big"

	"github.com/skip-mev/connect/v2/pkg/math"
)

// ConvertSquareRootX96Price converts the slot 0 sqrtPriceX96 value to a price. Note that this
// price is not scaled to the token decimals. This calculation is equivalent to:
//
// price = (sqrtPriceX96 / 2^96) ^ 2.
func ConvertSquareRootX96Price(
	sqrtPriceX96 *big.Int,
) *big.Float {
	// Convert the original sqrtPriceX96 to a big float to retain precision when dividing.
	sqrtPriceX96Float := new(big.Float).SetInt(sqrtPriceX96)

	// x96Float is the fixed-point precision for Uniswap V3 prices. This is equal to 2^96.
	x96Float := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil),
	)

	// Divide the sqrtPriceX96 by the fixed-point precision.
	sqrtPriceFloat := new(big.Float).Quo(sqrtPriceX96Float, x96Float)

	// Square the price to get the final result.
	return new(big.Float).Mul(sqrtPriceFloat, sqrtPriceFloat)
}

// ScalePrice scales the price to the desired ticker decimals. The price is normalized to
// the token decimals in the erc20 token contracts.
func ScalePrice(
	cfg PoolConfig,
	price *big.Float,
) *big.Float {
	// Adjust the price based on the difference between the token decimals in the erc20 token contracts.
	erc20ScalingFactor := math.GetScalingFactor(
		cfg.BaseDecimals,
		cfg.QuoteDecimals,
	)

	// Invert the price if the configuration specifies to do so.
	if cfg.Invert {
		scaledERC20AdjustedPrice := new(big.Float).Quo(price, erc20ScalingFactor)
		return new(big.Float).Quo(big.NewFloat(1), scaledERC20AdjustedPrice)
	}
	return new(big.Float).Mul(price, erc20ScalingFactor)
}
