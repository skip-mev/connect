package types_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
			name: "valid map no enabled tickers",
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
				EnabledTickers: nil,
			},
			expectErr: false,
		},
		{
			name: "valid map all enabled tickers",
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
				EnabledTickers: []string{ethusdt.String(), btcusdt.String(), usdcusd.String()},
			},
			expectErr: false,
		},
		{
			name: "invalid too many invalid tickers",
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
				EnabledTickers: []string{ethusdt.String(), usdtusd.String(), btcusdt.String(), btcusdt.String()},
			},
			expectErr: true,
		},
		{
			name: "invalid duplicate enabled tickers",
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
				EnabledTickers: []string{btcusdt.String(), btcusdt.String()},
			},
			expectErr: true,
		},
		{
			name: "invalid enabled ticker",
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
				EnabledTickers: []string{ethusdt.String(), usdtusd.String(), usdcusd.String()},
			},
			expectErr: true,
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
				EnabledTickers: []string{ethusdt.String(), btcusdt.String(), usdcusd.String()},
			},
			expectErr: true,
		},
		{
			name: "invalid ticker does not exist for a given path",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
				EnabledTickers: []string{btcusdt.String(), usdcusd.String()},
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
				EnabledTickers: []string{btcusdt.String(), usdcusd.String()},
			},
			expectErr: true,
		},
		{
			name: "invalid no providers for ticker",
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
					usdcusd.String(): usdcusdProviders,
				},
				EnabledTickers: []string{ethusdt.String(), btcusdt.String(), usdcusd.String()},
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
				EnabledTickers: []string{ethusdt.String(), btcusdt.String(), usdcusd.String()},
			},
			expectErr: true,
		},
		{
			name: "invalid path ticker not found in tickers map",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: slinkytypes.CurrencyPair{
											Base:  "ETHEREUM",
											Quote: "MOG",
										},
									},
									{
										CurrencyPair: slinkytypes.CurrencyPair{
											Base:  "MOG",
											Quote: "USDT",
										},
									},
								},
							},
						},
					},
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
				EnabledTickers: []string{ethusdt.String(), btcusdt.String(), usdcusd.String()},
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
