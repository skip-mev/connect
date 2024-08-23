package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/badprice"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/goodprice"
)

func (s *KeeperTestSuite) TestAddIncentives() {
	s.Run("can add an empty list of incentives", func() {
		err := s.incentivesKeeper.AddIncentives(s.ctx, nil)
		s.Require().NoError(err)

		// Check the count of incentives for each type.
		incentives, err := s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 0)
	})

	s.Run("can add a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		incentives := []types.Incentive{badPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// retrieve the incentive from the store
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		// check that the incentive is the same as the one we added
		i, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator.String(), i.Validator)
		s.Require().Equal(amount.String(), i.Amount)
	})

	s.Run("can add multiple incentives", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))
		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)

		incentives := []types.Incentive{badPrice1, badPrice2}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// retrieve the incentives from the store
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		// check that the incentives are the same as the ones we added
		i1, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(amount1.String(), i1.Amount)

		i2, ok := incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(amount2.String(), i2.Amount)
	})

	s.Run("can add single incentive of different types", func() {
		goodValidator := sdk.ValAddress([]byte("good_validator"))
		goodAmount := math.NewInt(100)
		goodPrice := goodprice.NewGoodPriceIncentive(goodValidator, goodAmount)

		badValidator := sdk.ValAddress([]byte("bad_validator"))
		badAmount := math.NewInt(200)
		badPrice := badprice.NewBadPriceIncentive(badValidator, badAmount)

		incentives := []types.Incentive{goodPrice, badPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// retrieve the incentives from the store
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i1, ok := incentives[0].(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(goodValidator.String(), i1.Validator)
		s.Require().Equal(goodAmount.String(), i1.Amount)

		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i2, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(badValidator.String(), i2.Validator)
		s.Require().Equal(badAmount.String(), i2.Amount)
	})

	s.Run("can add multiple incentives of different types", func() {
		goodValidator1 := sdk.ValAddress([]byte("good_validator1"))
		goodAmount1 := math.NewInt(100)
		goodPrice1 := goodprice.NewGoodPriceIncentive(goodValidator1, goodAmount1)

		goodValidator2 := sdk.ValAddress([]byte("good_validator2"))
		goodAmount2 := math.NewInt(200)
		goodPrice2 := goodprice.NewGoodPriceIncentive(goodValidator2, goodAmount2)

		badValidator1 := sdk.ValAddress([]byte("bad_validator1"))
		badAmount1 := math.NewInt(300)
		badPrice1 := badprice.NewBadPriceIncentive(badValidator1, badAmount1)

		badValidator2 := sdk.ValAddress([]byte("bad_validator2"))
		badAmount2 := math.NewInt(400)
		badPrice2 := badprice.NewBadPriceIncentive(badValidator2, badAmount2)

		incentives := []types.Incentive{goodPrice1, goodPrice2, badPrice1, badPrice2}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		// retrieve the incentives from the store
		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		i1, ok := incentives[0].(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(goodValidator1.String(), i1.Validator)
		s.Require().Equal(goodAmount1.String(), i1.Amount)

		i2, ok := incentives[1].(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(goodValidator2.String(), i2.Validator)
		s.Require().Equal(goodAmount2.String(), i2.Amount)

		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		i3, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(badValidator1.String(), i3.Validator)
		s.Require().Equal(badAmount1.String(), i3.Amount)

		i4, ok := incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(badValidator2.String(), i4.Validator)
		s.Require().Equal(badAmount2.String(), i4.Amount)
	})
}
