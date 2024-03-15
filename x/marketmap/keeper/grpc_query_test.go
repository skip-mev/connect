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
				Markets: make(map[string]types.Market),
			},
			LastUpdated: uint64(s.ctx.BlockHeight()),
			Version:     10,
			ChainId:     "test-chain",
		}

		s.Require().Equal(expected, resp)
	})

	s.Run("run query with state", func() {
		expectedMarketMap := types.MarketMap{
			Markets: make(map[string]types.Market),
		}
		for _, market := range markets {
			s.Require().NoError(s.keeper.CreateMarket(s.ctx, market))
			expectedMarketMap.Markets[market.Ticker.String()] = market
		}

		resp, err := qs.MarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		expected := &types.GetMarketMapResponse{
			MarketMap:   expectedMarketMap,
			LastUpdated: uint64(s.ctx.BlockHeight()),
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
	for _, market := range markets {
		s.Require().NoError(s.keeper.CreateMarket(s.ctx, market))
	}

	s.Run("run valid request", func() {
		resp, err := qs.LastUpdated(s.ctx, &types.GetLastUpdatedRequest{})
		s.Require().NoError(err)

		s.Require().Equal(uint64(s.ctx.BlockHeight()), resp.LastUpdated)
	})

	s.Run("run invalid nil request", func() {
		_, err := qs.LastUpdated(s.ctx, nil)
		s.Require().Error(err)
	})
}
