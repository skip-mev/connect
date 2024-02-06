package oracle

import (
	"math/big"
)

// ScaledDecimals is the standard number of decimal places each price will be converted to
// during the conversion process.
const ScaledDecimals = 18

func ScaleUpCurrencyPairPrice(decimals int64, price *big.Int) *big.Int {
	diff := int64(ScaledDecimals - decimals)
	exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(diff), nil)
	return new(big.Int).Mul(price, exp)
}

func ScaleDownCurrencyPair(decimals int64, price *big.Int) *big.Int {
	diff := int64(ScaledDecimals - decimals)
	exp := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(diff), nil)
	return new(big.Int).Div(price, exp)
}
