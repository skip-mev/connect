package coingecko_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

var (
	mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)
	btcmog = mmtypes.NewTicker("BTC", "MOG", 8, 1)
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
				constants.BITCOIN_USD,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs with multiple quotes",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
				constants.ETHEREUM_BITCOIN,
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
				constants.BITCOIN_USD,
				mogusd,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(coingecko.Name, coingecko.DefaultMarketConfig)
			require.NoError(t, err)

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
		expected types.PriceResponse
	}{
		{
			name: "single valid currency pair",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USD: {
						Value: big.NewInt(102025000000),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "single valid currency pair that did not get a price response",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("currency pair BITCOIN-USD did not get a response"), providertypes.ErrorWebSocketGeneral),
					},
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcmog: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("json error"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "bad price response",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid syntax"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "multiple bases with single quotes",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USD: {
						Value: big.NewInt(102025000000),
					},
					constants.ETHEREUM_USD: {
						Value: big.NewInt(102000000000),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "single base with multiple quotes",
			cps: []mmtypes.Ticker{
				constants.ETHEREUM_USD,
				constants.ETHEREUM_BITCOIN,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.ETHEREUM_USD: {
						Value: big.NewInt(102025000000),
					},
					constants.ETHEREUM_BITCOIN: {
						Value: big.NewInt(100000000),
					},
				},
				types.UnResolvedPrices{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(coingecko.Name, coingecko.DefaultMarketConfig)
			require.NoError(t, err)

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
