package keeper_test

import (
	"cosmossdk.io/math"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func (s *KeeperTestSuite) TestAddSLAs() {
	s.Run("can set no slas", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{})
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Empty(slas)
	})

	sla1 := slatypes.NewPriceFeedSLA(
		"testID",
		10,
		math.LegacyMustNewDecFromStr("0.1"),
		math.LegacyMustNewDecFromStr("0.1"),
		5,
		5,
	)

	s.Run("can set a single sla", func() {
		err := s.keeper.SetSLA(s.ctx, sla1)
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 1)
		s.Require().Equal(sla1, slas[0])
	})

	sla2 := slatypes.NewPriceFeedSLA(
		"testID2",
		20,
		math.LegacyMustNewDecFromStr("0.2"),
		math.LegacyMustNewDecFromStr("0.2"),
		10,
		10,
	)

	sla3 := slatypes.NewPriceFeedSLA(
		"testID3",
		30,
		math.LegacyMustNewDecFromStr("0.3"),
		math.LegacyMustNewDecFromStr("0.3"),
		15,
		15,
	)

	s.Run("can set multiple slas", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla1, sla2, sla3})
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 3)
		s.Require().Equal(sla1, slas[0])
		s.Require().Equal(sla2, slas[1])
		s.Require().Equal(sla3, slas[2])
	})
}

func (s *KeeperTestSuite) TestRemoveSLAs() {
	s.Run("can remove no slas", func() {
		err := s.keeper.RemoveSLAs(s.ctx, []string{})
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Empty(slas)

		err = s.keeper.RemoveSLA(s.ctx, "testID")
		s.Require().NoError(err)
	})

	sla := slatypes.NewPriceFeedSLA(
		"testID",
		10,
		math.LegacyMustNewDecFromStr("0.1"),
		math.LegacyMustNewDecFromStr("0.1"),
		5,
		5,
	)

	s.Run("can remove a single sla", func() {
		err := s.keeper.SetSLA(s.ctx, sla)
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 1)
		s.Require().Equal(sla, slas[0])

		err = s.keeper.RemoveSLAs(s.ctx, []string{sla.ID})
		s.Require().NoError(err)

		slas, err = s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Empty(slas)

		err = s.keeper.RemoveSLA(s.ctx, sla.ID)
		s.Require().NoError(err)
	})

	sla2 := slatypes.NewPriceFeedSLA(
		"testID2",
		20,
		math.LegacyMustNewDecFromStr("0.2"),
		math.LegacyMustNewDecFromStr("0.2"),
		10,
		10,
	)

	sla3 := slatypes.NewPriceFeedSLA(
		"testID3",
		30,
		math.LegacyMustNewDecFromStr("0.3"),
		math.LegacyMustNewDecFromStr("0.3"),
		15,
		15,
	)

	s.Run("can add several slas and remove one", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla, sla2, sla3})
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 3)

		err = s.keeper.RemoveSLAs(s.ctx, []string{sla.ID})
		s.Require().NoError(err)

		slas, err = s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 2)
		s.Require().Equal(sla2, slas[0])
		s.Require().Equal(sla3, slas[1])
	})

	s.Run("can add several slas and remove all", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla, sla2, sla3})
		s.Require().NoError(err)

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(slas, 3)

		err = s.keeper.RemoveSLAs(s.ctx, []string{sla.ID, sla2.ID, sla3.ID})
		s.Require().NoError(err)

		slas, err = s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Empty(slas)
	})
}
