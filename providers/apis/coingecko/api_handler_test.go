package coingecko_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var providerCfg = config.ProviderConfig{
	Name: coingecko.Name,
	API:  coingecko.DefaultAPIConfig,
	Market: config.MarketConfig{
		Name: coingecko.Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "bitcoin/usd",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ethereum/usd",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ethereum/btc",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
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
			name: "single valid currency pair",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs with multiple quotes",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd,btc&precision=18",
			expectedErr: false,
		},
		{
			name: "no supported bases",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("MOG", "USD"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "no supported quotes",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "some supported and non-supported currency pairs",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("MOG", "USD"),
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coingecko.NewAPIHandler(providerCfg)
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
			name: "single valid currency pair",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"bitcoin": {
		"usd": 1020.25
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
			name: "single valid currency pair that did not get a price response",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"bitcoin": {
		"btc" : 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("currency pair BITCOIN-USD did not get a response"),
				},
			),
		},
		{
			name: "unknown base",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"mog": {
		"usd": 1020.25,
		"btc": 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "unknown quote",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"bitcoin": {
		"mog": 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "unsupported base",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("MOG", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"mog": {
		"usd": 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{},
			),
		},
		{
			name: "unsupported quote",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"mog": {
		"usd": 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{},
			),
		},
		{
			name: "bad response",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			response: testutils.CreateResponseFromJSON(
				`
shout out my label thats me
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "MOG"): fmt.Errorf("json error"),
				},
			),
		},
		{
			name: "bad price response",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"bitcoin": {
		"usd": "$1020.25"
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[slinkytypes.CurrencyPair]error{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("invalid syntax"),
				},
			),
		},
		{
			name: "multiple bases with single quotes",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"bitcoin": {
		"usd": 1020.25
	},
	"ethereum": {
		"usd": 1020
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					slinkytypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(102025000000),
					},
					slinkytypes.NewCurrencyPair("ETHEREUM", "USD"): {
						Value: big.NewInt(102000000000),
					},
				},
				map[slinkytypes.CurrencyPair]error{},
			),
		},
		{
			name: "single base with multiple quotes",
			cps: []slinkytypes.CurrencyPair{
				slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
				slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"ethereum": {
		"usd": 1020.25,
		"btc": 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					slinkytypes.NewCurrencyPair("ETHEREUM", "USD"): {
						Value: big.NewInt(102025000000),
					},
					slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"): {
						Value: big.NewInt(100000000),
					},
				},
				map[slinkytypes.CurrencyPair]error{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coingecko.NewAPIHandler(providerCfg)
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
