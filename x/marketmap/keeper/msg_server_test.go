package keeper_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestCreateMarket() {
	msgServer := keeper.NewMsgServer(s.keeper)

	// create an initial market
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
	s.Require().NoError(err)
	s.Require().NotNil(resp)

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
}
