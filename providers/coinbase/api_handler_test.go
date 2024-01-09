package coinbase_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/skip-mev/slinky/providers/base/testutils"
	"github.com/skip-mev/slinky/providers/coinbase"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/require"
)

var config = coinbase.NewConfig(
	map[string]string{
		"BITCOIN":  "BTC",
		"USD":      "USD",
		"ETHEREUM": "ETH",
	},
)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []oracletypes.CurrencyPair
		url         string
		expectedErr bool
	}{
		{
			name: "valid",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			url:         "https://api.coinbase.com/v2/prices/BTC-USD/spot",
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "unknown base currency",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USD"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "unknown quote currency",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := coinbase.CoinBaseAPIHandler{
				Config: config,
			}

			url, err := h.CreateURL(tc.cps)
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.url, url)
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	testCases := []struct {
		name     string
		cps      []oracletypes.CurrencyPair
		response *http.Response
		expected providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
	}{
		{
			name: "valid",
			cps:  []oracletypes.CurrencyPair{oracletypes.NewCurrencyPair("BITCOIN", "USD")},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"amount": "1020.25",
		"currency": "USD"
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(102025000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
		},
		{
			name: "invalid quote currency",
			cps:  []oracletypes.CurrencyPair{oracletypes.NewCurrencyPair("BITCOIN", "USD")},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"amount": "1020.25",
		"currency": "MOG"
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("expected quote currency USD, got MOG"),
				},
			),
		},
		{
			name: "malformed response",
			cps:  []oracletypes.CurrencyPair{oracletypes.NewCurrencyPair("BITCOIN", "USD")},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"amount": "1020.25",
		"currency": "USD",
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("bad format"),
				},
			),
		},
		{
			name: "unable to parse float",
			cps:  []oracletypes.CurrencyPair{oracletypes.NewCurrencyPair("BITCOIN", "USD")},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"amount": "$1020.25",
		"currency": "USD"
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("bad format"),
				},
			),
		},
		{
			name: "unable to parse json",
			cps:  []oracletypes.CurrencyPair{oracletypes.NewCurrencyPair("BITCOIN", "USD")},
			response: testutils.CreateResponseFromJSON(
				`
toms obvious but not minimal language
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("bad format"),
				},
			),
		},
		{
			name: "multiple currency pairs to parse response for",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"amount": "1020.25",
		"currency": "USD"
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  fmt.Errorf("multiple cps"),
					oracletypes.NewCurrencyPair("ETHEREUM", "USD"): fmt.Errorf("multiple cps"),
				},
			),
		},
		{
			name: "quote currency is not supported",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"amount": "1020.25",
		"currency": "MOG"
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "MOG"): fmt.Errorf("unknown quote currency MOG"),
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := coinbase.CoinBaseAPIHandler{
				Config: config,
			}

			now := time.Now()
			resp := h.ParseResponse(tc.cps, tc.response)

			require.Len(t, resp.Resolved, len(tc.expected.Resolved))
			require.Len(t, resp.UnResolved, len(tc.expected.UnResolved))

			for cp, result := range tc.expected.Resolved {
				require.Contains(t, resp.Resolved, cp)
				r := resp.Resolved[cp]
				require.Equal(t, result.Value, r.Value)
				require.True(t, r.Timestamp.After(now))
			}

			for cp := range tc.expected.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}
