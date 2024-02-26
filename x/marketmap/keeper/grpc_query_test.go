package keeper_test

import (
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestMarketMap() {
	qs := keeper.NewQueryServer(s.keeper)
	s.ctx = s.ctx.WithChainID("test-chain")

	s.Run("invalid for nil request", func() {
		_, err := qs.MarketMap(s.ctx, nil)
		s.Require().Error(err)
	})

	s.Run("run query with no state", func() {
		resp, err := qs.MarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		expected := &types.GetMarketMapResponse{
			MarketMap: types.MarketMap{
				Tickers:   make(map[string]types.Ticker),
				Paths:     make(map[string]types.Paths),
				Providers: make(map[string]types.Providers),
			},
			LastUpdated: s.ctx.BlockHeight(),
			Version:     10,
			ChainId:     "test-chain",
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

		resp, err := qs.MarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		expected := &types.GetMarketMapResponse{
			MarketMap:   expectedMarketMap,
			LastUpdated: s.ctx.BlockHeight(),
			ChainId:     "test-chain",
		}

		s.Require().Equal(expected.LastUpdated, resp.LastUpdated)
		s.Require().Equal(expected.MarketMap, resp.MarketMap)
	})
}

func (s *KeeperTestSuite) TestParams() {
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

func (s *KeeperTestSuite) TestLastUpdated() {
	qs := keeper.NewQueryServer(s.keeper)
	// set initial states
	for _, ticker := range markets.tickers {
		marketPaths, ok := markets.paths[ticker.String()]
		s.Require().True(ok)
		marketProviders, ok := markets.providers[ticker.String()]
		s.Require().True(ok)
		s.Require().NoError(s.keeper.CreateMarket(s.ctx, ticker, marketPaths, marketProviders))
	}

	s.Run("run valid request", func() {
		resp, err := qs.LastUpdated(s.ctx, &types.GetLastUpdatedRequest{})
		s.Require().NoError(err)

		s.Require().Equal(s.ctx.BlockHeight(), resp.LastUpdated)
	})

	s.Run("run invalid nil request", func() {
		_, err := qs.LastUpdated(s.ctx, nil)
		s.Require().Error(err)
	})
}
