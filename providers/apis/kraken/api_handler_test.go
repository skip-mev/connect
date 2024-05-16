package kraken_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var (
	btcusd  = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSD",
	}
	btcusdt = types.DefaultProviderTicker{
		OffChainTicker: "BTCUSDT",
	}
	ethusdt = types.DefaultProviderTicker{
		OffChainTicker: "ETHUSDT",
	}
	ethusd  = types.DefaultProviderTicker{
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
			name: "valid single",
			cps: []types.ProviderTicker{
				btcusdt,
			},
			url:         "https://api.kraken.com/0/public/Ticker?pair=XBTUSDT",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []types.ProviderTicker{
				btcusdt,
				ethusdt,
			},
			url:         "https://api.kraken.com/0/public/Ticker?pair=XBTUSDT,ETHUSDT",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := kraken.NewAPIHandler(kraken.DefaultAPIConfig)
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
				`{"error":[],"result":{"XXBTZUSD":{"a":["64587.50000","2","2.000"],"b":["64587.40000","11","11.000"],"c":["64587.40000","0.01026127"],"v":["5866.14264484","6251.33408493"],"p":["64487.45123","64670.54770"],"t":[56819,62596],"l":["62356.50000","62356.50000"],"h":["68075.00000","68075.00000"],"o":"67600.00000"}}}`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(64587.4),
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
				`{"error":[],"result":{"XETHZUSD":{"a":["3338.95000","1","1.000"],"b":["3338.94000","246","246.000"],"c":["3338.08000","0.00702654"],"v":["33234.61736920","35692.20596751"],"p":["3310.12909","3324.16514"],"t":[25646,28278],"l":["3200.17000","3200.17000"],"h":["3547.76000","3547.76000"],"o":"3518.43000"},"XXBTZUSD":{"a":["64547.20000","4","4.000"],"b":["64547.10000","15","15.000"],"c":["64547.20000","0.00013362"],"v":["5869.92462186","6253.84063618"],"p":["64487.50403","64670.01016"],"t":[56856,62595],"l":["62356.50000","62356.50000"],"h":["68075.00000","68075.00000"],"o":"67600.00000"}}}`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(64547.2),
					},
					ethusd: {
						Value: big.NewFloat(3338.08),
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
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(
				`{"error":[],"result":{"XXBTZUSD":{"a":["$64587.50000","2","2.000"],"b":["$64587.40000","11","11.000"],"c":["$64587.40000","0.01026127"],"v":["5866.14264484","6251.33408493"],"p":["64487.45123","64670.54770"],"t":[56819,62596],"l":["62356.50000","62356.50000"],"h":["68075.00000","68075.00000"],"o":"67600.00000"}}}`,
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
			h, err := kraken.NewAPIHandler(kraken.DefaultAPIConfig)
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
		expected  kraken.ResponseBody
		expectErr bool
	}{
		{
			name: "valid single",
			response: testutils.CreateResponseFromJSON(
				`{"error":[],"result":{"XXBTZUSD":{"a":["64587.50000","2","2.000"],"b":["64587.40000","11","11.000"],"c":["64587.40000","0.01026127"],"v":["5866.14264484","6251.33408493"],"p":["64487.45123","64670.54770"],"t":[56819,62596],"l":["62356.50000","62356.50000"],"h":["68075.00000","68075.00000"],"o":"67600.00000"}}}`,
			),
			expected: kraken.ResponseBody{
				Errors: []string{},
				Tickers: map[string]kraken.TickerResult{
					"XXBTZUSD": {
						ClosePriceStats: []string{"64587.40000", "0.01026127"},
					},
				},
			}, expectErr: false,
		},
		{
			name: "valid multi",
			response: testutils.CreateResponseFromJSON(
				`{"error":[],"result":{"XETHZUSD":{"a":["3338.95000","1","1.000"],"b":["3338.94000","246","246.000"],"c":["3338.08000","0.00702654"],"v":["33234.61736920","35692.20596751"],"p":["3310.12909","3324.16514"],"t":[25646,28278],"l":["3200.17000","3200.17000"],"h":["3547.76000","3547.76000"],"o":"3518.43000"},"XXBTZUSD":{"a":["64547.20000","4","4.000"],"b":["64547.10000","15","15.000"],"c":["64547.20000","0.00013362"],"v":["5869.92462186","6253.84063618"],"p":["64487.50403","64670.01016"],"t":[56856,62595],"l":["62356.50000","62356.50000"],"h":["68075.00000","68075.00000"],"o":"67600.00000"}}}`,
			),
			expected: kraken.ResponseBody{
				Errors: []string{},
				Tickers: map[string]kraken.TickerResult{
					"XETHZUSD": {
						ClosePriceStats: []string{"3338.08000", "0.00702654"},
					},
					"XXBTZUSD": {
						ClosePriceStats: []string{"64547.20000", "0.00013362"},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := kraken.Decode(tc.response)
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, got, tc.expected)
		})
	}
}
