package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func (s *KeeperTestSuite) TestGetAllSLAsRequest() {
	sla1 := slatypes.NewPriceFeedSLA(
		"id",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	sla2 := slatypes.NewPriceFeedSLA(
		"id2",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	sla3 := slatypes.NewPriceFeedSLA(
		"id3",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	s.Run("can get empty set of SLAs", func() {
		req := &slatypes.GetAllSLAsRequest{}
		resp, err := s.queryServer.GetAllSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().Empty(resp.SLAs)
	})

	s.Run("can get all SLAs", func() {
		s.Require().NoError(s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla1, sla2, sla3}))

		req := &slatypes.GetAllSLAsRequest{}
		resp, err := s.queryServer.GetAllSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Len(resp.SLAs, 3)
		s.Require().Equal(sla1, resp.SLAs[0])
		s.Require().Equal(sla2, resp.SLAs[1])
		s.Require().Equal(sla3, resp.SLAs[2])

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 3)

		s.Require().Equal(resp.SLAs[0], slas[0])
		s.Require().Equal(resp.SLAs[1], slas[1])
		s.Require().Equal(resp.SLAs[2], slas[2])
	})
}

func (s *KeeperTestSuite) TestGetPriceFeedsRequest() {
	cp1 := slinkytypes.NewCurrencyPair("btc", "usd")

	consAddress1 := sdk.ConsAddress("consAddress1")
	consAddress2 := sdk.ConsAddress("consAddress2")

	id1 := "testId"
	priceFeed1, err := slatypes.NewPriceFeed(
		10,
		consAddress1,
		cp1,
		id1,
	)
	s.Require().NoError(err)
	priceFeed2, _ := slatypes.NewPriceFeed(
		10,
		consAddress2,
		cp1,
		id1,
	)
	s.Require().NoError(err)

	s.Run("can get empty set of price feeds", func() {
		req := &slatypes.GetPriceFeedsRequest{}
		resp, err := s.queryServer.GetPriceFeeds(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().Empty(resp.PriceFeeds)
	})

	s.Run("can get all price feeds", func() {
		err := s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		req := &slatypes.GetPriceFeedsRequest{
			ID: id1,
		}
		resp, err := s.queryServer.GetPriceFeeds(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Len(resp.PriceFeeds, 2)
		s.Require().Equal(priceFeed1, resp.PriceFeeds[0])
		s.Require().Equal(priceFeed2, resp.PriceFeeds[1])

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 2)

		s.Require().Equal(resp.PriceFeeds[0], feeds[0])
		s.Require().Equal(resp.PriceFeeds[1], feeds[1])
	})

	s.Run("multiple different feed IDs are set and returns only ones that correspond to it", func() {
		feed1, err := slatypes.NewPriceFeed(
			10,
			consAddress1,
			cp1,
			id1,
		)
		s.Require().NoError(err)

		feed2, err := slatypes.NewPriceFeed(
			10,
			consAddress1,
			cp1,
			"id2",
		)
		s.Require().NoError(err)

		feed3, err := slatypes.NewPriceFeed(
			10,
			consAddress1,
			cp1,
			"id3",
		)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed2)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, feed3)
		s.Require().NoError(err)

		req := &slatypes.GetPriceFeedsRequest{
			ID: "id2",
		}
		resp, err := s.queryServer.GetPriceFeeds(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Len(resp.PriceFeeds, 1)
		s.Require().Equal(feed2, resp.PriceFeeds[0])

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, "id2")
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)
		s.Require().Equal(resp.PriceFeeds[0], feeds[0])

		req = &slatypes.GetPriceFeedsRequest{
			ID: id1,
		}
		resp, err = s.queryServer.GetPriceFeeds(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().Len(resp.PriceFeeds, 1)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, id1)
		s.Require().NoError(err)
		s.Require().Len(feeds, 1)
		s.Require().Equal(resp.PriceFeeds[0], feeds[0])
	})
}

func (s *KeeperTestSuite) TestGetParamsRequest() {
	s.Run("can get default params", func() {
		req := &slatypes.ParamsRequest{}
		resp, err := s.queryServer.Params(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Equal(slatypes.DefaultParams(), resp.Params)

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(resp.Params, params)
	})

	s.Run("can get updated params", func() {
		params := slatypes.NewParams(false)
		err := s.keeper.SetParams(s.ctx, params)
		s.Require().NoError(err)

		req := &slatypes.ParamsRequest{}
		resp, err := s.queryServer.Params(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Equal(params, resp.Params)

		params, err = s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(resp.Params, params)
	})
}
