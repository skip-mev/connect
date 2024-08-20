package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/badprice"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/goodprice"
)

func (s *KeeperTestSuite) TestGetIncentivesByType() {
	s.Run("returns an error with empty type", func() {
		_, err := s.queryServer.GetIncentivesByType(s.ctx, nil)
		s.Require().Error(err)
	})

	s.Run("returns an error with unsupported type", func() {
		req := &types.GetIncentivesByTypeRequest{IncentiveType: "unsupported"}
		_, err := s.queryServer.GetIncentivesByType(s.ctx, req)
		s.Require().Error(err)
	})

	s.Run("returns an empty list with no incentives", func() {
		req := &types.GetIncentivesByTypeRequest{IncentiveType: badprice.BadPriceIncentiveType}
		resp, err := s.queryServer.GetIncentivesByType(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Entries, 0)
	})

	s.Run("returns a single incentive stored in the module", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		incentives := []types.Incentive{badPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		req := &types.GetIncentivesByTypeRequest{IncentiveType: badprice.BadPriceIncentiveType}
		resp, err := s.queryServer.GetIncentivesByType(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Entries, 1)

		// check that the incentive is the same as the one we added
		bz, err := badPrice.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz, resp.Entries[0])
	})

	s.Run("returns multiple incentives stored in the module", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))
		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)

		incentives := []types.Incentive{badPrice1, badPrice2}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		req := &types.GetIncentivesByTypeRequest{IncentiveType: badprice.BadPriceIncentiveType}
		resp, err := s.queryServer.GetIncentivesByType(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Entries, 2)

		// check that the incentives are the same as the ones we added
		bz1, err := badPrice1.Marshal()
		s.Require().NoError(err)

		bz2, err := badPrice2.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz1, resp.Entries[0])
		s.Require().Equal(bz2, resp.Entries[1])
	})
}

func (s *KeeperTestSuite) TestGetAllIncentives() {
	s.Run("returns an empty list with no incentives", func() {
		req := &types.GetAllIncentivesRequest{}
		resp, err := s.queryServer.GetAllIncentives(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Registry, 0)
	})

	s.Run("returns a single incentive stored in the module", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		incentives := []types.Incentive{badPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		req := &types.GetAllIncentivesRequest{}
		resp, err := s.queryServer.GetAllIncentives(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Registry, 1)

		s.Require().Equal(badprice.BadPriceIncentiveType, resp.Registry[0].IncentiveType)
		s.Require().Len(resp.Registry[0].Entries, 1)

		// check that the incentive is the same as the one we added
		bz, err := badPrice.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz, resp.Registry[0].Entries[0])
	})

	s.Run("returns multiple incentives stored in the module", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))
		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)

		incentives := []types.Incentive{badPrice1, badPrice2}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		req := &types.GetAllIncentivesRequest{}
		resp, err := s.queryServer.GetAllIncentives(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Registry, 1)

		s.Require().Equal(badprice.BadPriceIncentiveType, resp.Registry[0].IncentiveType)
		s.Require().Len(resp.Registry[0].Entries, 2)

		// check that the incentives are the same as the ones we added
		bz1, err := badPrice1.Marshal()
		s.Require().NoError(err)

		bz2, err := badPrice2.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz1, resp.Registry[0].Entries[0])
		s.Require().Equal(bz2, resp.Registry[0].Entries[1])
	})

	s.Run("returns a single incentive for each type", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))
		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice := badprice.NewBadPriceIncentive(validator1, amount1)
		goodPrice := goodprice.NewGoodPriceIncentive(validator2, amount2)

		incentives := []types.Incentive{badPrice, goodPrice}

		err := s.incentivesKeeper.AddIncentives(s.ctx, incentives)
		s.Require().NoError(err)

		req := &types.GetAllIncentivesRequest{}
		resp, err := s.queryServer.GetAllIncentives(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(resp.Registry, 2)

		s.Require().Equal(badprice.BadPriceIncentiveType, resp.Registry[0].IncentiveType)
		s.Require().Len(resp.Registry[0].Entries, 1)

		s.Require().Equal(goodprice.GoodPriceIncentiveType, resp.Registry[1].IncentiveType)
		s.Require().Len(resp.Registry[1].Entries, 1)

		// check that the incentives are the same as the ones we added
		bz1, err := badPrice.Marshal()
		s.Require().NoError(err)

		bz2, err := goodPrice.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz1, resp.Registry[0].Entries[0])
		s.Require().Equal(bz2, resp.Registry[1].Entries[0])
	})
}
