package math

import (
	"fmt"
	"math/big"
	"sort"

	"golang.org/x/exp/constraints"
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

// Abs returns the absolute value of a given number.
func Abs[V constraints.Signed](val V) V {
	if val < 0 {
		return -val
	}
	return val
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

// SortBigFloats is a stable slices sort for an array of big.Floats.
func SortBigFloats(values []*big.Float) {
	// Sort the values.
	sort.SliceStable(values, func(i, j int) bool {
		switch values[i].Cmp(values[j]) {
		case -1:
			return true
		case 1:
			return false
		default:
			return true
		}
	})
}

// CalculateMedian calculates the median from a list of big.Float. Returns an
// average if the number of values is even.
func CalculateMedian(values []*big.Float) *big.Float {
	if len(values) == 0 {
		return nil
	}
	SortBigFloats(values)

	middleIndex := len(values) / 2

	// Calculate the median.
	numValues := len(values)
	var median *big.Float
	if numValues%2 == 0 { // even
		median = new(big.Float).Add(values[middleIndex-1], values[middleIndex])
		median = median.Quo(median, new(big.Float).SetUint64(2))
	} else { // odd
		median = values[middleIndex]
	}

	return median
}
