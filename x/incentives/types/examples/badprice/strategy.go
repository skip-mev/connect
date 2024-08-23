package badprice

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// StakingKeeper defines the expected staking keeper interface required by
// the BadPriceIncentive strategy.
//
//go:generate mockery --name StakingKeeper --output ./../mocks --outpkg mocks --case underscore
type StakingKeeper interface {
	// GetValidatorStake returns the total amount of stake that a validator
	// currently has delegated to it.
	GetValidatorStake(ctx sdk.Context, val sdk.ValAddress) (stake math.Int, found bool)

	// Slash attempts to slash the validator at the given address by the
	// given amount. If the validator does not have sufficient stake to
	// slash, the slash will fail.
	Slash(ctx sdk.Context, val sdk.ValAddress, amount math.Int) error
}

// BadPriceIncentiveStrategy is the strategy function for the BadPriceIncentive
// type.
type BadPriceIncentiveStrategy struct {
	keeper StakingKeeper
}

// NewBadPriceIncentiveStrategy returns a new BadPriceIncentiveStrategy.
func NewBadPriceIncentiveStrategy(keeper StakingKeeper) *BadPriceIncentiveStrategy {
	return &BadPriceIncentiveStrategy{
		keeper: keeper,
	}
}

// GetStrategy returns the BadPriceIncentiveStrategy.
func (s *BadPriceIncentiveStrategy) GetStrategy() types.Strategy {
	return func(ctx sdk.Context, incentive types.Incentive) (types.Incentive, error) {
		// Cast the incentive to the concrete type.
		badPriceIncentive, ok := incentive.(*BadPriceIncentive)
		if !ok {
			return incentive, fmt.Errorf("invalid incentive type: %T", incentive)
		}

		validator, err := sdk.ValAddressFromBech32(badPriceIncentive.Validator)
		if err != nil {
			return incentive, nil
		}

		// Get the validator's current stake.
		stake, found := s.keeper.GetValidatorStake(ctx, validator)
		if !found {
			return incentive, nil
		}

		// Check the upper bound on what we can slash
		slashAmount, ok := math.NewIntFromString(badPriceIncentive.Amount)
		if !ok {
			return incentive, nil
		}

		if slashAmount.GT(stake) {
			slashAmount = stake
		}

		if err := s.keeper.Slash(ctx, validator, slashAmount); err != nil {
			return incentive, err
		}

		return nil, nil
	}
}
