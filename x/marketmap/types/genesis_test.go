package types_test

import (
	"testing"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"

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
			Tickers: types.TickersConfig{
				Tickers: []types.Ticker{
					ethusdt,
					btcusdt,
					usdcusd,
				},
			},
		}
		require.NoError(t, gs.ValidateBasic())
	})

	t.Run("bad genesis state", func(t *testing.T) {
		gs := types.GenesisState{
			Tickers: types.TickersConfig{
				Tickers: []types.Ticker{
					ethusdt,
					{
						CurrencyPair:     slinkytypes.CurrencyPair{},
						Decimals:         0,
						MinProviderCount: 0,
						Paths:            nil,
						Providers:        nil,
						Metadata_JSON:    "",
					},
				},
			},
		}
		require.Error(t, gs.ValidateBasic())
	})
}
