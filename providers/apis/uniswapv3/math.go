package uniswapv3

import (
	"math/big"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// ConvertSquareRootX96Price converts the slot 0 sqrtPriceX96 value to a price. Note that this
// price is not scaled to the token decimals. This calculation is equivalent to:
//
// price = (sqrtPriceX96 / 2^96) ^ 2
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

// ScalePrice scales the price to the respective token decimals. There are few different steps
// that are taken to scale the price:
//
//  1. Scale the price based on the difference between the token decimals in the erc20 token contracts.
//     If the difference is positive we are truncating the price, if the difference is negative we are
//     adding zeros to the price.
//  2. Adjust the price based on the difference between the desired precision and the base token decimals.
//
// This entire calculation is equvalent to the following:
//
// price = ((price / 10^(baseDecimals - quoteDecimals)) * 10^(tickerDecimals)) ^ -1
// where the ^ -1 is only applied if the configuration specifies to invert the price.
func ScalePrice(
	ticker mmtypes.Ticker,
	cfg PoolConfig,
	price *big.Float,
) *big.Float {
	// Adjust the price based on the difference between the token decimals in the erc20 token contracts.
	erc20ScalingFactor := GetScalingFactor(
		cfg.BaseDecimals,
		cfg.QuoteDecimals,
	)

	// Invert the price if the configuration specifies to do so.
	var scaledERC20AdjustedPrice *big.Float
	if cfg.Invert {
		scaledERC20AdjustedPrice = new(big.Float).Quo(price, erc20ScalingFactor)
		scaledERC20AdjustedPrice = new(big.Float).Quo(big.NewFloat(1), scaledERC20AdjustedPrice)
	} else {
		scaledERC20AdjustedPrice = new(big.Float).Mul(price, erc20ScalingFactor)
	}

	one := new(big.Float).SetInt(oracle.ScaledOne(ticker.Decimals))
	return new(big.Float).Mul(scaledERC20AdjustedPrice, one)
}

// GetScalingFactor returns the scaling factor for the price based on the difference between
// the token decimals in the erc20 token contracts. Please read over the Uniswap V3 math primer
// for more information on how this is utilized.
func GetScalingFactor(
	first, second int64,
) *big.Float {
	// Determine the scaling factor for the price.
	decimalDiff := first - second
	exp := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(10), big.NewInt(math.Abs(decimalDiff)), nil),
	)

	if decimalDiff > 0 {
		return exp
	}
	return new(big.Float).Quo(big.NewFloat(1), exp)
}
