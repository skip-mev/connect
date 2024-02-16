package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestCreateMarket() {
	msgServer := keeper.NewMsgServer(s.keeper)
	qs := keeper.NewQueryServer(s.keeper)

	// set params
	paramsMsg := &types.MsgParams{
		Authority: s.authority.String(),
		Params:    types.DefaultParams(),
	}
	paramsResp, err := msgServer.Params(s.ctx, paramsMsg)
	s.Require().NoError(err)
	s.Require().NotNil(paramsResp)

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
			CurrencyPair:     cp1,
			Decimals:         8,
			MinProviderCount: 1,
		},
		OffChainTicker: "BTC-ETH",
	}
	marketCfg1 := types.MarketConfig{
		Name: "provider1",
		TickerConfigs: map[string]types.TickerConfig{
			cp1.String(): btcEthTickerConfig,
		},
	}

	cp2 := slinkytypes.CurrencyPair{Base: "BTC", Quote: "USD"}
	aggCfg2 := types.PathsConfig{
		Ticker: types.Ticker{
			CurrencyPair:     cp1,
			Decimals:         8,
			MinProviderCount: 1,
		},
		Paths: []types.Path{
			{Operations: []types.Operation{{Ticker: types.Ticker{CurrencyPair: cp2}}}},
		},
	}
	btcUsdTickerConfig := types.TickerConfig{
		Ticker: types.Ticker{
			CurrencyPair:     cp2,
			Decimals:         8,
			MinProviderCount: 1,
		},
		OffChainTicker: "BTC-USD",
	}

	expectedMM := types.AggregateMarketConfig{
		MarketConfigs: map[string]types.MarketConfig{
			"provider1": {
				Name: "provider1",
				TickerConfigs: map[string]types.TickerConfig{
					cp1.String(): btcEthTickerConfig,
					cp2.String(): btcUsdTickerConfig,
				},
			},
		},
		TickerConfigs: map[string]types.PathsConfig{
			cp1.String(): aggCfg1,
			cp2.String(): aggCfg2,
		},
	}

	s.Require().NoError(s.keeper.CreateMarketConfig(s.ctx, marketCfg1))
	s.Require().NoError(s.keeper.CreateAggregationConfig(s.ctx, aggCfg1))

	// set a market in the map
	s.Run("unable to process nil request", func() {
		resp, err := msgServer.CreateMarket(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgCreateMarket{
			Signer: "invalid",
		}
		resp, err := msgServer.CreateMarket(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for market already exists", func() {
		msg := &types.MsgCreateMarket{
			Signer: types.DefaultParams().MarketAuthority,
			Ticker: btcEthTickerConfig.Ticker,
			ProvidersToOffChainTickers: map[string]string{
				"provider1": btcEthTickerConfig.OffChainTicker,
			},
			Paths: aggCfg1.Paths,
		}
		resp, err := msgServer.CreateMarket(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("valid adds market", func() {
		msg := &types.MsgCreateMarket{
			Signer: types.DefaultParams().MarketAuthority,
			Ticker: btcUsdTickerConfig.Ticker,
			ProvidersToOffChainTickers: map[string]string{
				"provider1": btcUsdTickerConfig.OffChainTicker,
			},
			Paths: aggCfg2.Paths,
		}
		resp, err := msgServer.CreateMarket(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		queryResp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)

		mm, err := s.keeper.GetMarketMap(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(*mm, queryResp.MarketMap)
		s.Require().Equal(expectedMM, queryResp.MarketMap)
	})
}

func (s *KeeperTestSuite) TestParams() {
	msgServer := keeper.NewMsgServer(s.keeper)

	s.Run("unable to process nil request", func() {
		resp, err := msgServer.Params(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgParams{
			Authority: "invalid",
		}
		resp, err := msgServer.Params(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for version lower than current versions", func() {
		msg := &types.MsgParams{
			Authority: s.authority.String(),
			Params: types.NewParams(
				types.DefaultMarketAuthority,
				0,
			),
		}
		resp, err := msgServer.Params(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("accepts a req with no params", func() {
		msg := &types.MsgParams{
			Authority: s.authority.String(),
		}
		resp, err := msgServer.Params(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
	})

	s.Run("accepts a req with default params", func() {
		msg := &types.MsgParams{
			Authority: s.authority.String(),
			Params:    types.DefaultParams(),
		}
		resp, err := msgServer.Params(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(msg.Params, params)
	})
}
