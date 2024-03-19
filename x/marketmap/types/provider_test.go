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

func TestProvidersEqual(t *testing.T) {
	cases := []struct {
		name  string
		p     types.Providers
		other types.Providers
		exp   bool
	}{
		{
			name: "equal",
			p: types.Providers{
				Providers: []types.ProviderConfig{
					{
						Name:           "mexc",
						OffChainTicker: "ticker",
					},
				},
			},
			other: types.Providers{
				Providers: []types.ProviderConfig{
					{
						Name:           "mexc",
						OffChainTicker: "ticker",
					},
				},
			},
			exp: true,
		},
		{
			name: "different length",
			p: types.Providers{
				Providers: []types.ProviderConfig{},
			},
			other: types.Providers{
				Providers: []types.ProviderConfig{
					{
						Name:           "mexc",
						OffChainTicker: "ticker",
					},
				},
			},
			exp: false,
		},
		{
			name: "different provider",
			p: types.Providers{
				Providers: []types.ProviderConfig{
					{
						Name:           "mexc",
						OffChainTicker: "ticker",
					},
				},
			},
			other: types.Providers{
				Providers: []types.ProviderConfig{
					{
						Name:           "binance",
						OffChainTicker: "ticker",
					},
				},
			},
			exp: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, tc.p.Equal(tc.other))
		})
	}
}

func TestProviderConfigEqual(t *testing.T) {
	cases := []struct {
		name  string
		pc    types.ProviderConfig
		other types.ProviderConfig
		exp   bool
	}{
		{
			name: "equal",
			pc: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
			},
			other: types.ProviderConfig{
				Name:           "mexc",
				OffChainTicker: "ticker",
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, tc.pc.Equal(tc.other))
		})
	}
}
