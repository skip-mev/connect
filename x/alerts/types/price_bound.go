package types

import (
	"fmt"

	"github.com/holiman/uint256"
)

// ValidateBasic validates the PriceBound, specifically that the high price-bound is greater than the low price-bound, and
// that the high / low price-bounds are valid uint256 values.
func (pb PriceBound) ValidateBasic() error {
	high, err := uint256.FromHex(pb.High)
	if err != nil {
		return err
	}

	low, err := uint256.FromHex(pb.Low)
	if err != nil {
		return err
	}

	if high.Cmp(low) <= 0 {
		return fmt.Errorf("high price-bound %s must be greater than low price-bound %s", pb.High, pb.Low)
	}

	return nil
}

// GetHighInt returns the high price-bound as a uint256.Int.
func (pb PriceBound) GetHighInt() (*uint256.Int, error) {
	return uint256.FromHex(pb.High)
}

// GetLowInt returns the low price-bound as a uint256.Int.
func (pb PriceBound) GetLowInt() (*uint256.Int, error) {
	return uint256.FromHex(pb.Low)
}
