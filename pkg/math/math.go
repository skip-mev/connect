package math

import (
	"math/big"
	"strconv"
)

// Min returns the minimum of two values.
func Min[V int | int64 | uint64 | int32 | uint32](vals ...V) V {
	if len(vals) == 0 {
		panic("cannot find minimum of empty slice")
	}

	min := vals[0]
	for _, val := range vals[1:] {
		if val < min {
			min = val
		}
	}
	return min
}

// Float64StringToBigInt converts a float64 string to a big.Int.
func Float64StringToBigInt(s string, decimals int64) (*big.Int, error) {
	floatNum, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}

	return Float64ToBigInt(floatNum, decimals), nil
}

// Float64ToBigInt converts a float64 to a big.Int.
//
// TODO: Is there a better approach to this?
func Float64ToBigInt(val float64, decimals int64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	factor := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	coin.SetInt(factor)

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}
