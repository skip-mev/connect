package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"

	"github.com/skip-mev/connect/v2/x/incentives/keeper"
	"github.com/skip-mev/connect/v2/x/incentives/types"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/badprice"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/goodprice"
)

func (s *KeeperTestSuite) TestExecuteByIncentiveType() {
	deleteCB := func(_ types.Incentive) (types.Incentive, error) {
		return nil, nil
	}

	updatePriceCB := func(incentive types.Incentive) (types.Incentive, error) {
		badPrice, ok := incentive.(*badprice.BadPriceIncentive)
		s.Require().True(ok)

		badPrice.Amount = math.NewInt(200000).String()
		return badPrice, nil
	}

	s.Run("can update an empty list of incentives", func() {
		err := s.incentivesKeeper.ExecuteByIncentiveType(s.ctx, &badprice.BadPriceIncentive{}, deleteCB)
		s.Require().NoError(err)

		// Check the count of incentives for each type.
		incentives, err := s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)
	})

	s.Run("can run a no-op on a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		incentives := []types.Incentive{badPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Update the incentives with the no-op callback.
		err = s.incentivesKeeper.ExecuteByIncentiveType(s.ctx, &badprice.BadPriceIncentive{}, deleteCB)
		s.Require().NoError(err)

		// Check that the incentive was removed from the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)
	})

	s.Run("can run a valid update on a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		incentives := []types.Incentive{badPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Check the incentive in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator.String(), i.Validator)
		s.Require().Equal(amount.String(), i.Amount)

		// Update the incentives with the update callback.
		err = s.incentivesKeeper.ExecuteByIncentiveType(s.ctx, &badprice.BadPriceIncentive{}, updatePriceCB)
		s.Require().NoError(err)

		// Check that the incentive was updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i, ok = incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator.String(), i.Validator)
		s.Require().Equal(math.NewInt(200000).String(), i.Amount)
	})

	s.Run("can update multiple incentives", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)

		incentives := []types.Incentive{badPrice1, badPrice2}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Check the incentives in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		i1, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(amount1.String(), i1.Amount)

		i2, ok := incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(amount2.String(), i2.Amount)

		// Update the incentives with the update callback.
		err = s.incentivesKeeper.ExecuteByIncentiveType(s.ctx, &badprice.BadPriceIncentive{}, updatePriceCB)
		s.Require().NoError(err)

		// Check that the incentives were updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		i1, ok = incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(math.NewInt(200000).String(), i1.Amount)

		i2, ok = incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(math.NewInt(200000).String(), i2.Amount)
	})

	s.Run("can update some incentives and remove others", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))
		validator3 := sdk.ValAddress([]byte("validator3"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)
		amount3 := math.NewInt(300)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)
		badPrice3 := badprice.NewBadPriceIncentive(validator3, amount3)

		incentives := []types.Incentive{badPrice1, badPrice2, badPrice3}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Check the incentives in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 3)

		i1, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(amount1.String(), i1.Amount)

		i2, ok := incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(amount2.String(), i2.Amount)

		i3, ok := incentives[2].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator3.String(), i3.Validator)
		s.Require().Equal(amount3.String(), i3.Amount)

		cb := func(incentive types.Incentive) (types.Incentive, error) {
			badPrice, ok := incentive.(*badprice.BadPriceIncentive)
			s.Require().True(ok)

			// If this is validator 2 we remove
			if badPrice.Validator == validator2.String() {
				return nil, nil
			}

			// Otherwise we update the price
			badPrice.Amount = math.NewInt(200000).String()
			return badPrice, nil
		}

		// Update the incentives with the update callback.
		err = s.incentivesKeeper.ExecuteByIncentiveType(s.ctx, &badprice.BadPriceIncentive{}, cb)
		s.Require().NoError(err)

		// Check that the incentives were updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		i1, ok = incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(math.NewInt(200000).String(), i1.Amount)

		i2, ok = incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator3.String(), i2.Validator)
		s.Require().Equal(math.NewInt(200000).String(), i2.Amount)
	})
}

