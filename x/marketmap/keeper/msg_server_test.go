package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/skip-mev/chaintestutil/sample"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestMsgServerCreateMarkets() {
	msgServer := keeper.NewMsgServer(s.keeper)

	// create initial markets
	msg := &types.MsgCreateMarkets{
		Authority: s.marketAuthorities[1],
		CreateMarkets: []types.Market{
			btcusdt,
			usdtusd,
		},
	}
	resp, err := msgServer.CreateMarkets(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// query the market map
	mm, err := s.keeper.GetAllMarkets(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(types.MarketMap{Markets: mm}, types.MarketMap{
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
			Authority: "invalid",
		}
		resp, err = msgServer.CreateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority (valid bech32)", func() {
		msg = &types.MsgCreateMarkets{
			Authority: sdk.AccAddress("invalid").String(),
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
			Authority: s.authority.String(),
			CreateMarkets: []types.Market{
				btcusdt,
			},
		}
		resp, err = msgServer.CreateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to create market with normalize by that is not on chain tickers", func() {
		msg = &types.MsgCreateMarkets{
			Authority: s.authority.String(),
			CreateMarkets: []types.Market{
				{
					Ticker: ethusdt.Ticker,
					ProviderConfigs: []types.ProviderConfig{
						{
							Name:           "kucoin",
							OffChainTicker: "eth-usdt",
							NormalizeByPair: &slinkytypes.CurrencyPair{
								Base:  "INVALID",
								Quote: "PAIR",
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

	// create initial markets
	createMsg := &types.MsgCreateMarkets{
		Authority: s.marketAuthorities[0],
		CreateMarkets: []types.Market{
			btcusdt,
			usdtusd,
		},
	}
	createResp, err := msgServer.CreateMarkets(s.ctx, createMsg)
	s.Require().NoError(err)
	s.Require().NotNil(createResp)

	// query the market map
	mm, err := s.keeper.GetAllMarkets(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(types.MarketMap{Markets: mm}, types.MarketMap{
		Markets: map[string]types.Market{
			btcusdt.Ticker.String(): btcusdt,
			usdtusd.Ticker.String(): usdtusd,
		},
	})

	// TODO: test hooks

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgUpdateMarkets{
			Authority: "invalid",
		}
		resp, err := msgServer.UpdateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority (valid bech32)", func() {
		msg := &types.MsgUpdateMarkets{
			Authority: sdk.AccAddress("invalid").String(),
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
			Authority: s.marketAuthorities[2],
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
		tickerUpdate.ProviderConfigs = []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "btc-usdc",
				NormalizeByPair: &slinkytypes.CurrencyPair{
					Base:  "INVALID",
					Quote: "PAIR",
				},
			},
		}

		msg := &types.MsgUpdateMarkets{
			Authority: s.authority.String(),
			UpdateMarkets: []types.Market{
				tickerUpdate,
			},
		}
		resp, err := msgServer.UpdateMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to update market if it does not already exist", func() {
		msg := &types.MsgUpdateMarkets{
			Authority: s.authority.String(),
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
		resp, err := msgServer.UpdateParams(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgParams{
			Authority: "invalid",
		}
		resp, err := msgServer.UpdateParams(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority (valid bech32)", func() {
		msg := &types.MsgParams{
			Authority: sdk.AccAddress("invalid").String(),
		}
		resp, err := msgServer.UpdateParams(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("accepts a req with valid params", func() {
		msg := &types.MsgParams{
			Authority: s.authority.String(),
			Params: types.Params{
				MarketAuthorities: []string{authtypes.NewModuleAddress(govtypes.ModuleName).String(), sample.Address(r)},
				Admin:             sample.Address(r),
			},
		}
		resp, err := msgServer.UpdateParams(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(msg.Params, params)
	})
}

func (s *KeeperTestSuite) TestMsgServerRemoveMarketAuthorities() {
	msgServer := keeper.NewMsgServer(s.keeper)

	s.Run("unable to process nil request", func() {
		resp, err := msgServer.RemoveMarketAuthorities(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgRemoveMarketAuthorities{
			Admin: "invalid",
		}
		resp, err := msgServer.RemoveMarketAuthorities(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("accepts a req that removes one authority", func() {
		msg := &types.MsgRemoveMarketAuthorities{
			Admin:           s.admin,
			RemoveAddresses: []string{s.marketAuthorities[0]},
		}
		resp, err := msgServer.RemoveMarketAuthorities(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		// check new authorities
		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(s.marketAuthorities[1:], params.MarketAuthorities)

		// reset
		s.Require().NoError(s.keeper.SetParams(s.ctx, types.Params{
			MarketAuthorities: s.marketAuthorities,
			Admin:             s.admin,
		}))
	})

	s.Run("accepts a req that removes multiple authorities", func() {
		msg := &types.MsgRemoveMarketAuthorities{
			Admin:           s.admin,
			RemoveAddresses: []string{s.marketAuthorities[0], s.marketAuthorities[2]},
		}
		resp, err := msgServer.RemoveMarketAuthorities(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		// check new authorities
		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal([]string{s.marketAuthorities[1]}, params.MarketAuthorities)

		// reset
		s.Require().NoError(s.keeper.SetParams(s.ctx, types.Params{
			MarketAuthorities: s.marketAuthorities,
			Admin:             s.admin,
		}))
	})

	s.Run("unable to accept a req that removes more authorities than exist in state", func() {
		msg := &types.MsgRemoveMarketAuthorities{
			Admin:           s.admin,
			RemoveAddresses: []string{sample.Address(r), sample.Address(r), sample.Address(r), sample.Address(r), sample.Address(r), sample.Address(r), sample.Address(r), sample.Address(r)},
		}
		resp, err := msgServer.RemoveMarketAuthorities(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)

		// reset
		s.Require().NoError(s.keeper.SetParams(s.ctx, types.Params{
			MarketAuthorities: s.marketAuthorities,
			Admin:             s.admin,
		}))
	})
}
