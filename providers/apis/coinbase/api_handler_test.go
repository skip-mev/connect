package coinbase_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var providerCfg = config.ProviderConfig{
	Name: coinbase.Name,
	API:  coinbase.DefaultAPIConfig,
	Market: config.MarketConfig{
		Name: coinbase.Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTC-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH-USD",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
		},
	},
}

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []slinkytypes.CurrencyPair
		url         string
		expectedErr bool
	}{
		{
			name: "valid",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			url:         "https://api.coinbase.com/v2/prices/BTC-USD/spot",
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "unknown base currency",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("MOG", "USD"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "unknown quote currency",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coinbase.NewAPIHandler(providerCfg)
			require.NoError(t, err)

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
		cps      []slinkytypes.CurrencyPair
		response *http.Response
		expected providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
	}{
		{
			name: "valid",
			cps:  []slinkytypes.CurrencyPair{slinkytypes.NewCurrencyPair("BITCOIN", "USD")},
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
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(102025000000),
					},
				},
				map[slinkytypes.CurrencyPair]error{},
			),
		},
		{
			name: "malformed response",
			cps:  []slinkytypes.CurrencyPair{slinkytypes.NewCurrencyPair("BITCOIN", "USD")},
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
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("bad format"),
				},
			),
		},
		{
			name: "unable to parse float",
			cps:  []slinkytypes.CurrencyPair{slinkytypes.NewCurrencyPair("BITCOIN", "USD")},
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
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("bad format"),
				},
			),
		},
		{
			name: "unable to parse json",
			cps:  []slinkytypes.CurrencyPair{slinkytypes.NewCurrencyPair("BITCOIN", "USD")},
			response: testutils.CreateResponseFromJSON(
				`
toms obvious but not minimal language
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("bad format"),
				},
			),
		},
		{
			name: "multiple currency pairs to parse response for",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
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
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"):  fmt.Errorf("multiple cps"),
					slinkytypes.NewCurrencyPair("ETHEREUM", "USD"): fmt.Errorf("multiple cps"),
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coinbase.NewAPIHandler(providerCfg)
			require.NoError(t, err)

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
