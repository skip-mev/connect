package keeper_test

import (
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

func (s *KeeperTestSuite) TestBeginBlocker() {
	s.Run("run with no state", func() {
		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))
	})

	s.Run("run with 1 in state - 1 removed", func() {
		// Create the currency pair.
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		cps, err := s.oracleKeeper.GetNumCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(1))
		removed, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removed, uint64(0))

		// Remove the currency pair.
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		cps, err = s.oracleKeeper.GetNumCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(0))
		removed, err = s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removed, uint64(1))

		// Begin blocker should reset the removed count.
		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))
	})

	s.Run("run with 2 in state - 1 removed", func() {
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.RemoveCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin1"}))
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{Base: "test", Quote: "coin2"}))

		s.Require().NoError(s.oracleKeeper.BeginBlocker(s.ctx))
		removes, err := s.oracleKeeper.GetNumRemovedCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(removes, uint64(0))

		cps, err := s.oracleKeeper.GetNumCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(cps, uint64(1))
	})
}
