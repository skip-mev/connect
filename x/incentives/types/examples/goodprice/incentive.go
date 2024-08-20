package goodprice

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

const (
	// GoodPriceIncentiveType is the type for the GoodPriceIncentive.
	GoodPriceIncentiveType = "good_price_update"
)

// Ensure that GoodPriceIncentive implements the Incentive interface
var _ types.Incentive = (*GoodPriceIncentive)(nil)

// NewGoodPriceIncentive returns a new GoodPriceIncentive.
//
// NOTE: THIS SHOULD NOT BE USED IN PRODUCTION. THIS IS ONLY FOR TESTING.
func NewGoodPriceIncentive(validator sdk.ValAddress, amount math.Int) *GoodPriceIncentive {
	return &GoodPriceIncentive{
		Validator: validator.String(),
		Amount:    amount.String(),
	}
}

// ValidateBasic does a basic stateless validation check that
// doesn't require access to any other information.
func (b *GoodPriceIncentive) ValidateBasic() error {
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
func (b *GoodPriceIncentive) Type() string {
	return GoodPriceIncentiveType
}

// Copy returns a copy of the incentive.
func (b *GoodPriceIncentive) Copy() types.Incentive {
	return &GoodPriceIncentive{
		Validator: b.Validator,
		Amount:    b.Amount,
	}
}
