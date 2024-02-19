package keeper_test

import (
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestQueryServer() {
	qs := keeper.NewQueryServer(s.keeper)

	s.Run("invalid for nil request", func() {
		_, err := qs.GetMarketMap(s.ctx, nil)
		s.Require().Error(err)
	})

	s.Run("run query with no state", func() {
		resp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		expected := &types.GetMarketMapResponse{
			MarketMap:   types.TickersConfig{Tickers: nil},
			LastUpdated: s.ctx.BlockHeight(),
		}

		s.Require().Equal(expected, resp)
	})
}
