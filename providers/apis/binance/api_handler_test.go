package binance_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var (
	btcusdt = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSDT",
	}
	ethusdt = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSDT",
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
				btcusdt,
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22ETHUSDT%22%5D",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := binance.NewAPIHandler(binance.DefaultNonUSAPIConfig)
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
				btcusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusdt: {
						Value: big.NewFloat(46707.03),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "valid multiple",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"},{"symbol":"ETHUSDT","price":"297.50000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusdt: {
						Value: big.NewFloat(46707.03),
					},
					ethusdt: {
						Value: big.NewFloat(297.5),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "bad response",
			cps: []types.ProviderTicker{
				btcusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`shout out my label that's me`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "bad price response",
			cps: []types.ProviderTicker{
				btcusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"$46707.03000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid syntax"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "no response",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
					ethusdt: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := binance.NewAPIHandler(binance.DefaultNonUSAPIConfig)
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

func TestDecode(t *testing.T) {
	testCases := []struct {
		name      string
		response  *http.Response
		expected  binance.Response
		expectErr bool
	}{
		{
			name: "valid single",
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"}]`,
			),
			expected: binance.Response{
				binance.Data{
					Symbol: "BTCUSDT",
					Price:  "46707.03000000",
				},
			},
			expectErr: false,
		},
		{
			name: "valid multi",
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"},{"symbol":"ETHUSDT","price":"707.03000000"}]`,
			),
			expected: binance.Response{
				binance.Data{
					Symbol: "BTCUSDT",
					Price:  "46707.03000000",
				},
				binance.Data{
					Symbol: "ETHUSDT",
					Price:  "707.03000000",
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := binance.Decode(tc.response)
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, got, tc.expected)
		})
	}
}