func (s *KeeperTestSuite) TestExecuteStrategies() {
	s.Run("can execute a strategy on an empty list of incentives", func() {
		err := s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)
	})

	s.Run("can execute a strategy on a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		// Add the incentive to the store.
		incentives := []types.Incentive{badPrice}
		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Mock the results of the staking keeper.
		s.stakingKeeper.On("GetValidatorStake", mock.Anything, validator).Return(amount, true).Once()
		s.stakingKeeper.On("Slash", mock.Anything, validator, amount).Return(nil).Once()

		// Execute the strategy.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentive was removed from the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)
	})

	s.Run("stores the incentive if the strategy returns an error", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		// Add the incentive to the store.
		incentives := []types.Incentive{badPrice}
		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Mock the results of the staking keeper.
		s.stakingKeeper.On("GetValidatorStake", mock.Anything, validator).Return(amount, true).Once()
		s.stakingKeeper.On("Slash", mock.Anything, validator, amount).Return(fmt.Errorf("slash error")).Once()

		// Execute the strategy.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().Error(err)

		// Check that the incentive was not removed from the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)
	})

	s.Run("can execute a strategy on multiple of the same incentive types", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)
		incentives := []types.Incentive{badPrice1, badPrice2}

		// Add the incentives to the store.
		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Mock the results of the staking keeper.
		s.stakingKeeper.On("GetValidatorStake", mock.Anything, validator1).Return(amount1, true).Once()
		s.stakingKeeper.On("Slash", mock.Anything, validator1, amount1).Return(nil).Once()

		s.stakingKeeper.On("GetValidatorStake", mock.Anything, validator2).Return(amount2, true).Once()
		s.stakingKeeper.On("Slash", mock.Anything, validator2, amount2).Return(nil).Once()

		// Execute the strategy.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentives were removed from the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)
	})

	s.Run("can execute a strategy on multiple different incentive types", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice := badprice.NewBadPriceIncentive(validator1, amount1)
		goodPrice := goodprice.NewGoodPriceIncentive(validator2, amount2)
		incentives := []types.Incentive{badPrice, goodPrice}

		// Add the incentives to the store.
		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Mock the results of the staking keeper.
		s.stakingKeeper.On("GetValidatorStake", mock.Anything, validator1).Return(amount1, true).Once()
		s.stakingKeeper.On("Slash", mock.Anything, validator1, amount1).Return(nil).Once()

		// Mock the results of the bank keeper.
		s.bankKeeper.On(
			"MintCoins",
			mock.Anything, mock.Anything, amount2,
		).Return(nil).Once()

		s.bankKeeper.On(
			"SendCoinsFromModuleToAccount",
			mock.Anything, mock.Anything, validator2, amount2,
		).Return(nil).Once()

		// Execute the strategy.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentives were removed from the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)

		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)
	})

	updateBadPriceStrategy := func(_ sdk.Context, incentive types.Incentive) (types.Incentive, error) {
		badPrice, ok := incentive.(*badprice.BadPriceIncentive)
		s.Require().True(ok)

		amount, ok := math.NewIntFromString(badPrice.Amount)
		s.Require().True(ok)

		amount = amount.Add(math.NewInt(100))
		badPrice.Amount = amount.String()

		return badPrice, nil
	}

	s.Run("can run an update strategy on a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		// create a new keeper with the updated strategy
		strategies := map[types.Incentive]types.Strategy{
			&badprice.BadPriceIncentive{}: updateBadPriceStrategy,
		}
		s.incentivesKeeper = keeper.NewKeeper(s.key, strategies)

		// Add the incentive to the store.
		incentives := []types.Incentive{badPrice}
		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Execute the strategy.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentive was updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator.String(), i.Validator)
		s.Require().Equal(math.NewInt(200).String(), i.Amount)

		// Execute the strategy again.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentive was updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i, ok = incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator.String(), i.Validator)
		s.Require().Equal(math.NewInt(300).String(), i.Amount)
	})

	updateGoodPriceStrategy := func(_ sdk.Context, incentive types.Incentive) (types.Incentive, error) {
		goodPrice, ok := incentive.(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)

		amount, ok := math.NewIntFromString(goodPrice.Amount)
		s.Require().True(ok)

		amount = amount.Add(math.NewInt(100))
		goodPrice.Amount = amount.String()

		return goodPrice, nil
	}

	s.Run("can run an update strategy on multiple different incentive types", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice := badprice.NewBadPriceIncentive(validator1, amount1)
		goodPrice := goodprice.NewGoodPriceIncentive(validator2, amount2)

		// create a new keeper with the updated strategy
		strategies := map[types.Incentive]types.Strategy{
			&badprice.BadPriceIncentive{}:   updateBadPriceStrategy,
			&goodprice.GoodPriceIncentive{}: updateGoodPriceStrategy,
		}
		s.incentivesKeeper = keeper.NewKeeper(s.key, strategies)

		// Add the incentives to the store.
		incentives := []types.Incentive{badPrice, goodPrice}
		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// Execute the strategy.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentives were updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i.Validator)
		s.Require().Equal(math.NewInt(200).String(), i.Amount)

		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i2, ok := incentives[0].(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(math.NewInt(300).String(), i2.Amount)

		// Execute the strategy again.
		err = s.incentivesKeeper.ExecuteStrategies(s.ctx)
		s.Require().NoError(err)

		// Check that the incentives were updated in the store.
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i, ok = incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i.Validator)
		s.Require().Equal(math.NewInt(300).String(), i.Amount)

		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i2, ok = incentives[0].(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(math.NewInt(400).String(), i2.Amount)
	})
}
