package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/mm2/types"
)

func TestProviderConfigValidateBasic(t *testing.T) {
	t.Run("valid config - pass", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "ticker",
			Metadata_JSON:  "",
		}
		require.NoError(t, pc.ValidateBasic())
	})
	t.Run("valid config inverted - pass", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "ticker",
			Invert:         true,
			Metadata_JSON:  "",
		}
		require.NoError(t, pc.ValidateBasic())
	})
	t.Run("valid config with index - pass", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "ticker",
			Index:          "index",
			Metadata_JSON:  "",
		}
		require.NoError(t, pc.ValidateBasic())
	})
	t.Run("invalid name - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "",
			OffChainTicker: "ticker",
			Metadata_JSON:  "",
		}
		require.Error(t, pc.ValidateBasic())
	})
	t.Run("invalid offchain ticker - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "",
			Metadata_JSON:  "",
		}
		require.Error(t, pc.ValidateBasic())
	})
	t.Run("invalid json - fail", func(t *testing.T) {
		pc := types.ProviderConfig{
			Name:           "mexc",
			OffChainTicker: "ticker",
			Metadata_JSON:  "invalid",
		}
		require.Error(t, pc.ValidateBasic())
	})
}

func TestProviderConfigEqual(t *testing.T) {
	cases := []struct {
		name  string
		pc    types.ProviderConfig
		other types.ProviderConfig
		exp   bool
	}{
		{
			name: "equal - basic",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Metadata_JSON:  "",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Metadata_JSON:  "",
			},
			exp: true,
		},
		{
			name: "equal - inverted with index",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Index:          "index",
				Invert:         true,
				Metadata_JSON:  "",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Index:          "index",
				Invert:         true,
				Metadata_JSON:  "",
			},
			exp: true,
		},
		{
			name: "equal - same metadata",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Metadata_JSON:  "{data: 1}",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Metadata_JSON:  "{data: 1}",
			},
			exp: true,
		},
		{
			name: "different name",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
			},
			other: types.ProviderConfig{
				Name:           "binance",
				OffChainTicker: "ticker",
			},
			exp: false,
		},
		{
			name: "different offchain ticker",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker2",
			},
			exp: false,
		},
		{
			name: "different invert",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Invert:         true,
				Metadata_JSON:  "",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Invert:         false,
				Metadata_JSON:  "",
			},
			exp: false,
		},
		{
			name: "different index",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Invert:         true,
				Index:          "",
				Metadata_JSON:  "",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Invert:         true,
				Index:          "index",
				Metadata_JSON:  "",
			},
			exp: false,
		},
		{
			name: "different metadata",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Metadata_JSON:  "{data: 1}",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
				Metadata_JSON:  "",
			},
			exp: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, tc.pc.Equal(tc.other))
		})
	}
}
