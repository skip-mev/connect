package keeper_test

import "github.com/skip-mev/slinky/x/oracle/types"

func (s *KeeperTestSuite) TestBeginBlocker() {
	s.Run("run with no state", func() {
		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetRemovedCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))
	})

	s.Run("run with invalid state 1 removed, 0 in state (cannot happen)", func() {
		s.Require().NoError(s.oracleKeeper.IncrementRemovedCPCounter(s.ctx))
		s.Require().Error(s.oracleKeeper.BeginBlocker(s.ctx))

		// reset state
		s.Require().NotPanics(func() {
			s.oracleKeeper.InitGenesis(s.ctx, *types.DefaultGenesisState())
		})
	})

	s.Run("run with 1 in state - 1 removed", func() {
		s.Require().NoError(s.oracleKeeper.IncrementRemovedCPCounter(s.ctx))
		s.Require().NoError(s.oracleKeeper.IncrementCPCounter(s.ctx))

		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetRemovedCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))

		cps, err := s.oracleKeeper.GetCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(0))

		// reset state
		s.Require().NotPanics(func() {
			s.oracleKeeper.InitGenesis(s.ctx, *types.DefaultGenesisState())
		})
	})

	s.Run("run with 2 in state - 1 removed", func() {
		s.Require().NoError(s.oracleKeeper.IncrementRemovedCPCounter(s.ctx))
		s.Require().NoError(s.oracleKeeper.IncrementCPCounter(s.ctx))
		s.Require().NoError(s.oracleKeeper.IncrementCPCounter(s.ctx))

		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetRemovedCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))

		cps, err := s.oracleKeeper.GetCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(1))

		// reset state
		s.Require().NotPanics(func() {
			s.oracleKeeper.InitGenesis(s.ctx, *types.DefaultGenesisState())
		})
	})
}
