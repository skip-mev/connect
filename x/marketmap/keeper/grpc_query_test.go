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
			Version:     10,
		}

		s.Require().Equal(expected, resp)
	})

	s.Run("run query with state", func() {
		expectedMarketMap := types.MarketMap{
			Tickers:   make(map[string]types.Ticker),
			Paths:     make(map[string]types.Paths),
			Providers: make(map[string]types.Providers),
		}
		for _, ticker := range markets.tickers {
			marketPaths, ok := markets.paths[ticker.String()]
			s.Require().True(ok)
			marketProviders, ok := markets.providers[ticker.String()]
			s.Require().True(ok)
			s.Require().NoError(s.keeper.CreateMarket(s.ctx, ticker, marketPaths, marketProviders))
			expectedMarketMap.Tickers[ticker.String()] = ticker
			expectedMarketMap.Paths[ticker.String()] = marketPaths
			expectedMarketMap.Providers[ticker.String()] = marketProviders
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

func (s *KeeperTestSuite) TestQueryServerParams() {
	params := types.DefaultParams()
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	qs := keeper.NewQueryServer(s.keeper)

	s.Run("run valid request", func() {
		resp, err := qs.Params(s.ctx, &types.ParamsRequest{})
		s.Require().NoError(err)

		s.Require().Equal(params, resp.Params)
	})

	s.Run("run invalid nil request", func() {
		_, err := qs.Params(s.ctx, nil)
		s.Require().Error(err)
	})
}
