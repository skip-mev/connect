package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

func (s *KeeperTestSuite) TestBeginBlocker() {
	s.Run("run with no state", func() {
		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetRemovedCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))
	})

	s.Run("run with 1 in state - 1 removed", func() {
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, slinkytypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, slinkytypes.CurrencyPair{Base: "test", Quote: "coin1"}))

		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetRemovedCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))

		cps, err := s.oracleKeeper.GetPrevBlockCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(0))
	})

	s.Run("run with 2 in state - 1 removed", func() {
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, slinkytypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, slinkytypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, slinkytypes.CurrencyPair{Base: "test", Quote: "coin2"}))

		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetRemovedCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))

		cps, err := s.oracleKeeper.GetPrevBlockCPCounter(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(1))
	})
}
