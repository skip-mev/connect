package providers

import (
	"math/big"
	"strconv"

	"github.com/holiman/uint256"
)

// Float64StringToUint256 converts a float64 string to a uint256.
func Float64StringToUint256(s string, decimals int) (*uint256.Int, error) {
	floatNum, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}

	return Float64ToUint256(floatNum, decimals), nil
}

// Float64ToBigInt converts a float64 to a uint256.
//
// NOTE: MustFromBig will panic only if there is overflow when
// converting the big.Int to a uint256.Int. This should never
// happen since uint256 should be large enough to handle pricing data.
func Float64ToUint256(val float64, decimals int) *uint256.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	factor := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	coin.SetInt(factor)

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return uint256.MustFromBig(result)
}
