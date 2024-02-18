package types_test

import (
	"testing"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	btcusdt = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "BITCOIN",
			Quote: "USDT",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				[]types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "BITCOIN",
							Quote: "USDT",
						},
					},
				},
			},
		},
	}

	usdtusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDT",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				[]types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
		},
	}

	usdcusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDC",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				[]types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "USDC",
							Quote: "USD",
						},
					},
				},
			},
		},
	}

	ethusdt = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "ETHEREUM",
			Quote: "USDT",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				[]types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "ETHEREUM",
							Quote: "USDT",
						},
					},
				},
			},
		},
	}
)

func TestPath(t *testing.T) {
	testCases := []struct {
		name   string
		path   types.Path
		target string
		expErr bool
	}{
		{
			name:   "empty path",
			path:   types.Path{},
			target: "",
			expErr: true,
		},
		{
			name: "invalid path with a single operation",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{},
					},
				},
			},
			target: "",
			expErr: true,
		},
		{
			name: "invalid path with multiple operations with a bad ticker in the route",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: slinkytypes.CurrencyPair{},
					},
				},
			},
			target: "",
			expErr: true,
		},
		{
			name: "invalid path with multiple operations and mismatching tickers",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: ethusdt.CurrencyPair,
					},
				},
			},
			target: "",
			expErr: true,
		},
		{
			name: "invalid path with multiple operations and cyclic graph",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: usdtusd.CurrencyPair,
					},
					{
						CurrencyPair: usdtusd.CurrencyPair,
						Invert:       true,
					},
				},
			},
			target: "",
			expErr: true,
		},
		{
			name: "valid path with a single operation",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
				},
			},
			target: "BITCOIN/USDT",
			expErr: false,
		},
		{
			name: "valid path with multiple operations",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: usdtusd.CurrencyPair,
					},
				},
			},
			target: "BITCOIN/USD",
			expErr: false,
		},
		{
			name: "valid path with multiple operations and inverted tickers",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: usdtusd.CurrencyPair,
						Invert:       true,
					},
					{
						CurrencyPair: btcusdt.CurrencyPair,
						Invert:       true,
					},
				},
			},
			target: "USD/BITCOIN",
			expErr: false,
		},
		{
			name: "valid path with multiple operations and inverted tickers",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: usdtusd.CurrencyPair,
					},
					{
						CurrencyPair: usdcusd.CurrencyPair,
						Invert:       true,
					},
				},
			},
			target: "BITCOIN/USDC",
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.path.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, tc.path.Match(tc.target))
			}
		})
	}
}

func TestOperation(t *testing.T) {
	t.Run("valid operation", func(t *testing.T) {
		cp := slinkytypes.CurrencyPair{
			Base:  "BITCOIN",
			Quote: "USDT",
		}

		_, err := types.NewOperation(cp, false)
		require.NoError(t, err)
	})

	t.Run("invalid operation", func(t *testing.T) {
		cp := slinkytypes.CurrencyPair{}
		_, err := types.NewOperation(cp, false)
		require.Error(t, err)
	})
}
