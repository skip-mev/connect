package types_test

import (
	"testing"

	"github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {
	t.Run("good empty genesis state", func(t *testing.T) {
		gs := types.GenesisState{}
		require.NoError(t, gs.ValidateBasic())
	})

	t.Run("bad genesis state", func(t *testing.T) {
		gs := types.GenesisState{
			Config: types.AggregateMarketConfig{
				TickerConfigs: map[string]types.PathsConfig{
					"BITCOIN/USDT": {},
				},
			},
		}
		require.Error(t, gs.ValidateBasic())
	})
}
