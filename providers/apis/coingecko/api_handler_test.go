package coingecko_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base/testutils"
)

var (
	btc_usd = coingecko.DefaultMarketConfig.MustGetProviderTicker(constants.BITCOIN_USD)
	eth_usd = coingecko.DefaultMarketConfig.MustGetProviderTicker(constants.ETHEREUM_USD)
	eth_btc = coingecko.DefaultMarketConfig.MustGetProviderTicker(constants.ETHEREUM_BITCOIN)
)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []types.ProviderTicker
		url         string
		expectedErr bool
	}{
		{
			name: "single valid currency pair",
			cps: []types.ProviderTicker{
				btc_usd,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			cps: []types.ProviderTicker{
				btc_usd,
				eth_usd,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs with multiple quotes",
			cps: []types.ProviderTicker{
				btc_usd,
				eth_usd,
				eth_btc,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd,btc&precision=18",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coingecko.NewAPIHandler(coingecko.DefaultAPIConfig)
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
			name: "single valid currency pair",
			cps: []types.ProviderTicker{
				btc_usd,
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
					btc_usd: {
						Value: math.Float64ToBigFloat(1020.25, types.DefaultTickerDecimals),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "single valid currency pair that did not get a price response",
			cps: []types.ProviderTicker{
				btc_usd,
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
					btc_usd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("currency pair BITCOIN-USD did not get a response"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "bad response",
			cps: []types.ProviderTicker{
				btc_usd,
			},
			response: testutils.CreateResponseFromJSON(
				`
shout out my label thats me
	`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btc_usd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("json error"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "bad price response",
			cps: []types.ProviderTicker{
				btc_usd,
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
					btc_usd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid syntax"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "multiple bases with single quotes",
			cps: []types.ProviderTicker{
				btc_usd,
				eth_usd,
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
					btc_usd: {
						Value: math.Float64ToBigFloat(1020.25, types.DefaultTickerDecimals),
					},
					eth_usd: {
						Value: math.Float64ToBigFloat(1020, types.DefaultTickerDecimals),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "single base with multiple quotes",
			cps: []types.ProviderTicker{
				eth_usd,
				eth_btc,
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
					eth_usd: {
						Value: math.Float64ToBigFloat(1020.25, types.DefaultTickerDecimals),
					},
					eth_btc: {
						Value: math.Float64ToBigFloat(1, types.DefaultTickerDecimals),
					},
				},
				types.UnResolvedPrices{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coingecko.NewAPIHandler(coingecko.DefaultAPIConfig)
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
				require.Equal(t, result.Value.SetPrec(types.DefaultTickerDecimals), r.Value.SetPrec(types.DefaultTickerDecimals))
				require.True(t, r.Timestamp.After(now))
			}

			for cp := range tc.expected.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}
