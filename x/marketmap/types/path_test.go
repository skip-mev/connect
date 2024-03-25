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
	}

	btcusdtPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						Provider: "kucoin",
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "BITCOIN",
							Quote: "USDT",
						},
					},
				},
			},
		},
	}

	btcusdtProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "btc-usdt",
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
	}

	usdtusdPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
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

	usdtusdProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdt-usd",
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
	}

	usdcusdPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						Provider: "kucoin",
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "USDC",
							Quote: "USD",
						},
					},
				},
			},
		},
	}

	usdcusdProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdc-usd",
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
	}

	ethusdtPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						Provider: "kucoin",
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "ETHEREUM",
							Quote: "USDT",
						},
					},
				},
			},
		},
	}

	ethusdtProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "eth-usdt",
			},
		},
	}

	usdteth = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDT",
			Quote: "ETHEREUM",
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	tickers = map[string]types.Ticker{
		btcusdt.String(): btcusdt,
		usdcusd.String(): usdcusd,
		usdtusd.String(): usdtusd,
		ethusdt.String(): ethusdt,
	}

	paths = map[string]types.Paths{
		btcusdt.String(): btcusdtPaths,
		usdcusd.String(): usdcusdPaths,
		usdtusd.String(): usdtusdPaths,
		ethusdt.String(): ethusdtPaths,
	}

	providers = map[string]types.Providers{
		btcusdt.String(): btcusdtProviders,
		usdcusd.String(): usdcusdProviders,
		usdtusd.String(): usdtusdProviders,
		ethusdt.String(): ethusdtProviders,
	}

	markets = struct {
		tickers   map[string]types.Ticker
		paths     map[string]types.Paths
		providers map[string]types.Providers
	}{
		tickers:   tickers,
		paths:     paths,
		providers: providers,
	}

	_ = markets
)

func TestPaths(t *testing.T) {
	testCases := []struct {
		name         string
		paths        types.Paths
		currencyPair slinkytypes.CurrencyPair
		expErr       bool
	}{
		{
			name:         "valid",
			paths:        btcusdtPaths,
			currencyPair: btcusdt.CurrencyPair,
			expErr:       false,
		},
		{
			name:         "invalid",
			paths:        types.Paths{},
			currencyPair: btcusdt.CurrencyPair,
			expErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.paths.ValidateBasic(tc.currencyPair)
			if tc.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestPathsEqual(t *testing.T) {
	cases := []struct {
		name  string
		paths types.Paths
		other types.Paths
		exp   bool
	}{
		{
			name:  "equal paths",
			paths: btcusdtPaths,
			other: btcusdtPaths,
			exp:   true,
		},
		{
			name:  "different length",
			paths: btcusdtPaths,
			other: types.Paths{},
			exp:   false,
		},
		{
			name:  "different paths",
			paths: btcusdtPaths,
			other: ethusdtPaths,
			exp:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.paths.Equal(tc.other), tc.exp)
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
			name: "invalid path with multiple operations inverted duplicate",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: ethusdt.CurrencyPair,
					},
					{
						CurrencyPair: usdteth.CurrencyPair,
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

func TestPathEqual(t *testing.T) {
	cases := []struct {
		name  string
		path  types.Path
		other types.Path
		exp   bool
	}{
		{
			name: "equal paths",
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
			other: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: usdtusd.CurrencyPair,
					},
				},
			},
			exp: true,
		},
		{
			name: "different length",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
				},
			},
			other: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
					{
						CurrencyPair: usdtusd.CurrencyPair,
					},
				},
			},
			exp: false,
		},
		{
			name: "different operations",
			path: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: btcusdt.CurrencyPair,
					},
				},
			},
			other: types.Path{
				Operations: []types.Operation{
					{
						CurrencyPair: usdtusd.CurrencyPair,
					},
				},
			},
			exp: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.path.Equal(tc.other), tc.exp)
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

func TestOperationEqual(t *testing.T) {
	cases := []struct {
		name  string
		op    types.Operation
		other types.Operation
		exp   bool
	}{
		{
			name: "equal operations",
			op: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   false,
			},
			other: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   false,
			},
			exp: true,
		},
		{
			name: "different base",
			op: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   false,
			},
			other: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "ETHEREUM",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   false,
			},
			exp: false,
		},
		{
			name: "different invert strategy",
			op: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   false,
			},
			other: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   true,
			},
			exp: false,
		},
		{
			name: "different provider",
			op: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "kucoin",
				Invert:   false,
			},
			other: types.Operation{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Provider: "binance",
				Invert:   false,
			},
			exp: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.op.Equal(tc.other), tc.exp)
		})
	}
}
