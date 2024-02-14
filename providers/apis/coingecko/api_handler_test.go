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
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	btcusd = mmtypes.NewTicker("BITCOIN", "USD", 8, 1)
	ethusd = mmtypes.NewTicker("ETHEREUM", "USD", 8, 1)
	ethbtc = mmtypes.NewTicker("ETHEREUM", "BITCOIN", 8, 1)
	mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)
	btcmog = mmtypes.NewTicker("BITCOIN", "MOG", 8, 1)

	marketConfig = mmtypes.MarketConfig{
		Name: coingecko.Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"BITCOIN/USD": {
				Ticker:         btcusd,
				OffChainTicker: "bitcoin/usd",
			},
			"ETHEREUM/USD": {
				Ticker:         ethusd,
				OffChainTicker: "ethereum/usd",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         ethbtc,
				OffChainTicker: "ethereum/btc",
			},
		},
	}
)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []mmtypes.Ticker
		url         string
		expectedErr bool
	}{
		{
			name: "single valid currency pair",
			cps: []mmtypes.Ticker{
				btcusd,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			cps: []mmtypes.Ticker{
				btcusd,
				ethusd,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs with multiple quotes",
			cps: []mmtypes.Ticker{
				btcusd,
				ethusd,
				ethbtc,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd,btc&precision=18",
			expectedErr: false,
		},
		{
			name: "no supported bases",
			cps: []mmtypes.Ticker{
				mogusd,
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "no supported quotes",
			cps: []mmtypes.Ticker{
				btcmog,
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "some supported and non-supported currency pairs",
			cps: []mmtypes.Ticker{
				btcusd,
				mogusd,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coingecko.NewAPIHandler(marketConfig, coingecko.DefaultAPIConfig)
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
		cps      []mmtypes.Ticker
		response *http.Response
		expected providertypes.GetResponse[mmtypes.Ticker, *big.Int]
	}{
		{
			name: "single valid currency pair",
			cps: []mmtypes.Ticker{
				btcusd,
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
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(102025000000),
					},
				},
				map[mmtypes.Ticker]error{},
			),
		},
		{
			name: "single valid currency pair that did not get a price response",
			cps: []mmtypes.Ticker{
				btcusd,
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
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					btcusd: fmt.Errorf("currency pair BITCOIN-USD did not get a response"),
				},
			),
		},
		{
			name: "bad response",
			cps: []mmtypes.Ticker{
				
				btcmog,
			},
			response: testutils.CreateResponseFromJSON(
				`
shout out my label thats me
	`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					btcmog: fmt.Errorf("json error"),
				},
			),
		},
		{
			name: "bad price response",
			cps: []mmtypes.Ticker{
				btcusd,
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
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					btcusd: fmt.Errorf("invalid syntax"),
				},
			),
		},
		{
			name: "multiple bases with single quotes",
			cps: []mmtypes.Ticker{
				btcusd,
				ethusd,
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
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(102025000000),
					},
					ethusd: {
						Value: big.NewInt(102000000000),
					},
				},
				map[mmtypes.Ticker]error{},
			),
		},
		{
			name: "single base with multiple quotes",
			cps: []mmtypes.Ticker{
				ethusd,
				ethbtc,
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
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{
					ethusd: {
						Value: big.NewInt(102025000000),
					},
					ethbtc: {
						Value: big.NewInt(100000000),
					},
				},
				map[mmtypes.Ticker]error{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coingecko.NewAPIHandler(marketConfig, coingecko.DefaultAPIConfig)
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
