package goodprice

import (
	"context"
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// BankKeeper defines the expected bank keeper interface required by
// the GoodPriceIncentive strategy.
//
//go:generate mockery --name BankKeeper --output ./../mocks --outpkg mocks --case underscore
type BankKeeper interface {
	// Mint coins and add to the module account before distributing them
	// to the recipient.
	MintCoins(
		ctx context.Context,
		moduleName string,
		amt math.Int,
	) error

	// Send coins from the module account to the recipient.
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.ValAddress,
		amt math.Int,
	) error
}

// GoodPriceIncentiveStrategy is the strategy function for the GoodPriceIncentive
// type.
type GoodPriceIncentiveStrategy struct {
	bk BankKeeper
}

// NewBadPriceIncentiveStrategy returns a new GoodPriceIncentiveStrategy.
func NewGoodPriceIncentiveStrategy(bk BankKeeper) *GoodPriceIncentiveStrategy {
	return &GoodPriceIncentiveStrategy{
		bk: bk,
	}
}

// GetStrategy returns the BadPriceIncentiveStrategy.
func (s *GoodPriceIncentiveStrategy) GetStrategy() types.Strategy {
	return func(ctx sdk.Context, incentive types.Incentive) (types.Incentive, error) {
		// Cast the incentive to the concrete type.
		badPriceIncentive, ok := incentive.(*GoodPriceIncentive)
		if !ok {
			return incentive, fmt.Errorf("invalid incentive type: %T", incentive)
		}

		validator, err := sdk.ValAddressFromBech32(badPriceIncentive.Validator)
		if err != nil {
			return nil, nil
		}

		amount, ok := math.NewIntFromString(badPriceIncentive.Amount)
		if !ok {
			return nil, nil
		}

		if amount.IsNegative() {
			return nil, nil
		}

		// Mint coins and add to the module account before distributing them
		// to the recipient.
		if err := s.bk.MintCoins(ctx, types.ModuleName, amount); err != nil {
			return incentive, err
		}

		// Send coins from the module account to the recipient.
		if err := s.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, validator, amount); err != nil {
			return incentive, err
		}

		return nil, nil
	}
}
