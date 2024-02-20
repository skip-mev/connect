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
								CurrencyPair: btcusdt.CurrencyPair,
							},
						},
					},
				},
				Providers: []types.ProviderConfig{
					{
						Name:           "binance",
						OffChainTicker: "btc-usd",
					},
					{
						Name:           "kucoin",
						OffChainTicker: "btcusd",
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
					Quote: "USD",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Paths: []types.Path{
					{
						Operations: []types.Operation{
							{
								CurrencyPair: btcusdt.CurrencyPair,
							},
							{
								CurrencyPair: usdtusd.CurrencyPair,
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
								CurrencyPair: btcusdt.CurrencyPair,
							},
							{
								CurrencyPair: ethusdt.CurrencyPair,
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
