package binance_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
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
			name: "valid single",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22ETHUSDT%22%5D",
			expectedErr: false,
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
			marketConfig, err := types.NewProviderMarketMap(binance.Name, binance.DefaultNonUSProviderConfig)
			require.NoError(t, err)

			h, err := binance.NewAPIHandler(marketConfig, binance.DefaultNonUSAPIConfig)
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

func TestCreateURL_US(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []mmtypes.Ticker
		url         string
		expectedErr bool
	}{
		{
			name: "valid single",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
			},
			url:         "https://api.binance.us/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
			},
			url:         "https://api.binance.us/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22ETHUSDT%22%5D",
			expectedErr: false,
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
			marketConfig, err := types.NewProviderMarketMap(binance.Name, binance.DefaultUSProviderConfig)
			require.NoError(t, err)

			h, err := binance.NewAPIHandler(marketConfig, binance.DefaultUSAPIConfig)
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
			name: "valid single",
			cps:  []mmtypes.Ticker{constants.BITCOIN_USDT},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USDT: {
						Value: big.NewInt(4670703000000),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "valid multiple",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"},{"symbol":"ETHUSDT","price":"297.50000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					constants.BITCOIN_USDT: {
						Value: big.NewInt(4670703000000),
					},
					constants.ETHEREUM_USDT: {
						Value: big.NewInt(29750000000),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "unsupported currency",
			cps: []mmtypes.Ticker{
				mogusd,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"MOGUSDT","price":"46707.03000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					mogusd: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "bad response",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
			},
			response: testutils.CreateResponseFromJSON(
				`shout out my label thats me`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USDT: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "bad price response",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"$46707.03000000"}]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USDT: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid syntax"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
		{
			name: "no response",
			cps: []mmtypes.Ticker{
				constants.BITCOIN_USDT,
				constants.ETHEREUM_USDT,
			},
			response: testutils.CreateResponseFromJSON(
				`[]`,
			),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					constants.BITCOIN_USDT: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
					constants.ETHEREUM_USDT: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorAPIGeneral),
					},
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marketConfig, err := types.NewProviderMarketMap(binance.Name, binance.DefaultUSProviderConfig)
			require.NoError(t, err)

			h, err := binance.NewAPIHandler(marketConfig, binance.DefaultNonUSAPIConfig)
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
