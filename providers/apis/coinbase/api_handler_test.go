package coinbase_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/base/testutils"
)

var (
	mogusd = types.DefaultProviderTicker{
		OffChainTicker: "MOGUSD",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSD",
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
			name: "valid",
			cps: []types.ProviderTicker{
				mogusd,
			},
			url:         "https://api.coinbase.com/v2/prices/BTC-USD/spot",
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []types.ProviderTicker{
				mogusd,
				ethusd,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coinbase.NewAPIHandler(coinbase.DefaultAPIConfig)
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
			name: "valid",
			cps: []types.ProviderTicker{
				mogusd,
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					mogusd: {
						Value: big.NewFloat(1020.25),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "malformed response",
			cps: []types.ProviderTicker{
				mogusd,
			},
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse float",
			cps: []types.ProviderTicker{
				mogusd,
			},
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
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse json",
			cps:  []types.ProviderTicker{mogusd},
			response: testutils.CreateResponseFromJSON(
				`
toms obvious but not minimal language
	`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coinbase.NewAPIHandler(coinbase.DefaultAPIConfig)
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
