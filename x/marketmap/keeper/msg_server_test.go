package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestCreateMarket() {
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

	queryResp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
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

	// set a market in the map
	s.Run("unable to process nil request", func() {
		resp, err := msgServer.UpdateMarketMap(s.ctx, nil)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	// TODO add with params
	// s.Run("unable to process for invalid authority", func() {
	//	msg := &types.MsgUpdateMarketMap{
	//		Signer: "invalid",
	//	}
	//	resp, err := msgServer.UpdateMarketMap(s.ctx, msg)
	//	s.Require().Error(err)
	// 	s.Require().Nil(resp)
	// })

	s.Run("unable to create market that already exists", func() {
		msg := &types.MsgUpdateMarketMap{
			Signer: s.authority.String(),
			CreateMarkets: []types.CreateMarket{
				{
					Ticker:    btcusdt,
					Providers: btcusdtProviders,
					Paths:     btcusdtPaths,
				},
			},
		}
		resp, err := msgServer.UpdateMarketMap(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("unable to create market with paths that are not on chain tickers", func() {
		msg := &types.MsgUpdateMarketMap{
			Signer: s.authority.String(),
			CreateMarkets: []types.CreateMarket{
				{
					Ticker:    ethusdt,
					Providers: ethusdtProviders,
					Paths: types.Paths{Paths: []types.Path{
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
		resp, err := msgServer.UpdateMarketMap(s.ctx, msg)
		s.Require().Error(err)
		s.Require().Nil(resp)
	})

	s.Run("update with a new market", func() {
		msg := &types.MsgUpdateMarketMap{
			Signer: s.authority.String(),
			CreateMarkets: []types.CreateMarket{
				{
					Ticker:    ethusdt,
					Providers: ethusdtProviders,
					Paths:     ethusdtPaths,
				},
			},
		}
		resp, err := msgServer.UpdateMarketMap(s.ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		queryResp, err := qs.GetMarketMap(s.ctx, &types.GetMarketMapRequest{})
		s.Require().NoError(err)
		s.Require().Equal(queryResp.MarketMap, types.MarketMap{
			Tickers: map[string]types.Ticker{
				btcusdt.String(): btcusdt,
				usdtusd.String(): usdtusd,
				ethusdt.String(): ethusdt,
			},
			Paths: map[string]types.Paths{
				btcusdt.String(): btcusdtPaths,
				usdtusd.String(): usdtusdPaths,
				ethusdt.String(): ethusdtPaths,
			},
			Providers: map[string]types.Providers{
				btcusdt.String(): btcusdtProviders,
				usdtusd.String(): usdtusdProviders,
				ethusdt.String(): ethusdtProviders,
			},
		})
	})
}
