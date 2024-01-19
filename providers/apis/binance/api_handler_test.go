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

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	providerCfg = config.ProviderConfig{
		Name: binance.Name,
		API: config.APIConfig{
			Enabled:    true,
			Timeout:    time.Second,
			Interval:   time.Second,
			MaxQueries: 1,
		},
		CurrencyPairConfig: config.CurrencyPairConfig{
			BaseOffChain: map[string]string{
				"BITCOIN": "BTC",
				"BINANCE": "BNB",
			},
			QuoteOffChain: map[string]string{
				"USDT": "USDT",
			},
		},
	}
)

func TestCreateURL(t *testing.T) {
	testCases := []struct {
		name        string
		cps         []oracletypes.CurrencyPair
		url         string
		expectedErr bool
	}{
		{
			name: "valid single",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "valid multiple",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				oracletypes.NewCurrencyPair("BINANCE", "USDT"),
			},
			url:         "https://api.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22BNBUSDT%22%5D",
			expectedErr: false,
		},
		{
			name: "unknown base currency",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USD"),
			},
			url:         "",
			expectedErr: true,
		},
		{
			name: "unknown quote currency",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			url:         "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := binance.APIHandler{
				ProviderConfig: providerCfg,
				BaseURL:        binance.BaseURL,
			}

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
		cps      []oracletypes.CurrencyPair
		response *http.Response
		expected providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
	}{
		{
			name: "valid single",
			cps:  []oracletypes.CurrencyPair{oracletypes.NewCurrencyPair("BITCOIN", "USDT")},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): {
						Value: big.NewInt(4670703000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
		},
		{
			name: "valid multiple",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
				oracletypes.NewCurrencyPair("BINANCE", "USDT"),
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"46707.03000000"},{"symbol":"BNBUSDT","price":"297.50000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): {
						Value: big.NewInt(4670703000000),
					},
					oracletypes.NewCurrencyPair("BINANCE", "USDT"): {
						Value: big.NewInt(29750000000),
					},
				},
				map[oracletypes.CurrencyPair]error{},
			),
		},
		{
			name: "unknown base",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"MOGUSDT","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "unknown quote",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCMOG","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): fmt.Errorf("no response"),
				},
			),
		},
		{
			name: "unsupported base",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("MOG", "USDT"),
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"MOGUSDT","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("MOG", "USDT"): fmt.Errorf("unknown base currency MOG"),
				},
			),
		},
		{
			name: "unsupported quote",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("USDT", "MOG"),
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"USDTMOG","price":"46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("USDT", "MOG"): fmt.Errorf("unknown quote currency MOG"),
				},
			),
		},
		{
			name: "bad response",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
			},
			response: testutils.CreateResponseFromJSON(
				`shout out my label thats me`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "MOG"): fmt.Errorf("json error"),
				},
			),
		},
		{
			name: "bad price response",
			cps: []oracletypes.CurrencyPair{
				oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			response: testutils.CreateResponseFromJSON(
				`[{"symbol":"BTCUSDT","price":"$46707.03000000"}]`,
			),
			expected: providertypes.NewGetResponse(
				map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{},
				map[oracletypes.CurrencyPair]error{
					oracletypes.NewCurrencyPair("BITCOIN", "USDT"): fmt.Errorf("invalid syntax"),
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := binance.APIHandler{
				ProviderConfig: providerCfg,
			}

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
			h := binance.APIHandler{
				ProviderConfig: providerCfg,
			}
			got, err := h.Decode(tc.response)
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, got, tc.expected)
		})
	}
}
