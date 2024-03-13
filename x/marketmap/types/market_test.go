package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestMarketMapValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		marketMap types.MarketMap
		expectErr bool
	}{
		{
			name:      "valid empty",
			marketMap: types.MarketMap{},
			expectErr: false,
		},
		{
			name: "valid map",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: false,
		},
		{
			name: "invalid mismatch ticker",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					usdtusd.String(): usdtusdProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: true,
		},
		{
			name: "invalid ticker does not exist for a given provider",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: true,
		},
		{
			name: "invalid ticker string does not match ticker ID",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					"invalid":        ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.marketMap.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
