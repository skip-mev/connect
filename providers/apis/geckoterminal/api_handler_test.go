package geckoterminal_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/geckoterminal"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var (
	mogusd  = types.DefaultProviderTicker{
		OffChainTicker: "MOGUSD",
	}
	pepeusd = types.DefaultProviderTicker{
		OffChainTicker: "PEPEUSD",
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
			name: "valid",
			cps: []types.ProviderTicker{
				mogusd,
			},
			url:         "https://api.geckoterminal.com/api/v2/simple/networks/eth/token_price/0xaaee1a9723aadb7afa2810263653a34ba2c21c7a",
			expectedErr: false,
		},
		{
			name: "multiple currency pairs",
			cps: []types.ProviderTicker{
				mogusd,
				pepeusd,
			},
			url:         "https://api.geckoterminal.com/api/v2/simple/networks/eth/token_price/0xaaee1a9723aadb7afa2810263653a34ba2c21c7a,0x6982508145454Ce325dDbE47a25d4ec3d2311933",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := geckoterminal.NewAPIHandler(geckoterminal.DefaultETHAPIConfig)
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
					mogusd: {
						Value: big.NewFloat(0.000000957896146138212),
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
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
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
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "incorrect attribute",
			cps: []types.ProviderTicker{
				mogusd,
			},
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
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "unable to parse json",
			cps: []types.ProviderTicker{
				mogusd,
			},
			response: testutils.CreateResponseFromJSON(
				`
toms obvious but not minimal language
			`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							fmt.Errorf("bad format"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "multiple currency pairs to parse response for",
			cps: []types.ProviderTicker{
				mogusd,
				pepeusd,
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
					mogusd: {
						Value: big.NewFloat(0.000000657896146138212),
					},
					pepeusd: {
						Value: big.NewFloat(0.000000957896146138212),
					},
				},
				types.UnResolvedPrices{},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := geckoterminal.NewAPIHandler(geckoterminal.DefaultETHAPIConfig)
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
