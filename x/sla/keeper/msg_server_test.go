package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func (s *KeeperTestSuite) TestMsgAddSLAs() {
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

	s.Run("accepts a req with valid SLA", func() {
		req := &slatypes.MsgAddSLAs{
			Authority: s.authority.String(),
			SLAs:      []slatypes.PriceFeedSLA{sla1},
		}
		resp, err := s.msgServer.AddSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 1)
		s.Require().Equal(sla1, slas[0])
	})

	s.Run("accepts a req with valid SLAs", func() {
		req := &slatypes.MsgAddSLAs{
			Authority: s.authority.String(),
			SLAs:      []slatypes.PriceFeedSLA{sla2, sla3},
		}
		resp, err := s.msgServer.AddSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 2)
		s.Require().Equal(sla2, slas[0])
		s.Require().Equal(sla3, slas[1])
	})
}

func (s *KeeperTestSuite) TestMsgRemoveSLAs() {
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

	s.Run("accepts a req with no SLAs", func() {
		req := &slatypes.MsgRemoveSLAs{
			Authority: s.authority.String(),
		}
		resp, err := s.msgServer.RemoveSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
	})

	s.Run("accepts a req removing one sla that does not exist", func() {
		req := &slatypes.MsgRemoveSLAs{
			Authority: s.authority.String(),
			IDs:       []string{"does not exist"},
		}
		resp, err := s.msgServer.RemoveSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
	})

	s.Run("removes a sla single sla", func() {
		s.Require().NoError(s.keeper.SetSLA(s.ctx, sla1))

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 1)
		s.Require().Equal(sla1, slas[0])

		req := &slatypes.MsgRemoveSLAs{
			Authority: s.authority.String(),
			IDs:       []string{sla1.ID},
		}
		resp, err := s.msgServer.RemoveSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		slas, err = s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 0)
	})

	s.Run("removes a sla single sla with some feeds in state", func() {
		cons1 := sdk.ConsAddress("cons1")
		cons2 := sdk.ConsAddress("cons2")
		cp1 := slinkytypes.NewCurrencyPair("BTC", "USD")

		feed1, err := slatypes.NewPriceFeed(10, cons1, cp1, sla1.ID)
		s.Require().NoError(err)

		feed2, err := slatypes.NewPriceFeed(10, cons2, cp1, sla1.ID)
		s.Require().NoError(err)

		s.Require().NoError(s.keeper.SetPriceFeed(s.ctx, feed1))
		s.Require().NoError(s.keeper.SetPriceFeed(s.ctx, feed2))

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, sla1.ID)
		s.Require().NoError(err)
		s.Require().Len(feeds, 2)

		s.Require().NoError(s.keeper.SetSLA(s.ctx, sla1))
		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 1)
		s.Require().Equal(sla1, slas[0])

		req := &slatypes.MsgRemoveSLAs{
			Authority: s.authority.String(),
			IDs:       []string{sla1.ID},
		}
		resp, err := s.msgServer.RemoveSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		slas, err = s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 0)

		feeds, err = s.keeper.GetAllPriceFeeds(s.ctx, sla1.ID)
		s.Require().NoError(err)
		s.Require().Len(feeds, 0)
	})

	s.Run("removes multiple slas", func() {
		s.Require().NoError(s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla1, sla2, sla3}))

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 3)
		s.Require().Equal(sla1, slas[0])
		s.Require().Equal(sla2, slas[1])
		s.Require().Equal(sla3, slas[2])

		req := &slatypes.MsgRemoveSLAs{
			Authority: s.authority.String(),
			IDs:       []string{sla1.ID, sla2.ID},
		}
		resp, err := s.msgServer.RemoveSLAs(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		slas, err = s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 1)
		s.Require().Equal(sla3, slas[0])
	})
}

func (s *KeeperTestSuite) TestMsgParams() {
	s.Run("accepts a req with no params", func() {
		req := &slatypes.MsgParams{
			Authority: s.authority.String(),
		}
		resp, err := s.msgServer.Params(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
	})

	s.Run("accepts a req with params", func() {
		req := &slatypes.MsgParams{
			Authority: s.authority.String(),
			Params:    slatypes.DefaultParams(),
		}
		resp, err := s.msgServer.Params(s.ctx, req)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(req.Params, params)
	})
}
