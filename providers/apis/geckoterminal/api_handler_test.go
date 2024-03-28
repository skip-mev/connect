package geckoterminal_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/geckoterminal"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

var popcat = mmtypes.NewTicker("POPCAT", "USD", 8, 1)

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
				constants.MOG_USD,
			},
			url:         "https://api.geckoterminal.com/api/v2/simple/networks/eth/token_price/0xaaee1a9723aadb7afa2810263653a34ba2c21c7a",
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USD,
				constants.ETHEREUM_USD,
			},
			url:         "https://api.geckoterminal.com/api/v2/simple/networks/eth/token_price/0xaaee1a9723aadb7afa2810263653a34ba2c21c7a,0x6982508145454Ce325dDbE47a25d4ec3d2311933",
			expectedErr: true,
		},
		{
			name: "unknown currency",
			cps: []mmtypes.Ticker{
				popcat,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(geckoterminal.Name, geckoterminal.DefaultETHProviderConfig)
			require.NoError(t, err)

			h, err := geckoterminal.NewAPIHandler(marketConfig, geckoterminal.DefaultETHAPIConfig)
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
				constants.MOG_USD,
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"id": "8ab62c52-6df2-4613-ad2d-dab08e9e4c8e",
		"type": "simple_token_price",
		"attributes": {
		"token_prices": {
				"0xaaee1a9723aadb7afa2810263653a34ba2c21c7a": "0.000000957896146138212"
			}
		}
	}
}
	`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.MOG_USD: {
						Value: big.NewInt(957896146138),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "malformed response",
			cps:  []mmtypes.Ticker{constants.MOG_USD},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"id": "8ab62c52-6df2-4613-ad2d-dab08e9e4c8e",
		"type": "simple_token_price",
		"attributes": {
		"token_prices": {
				"0xaaee1a9723aadb7afa2810263653a34ba2c21c7a": "0.000000957896146138212",
			}
		}
	}
}
			`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.MOG_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse float",
			cps:  []mmtypes.Ticker{constants.MOG_USD},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"id": "8ab62c52-6df2-4613-ad2d-dab08e9e4c8e",
		"type": "simple_token_price",
		"attributes": {
		"token_prices": {
				"0xaaee1a9723aadb7afa2810263653a34ba2c21c7a": "$0.000000957896146138212"
			}
		}
	}
}
			`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.MOG_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "incorrect attribute",
			cps:  []mmtypes.Ticker{constants.MOG_USD},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"id": "8ab62c52-6df2-4613-ad2d-dab08e9e4c8e",
		"type": "bad_price",
		"attributes": {
		"token_prices": {
				"0xaaee1a9723aadb7afa2810263653a34ba2c21c7a": "0.000000957896146138212"
			}
		}
	}
}
			`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.MOG_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse json",
			cps:  []mmtypes.Ticker{constants.MOG_USD},
			response: testutils.CreateResponseFromJSON(
				`
toms obvious but not minimal language
			`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.MOG_USD: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "multiple currency pairs to parse response for",
			cps: []mmtypes.Ticker{
				constants.MOG_USD,
				constants.PEPE_USD,
			},
			response: testutils.CreateResponseFromJSON(
				`
{
	"data": {
		"id": "8ab62c52-6df2-4613-ad2d-dab08e9e4c8e",
		"type": "simple_token_price",
		"attributes": {
			"token_prices": {
					"0xaaee1a9723aadb7afa2810263653a34ba2c21c7a": "0.000000657896146138212",
					"0x6982508145454Ce325dDbE47a25d4ec3d2311933": "0.000000957896146138212"
			}
		}
	}
}
			`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.MOG_USD: {
						Value: big.NewInt(657896146138),
					},
					constants.PEPE_USD: {
						Value: big.NewInt(957896146138),
					},
				},
				types.UnResolvedPrices{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(geckoterminal.Name, geckoterminal.DefaultETHProviderConfig)
			require.NoError(t, err)

			h, err := geckoterminal.NewAPIHandler(marketConfig, geckoterminal.DefaultETHAPIConfig)
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
