package bitstamp_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/bitstamp"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var (
	btcusd = types.DefaultProviderTicker{
		OffChainTicker: "BTC/USD",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "ETH/USD",
	}
)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []types.ProviderTicker
		url         string
		expectedErr bool
	}{
		{
			name:        "empty",
			cps:         []types.ProviderTicker{},
			url:         "",
			expectedErr: true,
		},
		{
			name: "valid single",
			cps: []types.ProviderTicker{
				btcusd,
			},
			url:         bitstamp.DefaultAPIConfig.Endpoints[0].URL,
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			url:         bitstamp.DefaultAPIConfig.Endpoints[0].URL,
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := bitstamp.NewAPIHandler(bitstamp.DefaultAPIConfig)
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
		cps      []types.ProviderTicker
		response *http.Response
		expected types.PriceResponse
	}{
		{
			name: "valid single",
			cps: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(
				`
[
	{
		"ask": "2211.00",
		"bid": "2188.97",
		"high": "2811.00",
		"last": "2211.00",
		"low": "2188.97",
		"open": "2211.00",
		"open_24": "2211.00",
		"pair": "BTC/USD",
		"percent_change_24": "13.57",
		"side": "0",
		"timestamp": "1643640186",
		"volume": "213.26801100",
		"vwap": "2189.80"
	}
]
				`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(2211.00),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "valid multiple",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			response: testutils.CreateResponseFromJSON(
				`
[
	{
		"ask": "2211.00",
		"bid": "2188.97",
		"high": "2811.00",
		"last": "2211.00",
		"low": "2188.97",
		"open": "2211.00",
		"open_24": "2211.00",
		"pair": "BTC/USD",
		"percent_change_24": "13.57",
		"side": "0",
		"timestamp": "1643640186",
		"volume": "213.26801100",
		"vwap": "2189.80"
	},
	{
		"ask": "2211.00",
		"bid": "2188.97",
		"high": "2811.00",
		"last": "420.69",
		"low": "2188.97",
		"open": "2211.00",
		"open_24": "2211.00",
		"pair": "ETH/USD",
		"percent_change_24": "13.57",
		"side": "0",
		"timestamp": "1643640186",
		"volume": "213.26801100",
		"vwap": "2189.80"
	}
]
				`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(2211.00),
					},
					ethusd: {
						Value: big.NewFloat(420.69),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "bad response",
			cps: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(
				`shout out my label that's me`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "bad price response",
			cps: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(
				`
[
	{
		"ask": "2211.00",
		"bid": "2188.97",
		"high": "2811.00",
		"last": "$2211.00",
		"low": "2188.97",
		"open": "2211.00",
		"open_24": "2211.00",
		"pair": "BTC/USD",
		"percent_change_24": "13.57",
		"side": "0",
		"timestamp": "1643640186",
		"volume": "213.26801100",
		"vwap": "2189.80"
	},
]
				`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid syntax"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "no response",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			response: testutils.CreateResponseFromJSON(
				`[]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
					ethusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := bitstamp.NewAPIHandler(bitstamp.DefaultAPIConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that createURL is executed before ParseResponse.
			_, err = h.CreateURL(tc.cps)
			require.NoError(t, err)

			now := time.Now()
			resp := h.ParseResponse(tc.cps, tc.response)

			require.Len(t, resp.Resolved, len(tc.expected.Resolved))
			require.Len(t, resp.UnResolved, len(tc.expected.UnResolved))

			for cp, result := range tc.expected.Resolved {
				require.Contains(t, resp.Resolved, cp)
				r := resp.Resolved[cp]
				require.Equal(t, result.Value.SetPrec(18), r.Value.SetPrec(18))
				require.True(t, r.Timestamp.After(now))
			}

			for cp := range tc.expected.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}
