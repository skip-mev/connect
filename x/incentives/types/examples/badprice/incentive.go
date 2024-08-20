package badprice

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

const (
	// BadPriceIncentiveType is the type for the BadPriceIncentive.
	BadPriceIncentiveType = "bad_price_update"
)

// Each Incentive type must implement the types.Incentive interface.
var _ types.Incentive = (*BadPriceIncentive)(nil)

// NewBadPriceIncentive creates a new BadPriceIncentive.
//
// NOTE: THIS SHOULD NOT BE USED IN PRODUCTION. THIS IS ONLY FOR TESTING.
func NewBadPriceIncentive(validator sdk.ValAddress, amount math.Int) *BadPriceIncentive {
	return &BadPriceIncentive{
		Validator: validator.String(),
		Amount:    amount.String(),
	}
}

// ValidateBasic does a basic stateless validation check that
// doesn't require access to any other information.
func (b *BadPriceIncentive) ValidateBasic() error {
	// You can add your custom validation logic here if needed.
	_, err := sdk.ValAddressFromBech32(b.Validator)
	if err != nil {
		return fmt.Errorf("invalid validator address %s: %w", b.Validator, err)
	}

	amount, ok := math.NewIntFromString(b.Amount)
	if !ok {
		return fmt.Errorf("invalid amount %s", b.Amount)
	}

	if amount.IsNegative() {
		return fmt.Errorf("amount %s cannot be negative", b.Amount)
	}

	return nil
}

// Type returns the type of the incentive.
func (b *BadPriceIncentive) Type() string {
	return BadPriceIncentiveType
}

// Copy returns a copy of the incentive.
func (b *BadPriceIncentive) Copy() types.Incentive {
	return &BadPriceIncentive{
		Validator: b.Validator,
		Amount:    b.Amount,
	}
}
