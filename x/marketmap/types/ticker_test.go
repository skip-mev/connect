package types_test

import (
	"testing"

	"github.com/skip-mev/slinky/testutil"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
			},
			expErr: false,
		},
		{
			name: "invalid metadata length",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Metadata_JSON:    testutil.RandomString(types.MaxMetadataJSONFieldLength + 1),
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

func TestTickerEqual(t *testing.T) {
	cases := []struct {
		name   string
		ticker types.Ticker
		other  types.Ticker
		exp    bool
	}{
		{
			name: "equal tickers",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			exp: true,
		},
		{
			name: "different base",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "ETHEREUM",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			exp: false,
		},
		{
			name: "different quote",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "BTC",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			exp: false,
		},
		{
			name: "different decimals",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         18,
				MinProviderCount: 1,
			},
			exp: false,
		},
		{
			name: "different min provider count",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 2,
			},
			exp: false,
		},

		{
			name: "different metadata",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Metadata_JSON:    `{"key": "value"}`,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Metadata_JSON:    `{"key": "value2"}`,
			},
			exp: false,
		},

		{
			name: "different enabled",
			ticker: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Enabled:          true,
			},
			other: types.Ticker{
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "BITCOIN",
					Quote: "USDT",
				},
				Decimals:         8,
				MinProviderCount: 1,
				Enabled:          false,
			},
			exp: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.ticker.Equal(tc.other), tc.exp)
		})
	}
}
