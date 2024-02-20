package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestQueryServeGetMarketMap() {
	cp1 := slinkytypes.CurrencyPair{Base: "BTC", Quote: "ETH"}
	aggCfg1 := types.PathsConfig{
		Ticker: types.Ticker{
			CurrencyPair:     cp1,
			Decimals:         0,
			MinProviderCount: 0,
		},
		Paths: []types.Path{
			{Operations: []types.Operation{{Ticker: types.Ticker{CurrencyPair: cp1}}}},
		},
	}
	btcEthTickerConfig := types.TickerConfig{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "BTC",
				Quote: "ETH",
			},
			Decimals:         8,
			MinProviderCount: 1,
		},
		OffChainTicker: "BTC-ETH",
	}
	marketCfg1 := types.MarketConfig{
		Name: "provider1",
		TickerConfigs: map[string]types.TickerConfig{
			"BTC/ETH": btcEthTickerConfig,
		},
	}

	expectedMM := types.AggregateMarketConfig{
		MarketConfigs: map[string]types.MarketConfig{
			"provider1": marketCfg1,
		},
		TickerConfigs: map[string]types.PathsConfig{
			cp1.String(): aggCfg1,
		},
	}

	qs := keeper.NewQueryServer(s.keeper)

	s.Run("run valid request", func() {
		const (
			testBlockHeight int64  = 9
			expectedVersion uint64 = 10
		)

		s.Require().NoError(s.keeper.CreateMarketConfig(s.ctx.WithBlockHeight(testBlockHeight), marketCfg1))
		s.Require().NoError(s.keeper.CreateAggregationConfig(s.ctx.WithBlockHeight(testBlockHeight), aggCfg1))

		resp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		mm, err := s.keeper.GetMarketMap(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(*mm, resp.MarketMap)
		s.Require().Equal(expectedMM, resp.MarketMap)

		// check if last updated is the ctx value used for the keeper writes.
		s.Require().Equal(testBlockHeight, resp.LastUpdated)

		s.Require().Equal(expectedVersion, resp.Version)
	})

	s.Run("run invalid nil request", func() {
		_, err := qs.GetMarketMap(s.ctx, nil)
		s.Require().Error(err)
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
