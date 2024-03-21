package coinbase_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/base/testutils"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var mogusd = mmtypes.NewTicker("MOG", "USD", 8, 1)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []mmtypes.Ticker
		url         string
		expectedErr bool
	}{
		{
			name: "valid",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
			},
			url:         "https://api.coinbase.com/v2/prices/BTC-USD/spot",
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "unknown currency",
			cps: []mmtypes.Ticker{
				mogusd,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(coinbase.Name, coinbase.DefaultProviderConfig)
			require.NoError(t, err)

			h, err := coinbase.NewAPIHandler(marketConfig, coinbase.DefaultAPIConfig)
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
			name: "valid",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
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
					constants.BITCOIN_USD: {
						Value: big.NewInt(102025000000),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "malformed response",
			cps:  []mmtypes.Ticker{constants.BITCOIN_USD},
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
					constants.BITCOIN_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse float",
			cps:  []mmtypes.Ticker{constants.BITCOIN_USD},
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
					constants.BITCOIN_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse json",
			cps:  []mmtypes.Ticker{constants.BITCOIN_USD},
			response: testutils.CreateResponseFromJSON(
				`
toms obvious but not minimal language
	`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "multiple currency pairs to parse response for",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
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
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("multiple cps"), providertypes.ErrorAPIGeneral),
					},
					constants.ETHEREUM_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("multiple cps"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(coinbase.Name, coinbase.DefaultProviderConfig)
			require.NoError(t, err)

			h, err := coinbase.NewAPIHandler(marketConfig, coinbase.DefaultAPIConfig)
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
