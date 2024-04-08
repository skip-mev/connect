package math

import (
	"fmt"
	"math/big"
)

// Min returns the minimum of two values.
func Min[V int | int64 | uint64 | int32 | uint32](vals ...V) V {
	if len(vals) == 0 {
		panic("cannot find minimum of empty slice")
	}

	minimum := vals[0]
	for _, val := range vals[1:] {
		if val < minimum {
			minimum = val
		}
	}
	return minimum
}

// Float64StringToBigInt converts a float64 string to a big.Int.
func Float64StringToBigInt(s string, decimals uint64) (*big.Int, error) {
	bigFloat := new(big.Float)
	_, _, err := bigFloat.Parse(s, 10)
	if err != nil {
		return nil, err
	}

	return BigFloatToBigInt(bigFloat, decimals), nil
}

// Float64ToBigInt converts a float64 to a big.Int.
func Float64ToBigInt(val float64, decimals uint64) *big.Int {
	bigVal := new(big.Float)
	bigVal.SetFloat64(val)

	return BigFloatToBigInt(bigVal, decimals)
}

// BigFloatToBigInt converts a big.Float to a big.Int.
func BigFloatToBigInt(f *big.Float, decimals uint64) *big.Int {
	bigFloat := new(big.Float)
	factor := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	bigFloat.SetInt(factor)

	f.Mul(f, bigFloat)

	result := new(big.Int)
	f.Int(result) // store converted number in result

	return result
}

// Float64StringToBigFloat converts a float64 string to a big.Float.
func Float64StringToBigFloat(s string) (*big.Float, error) {
	bigFloat := new(big.Float)
	_, ok := bigFloat.SetString(s)
	if !ok {
		return nil, fmt.Errorf("failed to set big.Float from string: %s", s)
	}
	return bigFloat, nil
}

// ScaleBigFloat scales a big.Float by the given decimals.
func ScaleBigFloat(f *big.Float, decimals uint64) *big.Float {
	bigFloat := new(big.Float)
	factor := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	bigFloat.SetInt(factor)

	f.Mul(f, bigFloat)

	return f
}
