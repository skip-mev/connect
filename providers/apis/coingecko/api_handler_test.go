package coingecko_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	config = coingecko.Config{
		SupportedBases: map[string]string{
			"BITCOIN":  "bitcoin",
			"ETHEREUM": "ethereum",
		},
		SupportedQuotes: map[string]string{
			"USD":     "usd",
			"BITCOIN": "btc",
		},
	}
	apiConfig = coingecko.Config{
		APIKey: "test-key",
		SupportedBases: map[string]string{
			"BITCOIN":  "bitcoin",
			"ETHEREUM": "ethereum",
		},
		SupportedQuotes: map[string]string{
			"USD":     "usd",
			"BITCOIN": "btc",
		},
	}
)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []oracletypes.CurrencyPair
		config      coingecko.Config
		url         string
		expectedErr bool
	}{
		{
			name: "single valid currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			config:      config,
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			config:      config,
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs with multiple quotes",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			config:      config,
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd,btc&precision=18",
			expectedErr: false,
		},
		{
			name: "single valid currency pair with api key",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			config:      apiConfig,
			url:         "https://pro-api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18&x-cg-pro-api-key=test-key",
			expectedErr: false,
		},
		{
			name: "no supported bases",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USD"),
			},
			config:      config,
			url:         "",
			expectedErr: true,
		},
		{
			name: "no supported quotes",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			config:      config,
			url:         "",
			expectedErr: true,
		},
		{
			name: "some supported and non-supported currency pairs",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("MOG", "USD"),
			},
			config:      config,
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := coingecko.CoinGeckoAPIHandler{
				Config: tc.config,
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
			name: "single valid currency pair",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(102025000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
		},
		{
			name: "single valid currency pair that did not get a price response",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("currency pair BITCOIN-USD did not get a response"),
				},
			),
		},
		{
			name: "unknown base",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "unknown quote",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "unsupported base",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("MOG", "USD"): fmt.Errorf("unknown base currency MOG"),
				},
			),
		},
		{
			name: "unsupported quote",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "MOG"): fmt.Errorf("unknown quote currency MOG"),
				},
			),
		},
		{
			name: "bad response",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			response: testutils.CreateResponseFromJSON(
				`
shout out my label thats me
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "MOG"): fmt.Errorf("json error"),
				},
			),
		},
		{
			name: "bad price response",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): fmt.Errorf("invalid syntax"),
				},
			),
		},
		{
			name: "multiple bases with single quotes",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
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
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(102025000000),
					},
					oracletypes.NewCurrencyPair("ETHEREUM", "USD"): {
						Value: big.NewInt(102000000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
		},
		{
			name: "single base with multiple quotes",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				oracletypes.NewCurrencyPair("BITCOIN", "BITCOIN"),
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"bitcoin": {
		"usd": 1020.25,
		"btc": 1
	}
}
	`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): {
						Value: big.NewInt(102025000000),
					},
					oracletypes.NewCurrencyPair("BITCOIN", "BITCOIN"): {
						Value: big.NewInt(100000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := coingecko.CoinGeckoAPIHandler{
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
