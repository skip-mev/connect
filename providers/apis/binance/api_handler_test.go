package binance_test

// NOTE: some binance tests currently use the binance.us endpoints because the standard binance endpoints
// are georestricted

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	btcusdt = mmtypes.NewTicker("BITCOIN", "USDT", 8, 1)
	bnbusdt = mmtypes.NewTicker("BINANCE", "USDT", 8, 1)
	mogusdt = mmtypes.NewTicker("MOG", "USDT", 8, 1)

	marketConfig = mmtypes.MarketConfig{
		Name: binance.Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"BITCOIN/USDT": {
				Ticker:         btcusdt,
				OffChainTicker: "BTCUSDT",
			},
			"BINANCE/USDT": {
				Ticker:         bnbusdt,
				OffChainTicker: "BNBUSDT",
			},
		},
	}
)

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
				btcusdt,
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []mmtypes.Ticker{
				btcusdt,
				bnbusdt,
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22BNBUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "unknown currency",
			cps: []mmtypes.Ticker{
				mogusdt,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
				btcusdt,
			},
			url:         "https://api.binance.us/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []mmtypes.Ticker{
				btcusdt,
				bnbusdt,
			},
			url:         "https://api.binance.us/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22BNBUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "unknown currency",
			cps: []mmtypes.Ticker{
				mogusdt,
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
		expected providertypes.GetResponse[mmtypes.Ticker, *big.Int]
	}{
		{
			name: "valid single",
			cps:  []mmtypes.Ticker{btcusdt},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{
					btcusdt: {
						Value: big.NewInt(4670703000000),
					},
				},
				map[mmtypes.Ticker]error{},
			),
		},
		{
			name: "valid multiple",
			cps: []mmtypes.Ticker{
				btcusdt,
				bnbusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"},{"symbol":"BNBUSDT","price":"297.50000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{
					btcusdt: {
						Value: big.NewInt(4670703000000),
					},
					bnbusdt: {
						Value: big.NewInt(29750000000),
					},
				},
				map[mmtypes.Ticker]error{},
			),
		},
		{
			name: "unsupported currency",
			cps: []mmtypes.Ticker{
				mogusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"MOGUSDT","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					mogusdt: fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "bad response",
			cps: []mmtypes.Ticker{
				btcusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`shout out my label thats me`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					btcusdt: fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "bad price response",
			cps: []mmtypes.Ticker{
				btcusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"$46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					btcusdt: fmt.Errorf("invalid syntax"),
				},
			),
		},
		{
			name: "no response",
			cps: []mmtypes.Ticker{
				btcusdt,
				bnbusdt,
			},
			response: testutils.CreateResponseFromJSON(
				`[]`,
			),
			expected: providertypes.NewGetResponse(
				map[mmtypes.Ticker]providertypes.Result[*big.Int]{},
				map[mmtypes.Ticker]error{
					btcusdt: fmt.Errorf("no response"),
					bnbusdt: fmt.Errorf("no response"),
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
				`[{"symbol":"BTCUSDT","price":"46707.03000000"},{"symbol":"BNBUSDT","price":"707.03000000"}]`,
			),
			expected: binance.Response{
				binance.Data{
					Symbol: "BTCUSDT",
					Price:  "46707.03000000",
				},
				binance.Data{
					Symbol: "BNBUSDT",
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
