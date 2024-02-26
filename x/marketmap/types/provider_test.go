package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestProviderConfigValidateBasic(t *testing.T) {
	t.Run("valid config - pass", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "ticker",
		}
		require.NoError(t, pc.ValidateBasic())
	})
	t.Run("invalid name - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "",
			OffChainTicker: "ticker",
		}
		require.Error(t, pc.ValidateBasic())
	})
	t.Run("valid offchain ticker - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "",
		}
		require.Error(t, pc.ValidateBasic())
	})
}
