package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestGenesisState(t *testing.T) {
	t.Run("invalid empty genesis state - fail", func(t *testing.T) {
		gs := types.GenesisState{}
		require.Error(t, gs.ValidateBasic())
	})

	t.Run("invalid params - fail", func(t *testing.T) {
		gs := types.DefaultGenesisState()

		gs.Params.MarketAuthority = "invalid"
		require.Error(t, gs.ValidateBasic())
	})

	t.Run("good populated genesis state", func(t *testing.T) {
		gs := types.GenesisState{
			MarketMap: types.MarketMap{
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
			},
			Params: types.DefaultParams(),
		}
		require.NoError(t, gs.ValidateBasic())
	})

	t.Run("bad genesis state - mistmatch", func(t *testing.T) {
		gs := types.GenesisState{
			MarketMap: types.MarketMap{
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
					usdtusd.String(): usdtusdProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
		}
		require.Error(t, gs.ValidateBasic())
	})
}
