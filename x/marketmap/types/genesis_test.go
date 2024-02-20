package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestGenesisState(t *testing.T) {
	t.Run("good empty genesis state", func(t *testing.T) {
		gs := types.GenesisState{}
		require.NoError(t, gs.ValidateBasic())
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
					ethusdt.String(): {ethusdt.Paths},
					btcusdt.String(): {btcusdt.Paths},
					usdcusd.String(): {usdcusd.Paths},
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): {ethusdt.Providers},
					btcusdt.String(): {btcusdt.Providers},
					usdcusd.String(): {usdcusd.Providers},
				},
			},
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
					ethusdt.String(): {ethusdt.Paths},
					btcusdt.String(): {btcusdt.Paths},
					usdcusd.String(): {usdcusd.Paths},
				},
				Providers: map[string]types.Providers{
					usdtusd.String(): {usdtusd.Providers},
					btcusdt.String(): {btcusdt.Providers},
					usdcusd.String(): {usdcusd.Providers},
				},
			},
		}
		require.Error(t, gs.ValidateBasic())
	})
}
