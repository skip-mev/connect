package types_test

import (
	"testing"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestTicker(t *testing.T) {
	testCases := []struct {
		name   string
		ticker types.Ticker
		expErr bool
	}{
		{
			name: "valid ticker",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
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
			name: "valid ticker multiple paths",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								Ticker: btcusdt,
							},
							{
								Ticker: usdtusd,
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name: "invalid paths",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								Ticker: btcusdt,
							},
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
			name: "empty base",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			expErr: true,
		},
		{
			name: "empty quote",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			expErr: true,
		},
		{
			name: "invalid base",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "bitcoin",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			expErr: true,
		},
		{
			name: "invalid quote",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "usdt",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			expErr: true,
		},
		{
			name: "invalid decimals",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         0,
				MinProviderCount: 1,
			},
			expErr: true,
		},
		{
			name: "invalid min provider count",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 0,
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ticker.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestTickerConfig(t *testing.T) {
	testCases := []struct {
		name         string
		tickerConfig types.TickerConfig
		expErr       bool
	}{
		{
			name: "valid ticker config",
			tickerConfig: types.TickerConfig{
				Ticker: types.Ticker{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "BITCOIN",
						Quote: "USDT",
					},
					Decimals:         8,
					MinProviderCount: 1,
				},
				OffChainTicker: "BTC/USDT",
			},
			expErr: false,
		},
		{
			name: "invalid ticker",
			tickerConfig: types.TickerConfig{
				Ticker: types.Ticker{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "",
						Quote: "USDT",
					},
					Decimals:         8,
					MinProviderCount: 1,
				},
				OffChainTicker: "BTC/USDT",
			},
			expErr: true,
		},
		{
			name: "invalid off chain ticker",
			tickerConfig: types.TickerConfig{
				Ticker: types.Ticker{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "BITCOIN",
						Quote: "USDT",
					},
					Decimals:         8,
					MinProviderCount: 1,
				},
				OffChainTicker: "",
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.tickerConfig.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
