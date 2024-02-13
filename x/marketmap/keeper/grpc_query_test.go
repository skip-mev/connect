package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestQueryServer() {
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
		s.Require().NoError(s.keeper.CreateMarketConfig(s.ctx, marketCfg1))
		s.Require().NoError(s.keeper.CreateAggregationConfig(s.ctx, aggCfg1))

		resp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		mm, err := s.keeper.GetMarketMap(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(*mm, resp.MarketMap)
		s.Require().Equal(expectedMM, resp.MarketMap)
	})

	s.Run("run invalid nil request", func() {
		_, err := qs.GetMarketMap(s.ctx, nil)
		s.Require().Error(err)
	})
}
