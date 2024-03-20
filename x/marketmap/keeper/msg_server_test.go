package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestMsgServerCreateMarkets() {
	msgServer := keeper.NewMsgServer(s.keeper)
	qs := keeper.NewQueryServer(s.keeper)

	// create initial markets
	msg := &types.MsgCreateMarkets{
		Signer: s.authority.String(),
		CreateMarkets: []types.Market{
			{
				Ticker:    btcusdt.Ticker,
				Providers: btcusdt.Providers,
				Paths:     btcusdt.Paths,
			},
			{
				Ticker:    usdtusd.Ticker,
				Providers: usdtusd.Providers,
				Paths:     usdtusd.Paths,
			},
		},
	}
	resp, err := msgServer.CreateMarkets(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// query the market map
	queryResp, err := qs.MarketMap(s.ctx, &types.GetMarketMapRequest{})
	s.Require().NoError(err)
	s.Require().Equal(queryResp.MarketMap, types.MarketMap{
		Markets: map[string]types.Market{
			btcusdt.Ticker.String(): btcusdt,
			usdtusd.Ticker.String(): usdtusd,
		},
	})

	// query the oracle module to see if they were created via hooks
	cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
	s.Require().Equal([]slinkytypes.CurrencyPair{btcusdt.Ticker.CurrencyPair, usdtusd.Ticker.CurrencyPair}, cps)

	s.Run("unable to process for invalid authority", func() {
		msg = &types.MsgCreateMarkets{
			Signer: "invalid",
		}
		resp, err = msgServer.CreateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	// set a market in the map
	s.Run("unable to process nil request", func() {
		resp, err = msgServer.CreateMarkets(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to create market that already exists", func() {
		msg = &types.MsgCreateMarkets{
			Signer: s.authority.String(),
			CreateMarkets: []types.Market{
				{
					Ticker:    btcusdt.Ticker,
					Providers: btcusdt.Providers,
					Paths:     btcusdt.Paths,
				},
			},
		}
		resp, err = msgServer.CreateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to create market with paths that are not on chain tickers", func() {
		msg = &types.MsgCreateMarkets{
			Signer: s.authority.String(),
			CreateMarkets: []types.Market{
				{
					Ticker:    ethusdt.Ticker,
					Providers: ethusdt.Providers,
					Paths: types.Paths{
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: slinkytypes.CurrencyPair{
											Base:  "ETHEREUM",
											Quote: "MOG",
										},
									},
									{
										CurrencyPair: slinkytypes.CurrencyPair{
											Base:  "MOG",
											Quote: "USDT",
										},
									},
								},
							},
						},
					},
				},
			},
		}
		resp, err = msgServer.CreateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})
}

func (s *KeeperTestSuite) TestMsgServerUpdateMarkets() {
	msgServer := keeper.NewMsgServer(s.keeper)
	qs := keeper.NewQueryServer(s.keeper)

	// create initial markets
	createMsg := &types.MsgCreateMarkets{
		Signer: s.authority.String(),
		CreateMarkets: []types.Market{
			btcusdt,
			usdtusd,
		},
	}
	createResp, err := msgServer.CreateMarkets(s.ctx, createMsg)
	s.Require().NoError(err)
	s.Require().NotNil(createResp)

	// query the market map
	queryResp, err := qs.MarketMap(s.ctx, &types.GetMarketMapRequest{})
	s.Require().NoError(err)
	s.Require().Equal(queryResp.MarketMap, types.MarketMap{
		Markets: map[string]types.Market{
			btcusdt.Ticker.String(): btcusdt,
			usdtusd.Ticker.String(): usdtusd,
		},
	})

	// query the oracle module to see if they were created via hooks
	cps := s.oracleKeeper.GetAllCurrencyPairs(s.ctx)
	s.Require().Equal([]slinkytypes.CurrencyPair{btcusdt.Ticker.CurrencyPair, usdtusd.Ticker.CurrencyPair}, cps)

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgUpdateMarkets{
			Signer: "invalid",
		}
		resp, err := msgServer.UpdateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	// set a market in the map
	s.Run("unable to process nil request", func() {
		resp, err := msgServer.UpdateMarkets(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("able to update market that already exists", func() {
		tickerUpdate := btcusdt
		tickerUpdate.Ticker.Decimals = 1

		msg := &types.MsgUpdateMarkets{
			Signer: s.authority.String(),
			UpdateMarkets: []types.Market{
				tickerUpdate,
			},
		}
		resp, err := msgServer.UpdateMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
	})

	s.Run("unable to update when paths refers to non-existent ticker", func() {
		tickerUpdate := btcusdt
		tickerUpdate.Ticker.Decimals = 1

		//nolint: gocritic
		updatePaths := append(btcusdt.Paths.Paths, types.Path{
			Operations: []types.Operation{
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "COINBIT",
						Quote: "USDT",
					},
					Invert: false,
				},
			},
		},
		)

		msg := &types.MsgUpdateMarkets{
			Signer: s.authority.String(),
			UpdateMarkets: []types.Market{
				{
					Ticker:    tickerUpdate.Ticker,
					Providers: btcusdt.Providers,
					Paths:     types.Paths{Paths: updatePaths},
				},
			},
		}
		resp, err := msgServer.UpdateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to update market if it does not already exist", func() {
		msg := &types.MsgUpdateMarkets{
			Signer: s.authority.String(),
			UpdateMarkets: []types.Market{
				ethusdt,
			},
		}
		resp, err := msgServer.UpdateMarkets(s.ctx, msg)
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
