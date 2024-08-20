package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestGenesisState(t *testing.T) {
	t.Run("invalid empty genesis state - fail", func(t *testing.T) {
		gs := types.GenesisState{}
		require.Error(t, gs.ValidateBasic())
	})

	t.Run("invalid params - fail", func(t *testing.T) {
		gs := types.DefaultGenesisState()

		gs.Params.MarketAuthorities = []string{"invalid"}
		require.Error(t, gs.ValidateBasic())
	})

	t.Run("good populated genesis state", func(t *testing.T) {
		gs := types.GenesisState{
			MarketMap: types.MarketMap{
				Markets: markets,
			},
			Params: types.DefaultParams(),
		}
		require.NoError(t, gs.ValidateBasic())
	})
}
