package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestMsgServerCreateMarket() {
	msgServer := keeper.NewMsgServer(s.keeper)
	qs := keeper.NewQueryServer(s.keeper)

	// create initial markets
	msg := &types.MsgUpdateMarketMap{
		Signer: s.authority.String(),
		CreateMarkets: []types.CreateMarket{
			{
				Ticker:    btcusdt,
				Providers: btcusdtProviders,
				Paths:     btcusdtPaths,
			},
			{
				Ticker:    usdtusd,
				Providers: usdtusdProviders,
				Paths:     usdtusdPaths,
			},
		},
	}
	resp, err := msgServer.UpdateMarketMap(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// query the market map
	queryResp, err := qs.MarketMap(s.ctx, &types.GetMarketMapRequest{})
	s.Require().NoError(err)
	s.Require().Equal(queryResp.MarketMap, types.MarketMap{
		Tickers: map[string]types.Ticker{
			btcusdt.String(): btcusdt,
			usdtusd.String(): usdtusd,
		},
		Paths: map[string]types.Paths{
			btcusdt.String(): btcusdtPaths,
			usdtusd.String(): usdtusdPaths,
		},
		Providers: map[string]types.Providers{
			btcusdt.String(): btcusdtProviders,
			usdtusd.String(): usdtusdProviders,
		},
	})

	// query the oracle module to see if they were created via hooks
	cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
	s.Require().Equal([]slinkytypes.CurrencyPair{btcusdt.CurrencyPair, usdtusd.CurrencyPair}, cps)

	s.Run("unable to process for invalid authority", func() {
		msg = &types.MsgUpdateMarketMap{
			Signer: "invalid",
		}
		resp, err = msgServer.UpdateMarketMap(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	// set a market in the map
	s.Run("unable to process nil request", func() {
		resp, err = msgServer.UpdateMarketMap(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to create market that already exists", func() {
		msg = &types.MsgUpdateMarketMap{
			Signer: s.authority.String(),
			CreateMarkets: []types.CreateMarket{
				{
					Ticker:    btcusdt,
					Providers: btcusdtProviders,
					Paths:     btcusdtPaths,
				},
			},
		}
		resp, err = msgServer.UpdateMarketMap(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to create market with paths that are not on chain tickers", func() {
		msg = &types.MsgUpdateMarketMap{
			Signer: s.authority.String(),
			CreateMarkets: []types.CreateMarket{
				{
					Ticker:    ethusdt,
					Providers: ethusdtProviders,
					Paths: types.Paths{
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: slinkytypes.CurrencyPair{
											Base:      "ETHEREUM",
											Quote:     "MOG",
											Delimiter: slinkytypes.DefaultDelimiter,
										},
									},
									{
										CurrencyPair: slinkytypes.CurrencyPair{
											Base:      "MOG",
											Quote:     "USDT",
											Delimiter: slinkytypes.DefaultDelimiter,
										},
									},
								},
							},
						},
					},
				},
			},
		}
		resp, err = msgServer.UpdateMarketMap(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})
}

func (s *KeeperTestSuite) TestMsgServerParams() {
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

	s.Run("unable to process a req with no params", func() {
		msg := &types.MsgParams{
			Authority: s.authority.String(),
		}
		resp, err := msgServer.Params(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("accepts a req with valid params", func() {
		msg := &types.MsgParams{
			Authority: s.authority.String(),
			Params: types.NewParams(
				types.DefaultMarketAuthority,
				11,
			),
		}
		resp, err := msgServer.Params(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(msg.Params, params)
	})
}
