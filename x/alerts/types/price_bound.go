package types

import (
	"fmt"
	"math/big"
)

// ValidateBasic validates the PriceBound, specifically that the high price-bound is greater than the low price-bound, and
// that the high / low price-bounds are valid big.Int values.
func (pb *PriceBound) ValidateBasic() error {
	high, converted := new(big.Int).SetString(pb.High, 10)
	if !converted {
		return fmt.Errorf("invalid high price-bound %s", pb.High)
	}

	low, converted := new(big.Int).SetString(pb.Low, 10)
	if !converted {
		return fmt.Errorf("invalid low price-bound %s", pb.Low)
	}

	if high.Cmp(low) <= 0 {
		return fmt.Errorf("high price-bound %s must be greater than low price-bound %s", pb.High, pb.Low)
	}

	return nil
}

// GetHighInt returns the high price-bound as a big.Int.
func (pb *PriceBound) GetHighInt() (*big.Int, error) {
	high, converted := new(big.Int).SetString(pb.High, 10)
	if !converted {
		return nil, fmt.Errorf("invalid high price-bound %s", pb.High)
	}

	return high, nil
}

// GetLowInt returns the low price-bound as a big.Int.
func (pb *PriceBound) GetLowInt() (*big.Int, error) {
	low, converted := new(big.Int).SetString(pb.Low, 10)
	if !converted {
		return nil, fmt.Errorf("invalid low price-bound %s", pb.Low)
	}

	return low, nil
}
