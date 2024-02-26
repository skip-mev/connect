package keeper_test

import (
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func (s *KeeperTestSuite) TestInitGenesisInvalidGenesis() {
	s.Run("test that init genesis with invalid genesis params - fails", func() {
		// create a fake genesis state with invalid params
		gs := types.GenesisState{
			MarketMap: types.DefaultGenesisState().MarketMap,
			Params: types.Params{
				MarketAuthority: "invalid",
				Version:         0,
			},
		}

		// assert that InitGenesis panics
		s.Require().Panics(func() {
			s.keeper.InitGenesis(s.ctx, gs)
		})
	})

	s.Run("test that init genesis with invalid lastupdated - fails", func() {
		// create a fake genesis state with invalid lastupdated
		gs := types.DefaultGenesisState()
		gs.LastUpdated = -1

		// assert that InitGenesis panics
		s.Require().Panics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})
	})

	s.Run("test that init genesis with invalid duplicate runs - fails", func() {
		// create a valid genesis
		gs := types.DefaultGenesisState()

		gs.MarketMap = types.MarketMap{
			Tickers: map[string]types.Ticker{
				ethusdt.String(): ethusdt,
				btcusdt.String(): btcusdt,
				usdcusd.String(): usdcusd,
			},
			Paths: map[string]types.Paths{
				ethusdt.String(): ethusdtPaths,
				btcusdt.String(): btcusdtPaths,
				usdcusd.String(): usdcusdPaths,
			},
			Providers: map[string]types.Providers{
				ethusdt.String(): ethusdtProviders,
				btcusdt.String(): btcusdtProviders,
				usdcusd.String(): usdcusdProviders,
			},
		}

		// assert that InitGenesis panics
		s.Require().Panics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
			s.keeper.InitGenesis(s.ctx, *gs)
		})
	})
}

func (s *KeeperTestSuite) TestInitGenesisValid() {
	s.Run("init valid default genesis", func() {
		gs := types.DefaultGenesisState()

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})
	})

	s.Run("init valid genesis with fields", func() {
		gs := types.DefaultGenesisState()
		gs.MarketMap = types.MarketMap{
			Tickers: map[string]types.Ticker{
				ethusdt.String(): ethusdt,
				btcusdt.String(): btcusdt,
				usdcusd.String(): usdcusd,
			},
			Paths: map[string]types.Paths{
				ethusdt.String(): ethusdtPaths,
				btcusdt.String(): btcusdtPaths,
				usdcusd.String(): usdcusdPaths,
			},
			Providers: map[string]types.Providers{
				ethusdt.String(): ethusdtProviders,
				btcusdt.String(): btcusdtProviders,
				usdcusd.String(): usdcusdProviders,
			},
		}

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})

		gotTickers, err := s.keeper.GetAllTickersMap(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(gs.MarketMap.Tickers, gotTickers)

		gotPaths, err := s.keeper.GetAllPathsMap(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(gs.MarketMap.Paths, gotPaths)

		gotProviders, err := s.keeper.GetAllProvidersMap(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(gs.MarketMap.Providers, gotProviders)

	})
}
