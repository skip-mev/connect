package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
	}

	usdtusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDT",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	usdcusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDC",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	ethusdt = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "ETHEREUM",
			Quote: "USDT",
		},
		Decimals:         8,
		MinProviderCount: 1,
	}
)

func TestPathsConfig(t *testing.T) {
	testCases := []struct {
		name   string
		config types.PathsConfig
		expErr bool
	}{
		{
			name: "valid paths config",
			config: types.PathsConfig{
				Ticker: btcusdt,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								Ticker: btcusdt,
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name: "empty paths",
			config: types.PathsConfig{
				Ticker: btcusdt,
				Paths:  []types.Path{},
			},
			expErr: true,
		},
		{
			name: "invalid ticker",
			config: types.PathsConfig{
				Ticker: types.Ticker{},
				Paths:  []types.Path{},
			},
			expErr: true,
		},
		{
			name: "invalid path",
			config: types.PathsConfig{
				Ticker: btcusdt,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								Ticker: types.Ticker{},
							},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "invalid path with mismatching ticker",
			config: types.PathsConfig{
				Ticker: btcusdt,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								Ticker: ethusdt,
							},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "duplicate path",
			config: types.PathsConfig{
				Ticker: btcusdt,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								Ticker: btcusdt,
							},
						},
					},
					{
						Operations: []types.Operation{
							{
								Ticker: btcusdt,
							},
						},
					},
				},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

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
						Ticker: types.Ticker{},
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
						Ticker: btcusdt,
					},
					{
						Ticker: types.Ticker{},
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
						Ticker: btcusdt,
					},
					{
						Ticker: ethusdt,
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
						Ticker: btcusdt,
					},
					{
						Ticker: usdtusd,
					},
					{
						Ticker: usdtusd,
						Invert: true,
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
						Ticker: btcusdt,
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
						Ticker: btcusdt,
					},
					{
						Ticker: usdtusd,
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
						Ticker: usdtusd,
						Invert: true,
					},
					{
						Ticker: btcusdt,
						Invert: true,
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
						Ticker: btcusdt,
					},
					{
						Ticker: usdtusd,
					},
					{
						Ticker: usdcusd,
						Invert: true,
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
		ticker := types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "BITCOIN",
				Quote: "USDT",
			},
			Decimals:         8,
			MinProviderCount: 1,
		}

		_, err := types.NewOperation(ticker, false)
		require.NoError(t, err)
	})

	t.Run("invalid operation", func(t *testing.T) {
		ticker := types.Ticker{}
		_, err := types.NewOperation(ticker, false)
		require.Error(t, err)
	})
}
