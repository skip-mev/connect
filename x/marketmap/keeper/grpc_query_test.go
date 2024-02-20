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
			MarketMap: types.MarketMap{
				Tickers:   make(map[string]types.Ticker),
				Paths:     make(map[string]types.Paths),
				Providers: make(map[string]types.Providers),
			},
			LastUpdated: s.ctx.BlockHeight(),
		}

		s.Require().Equal(expected, resp)
	})

	s.Run("run query with state", func() {
		expectedMarketMap := types.MarketMap{
			Tickers:   make(map[string]types.Ticker),
			Paths:     make(map[string]types.Paths),
			Providers: make(map[string]types.Providers),
		}
		for _, ticker := range tickers {
			s.Require().NoError(s.keeper.CreateMarket(s.ctx, ticker, types.Paths{Paths: ticker.Paths}, types.Providers{Providers: ticker.Providers}))
			expectedMarketMap.Tickers[ticker.String()] = ticker
			expectedMarketMap.Paths[ticker.String()] = types.Paths{Paths: ticker.Paths}
			expectedMarketMap.Providers[ticker.String()] = types.Providers{Providers: ticker.Providers}

		}

		resp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		expected := &types.GetMarketMapResponse{
			MarketMap:   expectedMarketMap,
			LastUpdated: s.ctx.BlockHeight(),
		}

		s.Require().Equal(expected.LastUpdated, resp.LastUpdated)
		s.Require().Equal(expected.MarketMap, resp.MarketMap)
	})
}
