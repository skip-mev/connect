package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/testutil"
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
	t.Run("invalid empty name - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "",
			OffChainTicker: "ticker",
		}
		require.Error(t, pc.ValidateBasic())
	})
	t.Run("invalid empty offchain ticker - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "",
		}
		require.Error(t, pc.ValidateBasic())
	})
	t.Run("invalid too long name - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           testutil.RandomString(types.MaxProviderNameFieldLength + 1),
			OffChainTicker: "ticker",
		}
		require.Error(t, pc.ValidateBasic())
	})
	t.Run("invalid too long offchain ticker - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: testutil.RandomString(types.MaxProviderTickerFieldLength + 1),
		}
		require.Error(t, pc.ValidateBasic())
	})
}
