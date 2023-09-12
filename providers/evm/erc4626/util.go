package erc4626

import "math/big"

func getUnitValueFromDecimals(decimals uint64) *big.Int {
	return big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
}
