package keeper_test

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/mock"

	"github.com/skip-mev/chaintestutil/sample"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/marketmap/keeper"
	"github.com/skip-mev/connect/v2/x/marketmap/types"
	mmmocks "github.com/skip-mev/connect/v2/x/marketmap/types/mocks"
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
	s.Require().Equal([]connecttypes.CurrencyPair{btcusdt.Ticker.CurrencyPair, usdtusd.Ticker.CurrencyPair}, cps)

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
							NormalizeByPair: &connecttypes.CurrencyPair{
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
				NormalizeByPair: &connecttypes.CurrencyPair{
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

func (s *KeeperTestSuite) TestMsgServerUpsertMarkets() {
	hooks := mmmocks.NewMarketMapHooks(s.T())

	// init keeper w/ mocked hooks
	s.keeper = s.initKeeperWithHooks(hooks)

	msgServer := keeper.NewMsgServer(s.keeper)

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgUpsertMarkets{
			Authority: "invalid",
		}
		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority (valid bech32)", func() {
		msg := &types.MsgUpsertMarkets{
			Authority: sdk.AccAddress("invalid").String(),
		}
		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	// set a market in the map
	s.Run("unable to process nil request", func() {
		resp, err := msgServer.UpsertMarkets(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("if a market does not exist, create it + error if hook errors", func() {
		msg := &types.MsgUpsertMarkets{
			Authority: s.marketAuthorities[0],
			Markets: []types.Market{
				usdcusd,
			},
		}

		err := fmt.Errorf("hook error")
		hooks.On("AfterMarketCreated", s.ctx, usdcusd).Return(err).Once()

		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("if a market does not exist, create it", func() {
		msg := &types.MsgUpsertMarkets{
			Authority: s.marketAuthorities[0],
			Markets: []types.Market{
				btcusdt,
			},
		}

		hooks.On("AfterMarketCreated", mock.Anything, btcusdt).Return(nil).Once()

		s.ctx = s.ctx.WithBlockHeight(12)

		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Len(resp.MarketUpdates, 0)

		// check that the market now exists
		found, err := s.keeper.HasMarket(s.ctx, btcusdt.Ticker.String())
		s.Require().NoError(err)
		s.Require().True(found)

		// check that last updated is correct
		lastUpdated, err := s.keeper.GetLastUpdated(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(uint64(s.ctx.BlockHeight()), lastUpdated) //nolint:gosec

		// check that the emitted events are correct (get the last event)
		event := s.ctx.EventManager().Events()[len(s.ctx.EventManager().Events())-1]
		s.Require().Equal(types.EventTypeCreateMarket, event.Type)

		// require that attributes are correct
		s.Require().Equal(btcusdt.Ticker.String(), event.Attributes[0].Value)
		s.Require().Equal(strconv.FormatUint(btcusdt.Ticker.Decimals, 10), event.Attributes[1].Value)
		s.Require().Equal(strconv.FormatUint(btcusdt.Ticker.MinProviderCount, 10), event.Attributes[2].Value)
		s.Require().Equal(btcusdt.Ticker.Metadata_JSON, event.Attributes[3].Value)
	})

	s.Run("if a market exists, update it but fail when the hook fails", func() {
		msg := &types.MsgUpsertMarkets{
			Authority: s.marketAuthorities[0],
			Markets: []types.Market{
				usdtusd,
				btcusdt,
			},
		}

		err := fmt.Errorf("hook error")
		hooks.On("AfterMarketCreated", mock.Anything, usdtusd).Return(nil).Once()
		hooks.On("AfterMarketUpdated", mock.Anything, btcusdt).Return(err).Once()

		s.ctx = s.ctx.WithBlockHeight(13)

		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("if a market exists, update it", func() {
		btcusdt.Ticker.MinProviderCount = 5

		msg := &types.MsgUpsertMarkets{
			Authority: s.marketAuthorities[0],
			Markets: []types.Market{
				ethusdt,
				btcusdt,
			},
		}

		hooks.On("AfterMarketUpdated", mock.Anything, btcusdt).Return(nil).Once()
		hooks.On("AfterMarketCreated", mock.Anything, ethusdt).Return(nil).Once()

		s.ctx = s.ctx.WithBlockHeight(13)

		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		s.Require().Len(resp.MarketUpdates, 0)

		// check that the market still exists
		found, err := s.keeper.HasMarket(s.ctx, btcusdt.Ticker.String())
		s.Require().NoError(err)
		s.Require().True(found)

		// check that the market now exists
		found, err = s.keeper.HasMarket(s.ctx, ethusdt.Ticker.String())
		s.Require().NoError(err)
		s.Require().True(found)

		// check that last updated is correct
		lastUpdated, err := s.keeper.GetLastUpdated(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(uint64(s.ctx.BlockHeight()), lastUpdated) //nolint:gosec

		// check that the emitted events are correct (get the last event)
		event := s.ctx.EventManager().Events()[len(s.ctx.EventManager().Events())-1]
		s.Require().Equal(types.EventTypeUpdateMarket, event.Type)

		// require that attributes are correct
		s.Require().Equal(btcusdt.Ticker.String(), event.Attributes[0].Value)
		s.Require().Equal(strconv.FormatUint(btcusdt.Ticker.Decimals, 10), event.Attributes[1].Value)
		s.Require().Equal(strconv.FormatUint(btcusdt.Ticker.MinProviderCount, 10), event.Attributes[2].Value)
		s.Require().Equal(btcusdt.Ticker.Metadata_JSON, event.Attributes[3].Value)
	})

	s.Run("if a market breaks the market-map, fail", func() {
		msg := &types.MsgUpsertMarkets{
			Authority: s.marketAuthorities[0],
			Markets: []types.Market{
				{
					Ticker: ethusdt.Ticker,
					ProviderConfigs: []types.ProviderConfig{
						{
							Name:           "kucoin",
							OffChainTicker: "eth-usdt",
							NormalizeByPair: &connecttypes.CurrencyPair{
								Base:  "INVALID",
								Quote: "PAIR",
							},
						},
					},
				},
			},
		}
		hooks.On("AfterMarketUpdated", mock.Anything, mock.Anything).Return(nil).Once()

		s.ctx = s.ctx.WithBlockHeight(13)

		resp, err := msgServer.UpsertMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})
}

func (s *KeeperTestSuite) TestMsgServerRemoveMarkets() {
	msgServer := keeper.NewMsgServer(s.keeper)

	s.Run("unable to process nil request", func() {
		resp, err := msgServer.RemoveMarkets(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to process for invalid authority", func() {
		msg := &types.MsgRemoveMarkets{
			Admin: "invalid",
		}
		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("only remove existing markets - no error", func() {
		msg := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{"BTC/USD", "ETH/USDT"},
		}
		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().Equal([]string{}, resp.DeletedMarkets)
	})

	s.Run("unable to remove non-existent market - single", func() {
		msg := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{"BTC/USD"},
		}
		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().Equal([]string{}, resp.DeletedMarkets)
	})

	s.Run("able to remove disabled market", func() {
		copyBTC := btcusdt
		copyBTC.Ticker.Enabled = false

		msg := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{copyBTC.Ticker.String()},
		}

		err := s.keeper.CreateMarket(s.ctx, copyBTC)
		s.Require().NoError(err)

		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().Equal([]string{copyBTC.Ticker.String()}, resp.DeletedMarkets)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyBTC.Ticker.String())
		s.Require().Error(err)
	})

	s.Run("do not remove enabled market", func() {
		copyBTC := btcusdt
		copyBTC.Ticker.Enabled = true

		err := s.keeper.CreateMarket(s.ctx, copyBTC)
		s.Require().NoError(err)

		msg := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{copyBTC.Ticker.String()},
		}

		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)

		// market should exist
		_, err = s.keeper.GetMarket(s.ctx, copyBTC.Ticker.String())
		s.Require().NoError(err)

		// update market to be disabled
		copyBTC.Ticker.Enabled = false

		err = s.keeper.UpdateMarket(s.ctx, copyBTC)
		s.Require().NoError(err)

		// remove
		resp, err = msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().Equal([]string{copyBTC.Ticker.String()}, resp.DeletedMarkets)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyBTC.Ticker.String())
		s.Require().Error(err)
	})

	s.Run("resulting state is invalid - 1", func() {
		// add a market that depends on the btc market
		copyBTC := btcusdt

		err := s.keeper.CreateMarket(s.ctx, copyBTC)
		s.Require().NoError(err)

		copyETH := ethusdt
		copyETH.ProviderConfigs = []types.ProviderConfig{
			{
				Name:           "normalized",
				OffChainTicker: "normalized",
				NormalizeByPair: &connecttypes.CurrencyPair{
					Base:  copyBTC.Ticker.CurrencyPair.Base,
					Quote: copyBTC.Ticker.CurrencyPair.Quote,
				},
			},
		}

		err = s.keeper.CreateMarket(s.ctx, copyETH)
		s.Require().NoError(err)

		msgRemoveBTC := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{copyBTC.Ticker.String()},
		}

		msgRemoveETH := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{copyETH.Ticker.String()},
		}

		resp, err := msgServer.RemoveMarkets(s.ctx, msgRemoveBTC)
		s.Require().Error(err)
		s.Require().Nil(resp)

		// remove dependent market first for valid state in 2 transaction
		resp, err = msgServer.RemoveMarkets(s.ctx, msgRemoveETH)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyETH.Ticker.String())
		s.Require().Error(err)
	})

	s.Run("remove both markets in one tx - no dependency", func() {
		// add a market that depends on the btc market
		copyBTC := btcusdt
		copyETH := ethusdt

		err := s.keeper.CreateMarket(s.ctx, copyBTC)
		s.Require().NoError(err)

		err = s.keeper.CreateMarket(s.ctx, copyETH)
		s.Require().NoError(err)

		msg := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{copyBTC.Ticker.String(), copyETH.Ticker.String()},
		}

		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyBTC.Ticker.String())
		s.Require().Error(err)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyETH.Ticker.String())
		s.Require().Error(err)
	})

	s.Run("remove both markets in one tx - with dependency", func() {
		// add a market that depends on the btc market
		copyBTC := btcusdt

		err := s.keeper.CreateMarket(s.ctx, copyBTC)
		s.Require().NoError(err)

		copyETH := ethusdt
		copyETH.ProviderConfigs = []types.ProviderConfig{
			{
				Name:           "normalized",
				OffChainTicker: "normalized",
				NormalizeByPair: &connecttypes.CurrencyPair{
					Base:  copyBTC.Ticker.CurrencyPair.Base,
					Quote: copyBTC.Ticker.CurrencyPair.Quote,
				},
			},
		}

		err = s.keeper.CreateMarket(s.ctx, copyETH)
		s.Require().NoError(err)

		msg := &types.MsgRemoveMarkets{
			Admin:   s.admin,
			Markets: []string{copyBTC.Ticker.String(), copyETH.Ticker.String()},
		}

		resp, err := msgServer.RemoveMarkets(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyBTC.Ticker.String())
		s.Require().Error(err)

		// market should not exist
		_, err = s.keeper.GetMarket(s.ctx, copyETH.Ticker.String())
		s.Require().Error(err)
	})
}
