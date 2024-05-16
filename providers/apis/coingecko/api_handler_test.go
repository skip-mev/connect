package coingecko_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base/testutils"
)

var (
	btcusd = types.DefaultProviderTicker{
		OffChainTicker: "btc/usd",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "eth/usd",
	}
	ethbtc = types.DefaultProviderTicker{
		OffChainTicker: "eth/btc",
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
			name: "single valid currency pair",
			cps: []types.ProviderTicker{
				btcusd,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			url:         "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd&precision=18",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs with multiple quotes",
			cps: []types.ProviderTicker{
				btcusd,
				ethusd,
				ethbtc,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(1020.25),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "single valid currency pair that did not get a price response",
			cps: []types.ProviderTicker{
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("currency pair BITCOIN-USD did not get a response"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "bad response",
			cps: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(
				`
shout out my label that's me
	`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("json error"), providertypes.ErrorWebSocketGeneral),
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
					btcusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid syntax"), providertypes.ErrorWebSocketGeneral),
					},
				},
			),
		},
		{
			name: "multiple bases with single quotes",
			cps: []types.ProviderTicker{
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(1020.25),
					},
					ethusd: {
						Value: big.NewFloat(1020),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "single base with multiple quotes",
			cps: []types.ProviderTicker{
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					ethusd: {
						Value: big.NewFloat(1020.25),
					},
					ethbtc: {
						Value: big.NewFloat(1),
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
