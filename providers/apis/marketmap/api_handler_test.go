package marketmap_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/base/testutils"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

var (
	chains = []types.Chain{
		{
			ChainID: "dYdX",
		},
		{
			ChainID: "osmosis",
		},
	}

	goodMarketMap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
		},
		Providers: map[string]mmtypes.Providers{
			constants.BITCOIN_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           "coinbase",
						OffChainTicker: "BTC/USD",
					},
				},
			},
		},
	}

	badMarketMap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
		},
	}
)

func TestCreateURL(t *testing.T) {
	apiHandler, err := marketmap.NewAPIHandler(marketmap.DefaultAPIConfig)
	require.NoError(t, err)

	t.Run("errors when there are multiple chains inputted", func(t *testing.T) {
		_, err := apiHandler.CreateURL(chains)
		require.Error(t, err)
	})

	t.Run("returns the URL when there is only one chain inputted", func(t *testing.T) {
		url, err := apiHandler.CreateURL(chains[:1])
		require.NoError(t, err)
		require.Equal(t, marketmap.DefaultAPIConfig.URL, url)
	})
}

func TestParseResponse(t *testing.T) {
	cases := []struct {
		name     string
		chains   []types.Chain
		resp     func() *http.Response
		expected types.MarketMapResponse
	}{
		{
			name:   "errors when too many chains are inputted",
			chains: chains,
			resp: func() *http.Response {
				return &http.Response{}
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: fmt.Errorf("expected one chain, got 2"),
					chains[1]: fmt.Errorf("expected one chain, got 2"),
				},
			},
		},
		{
			name:   "errors when the response is nil",
			chains: chains[:1],
			resp: func() *http.Response {
				return nil
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: fmt.Errorf("nil response"),
				},
			},
		},
		{
			name:   "errors when the response body cannot be parsed",
			chains: chains[:1],
			resp: func() *http.Response {
				return testutils.CreateResponseFromJSON("invalid json")
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: fmt.Errorf("failed to parse market map response"),
				},
			},
		},
		{
			name:   "errors when the market map response is invalid",
			chains: chains[:1],
			resp: func() *http.Response {
				resp := mmtypes.GetMarketMapResponse{
					MarketMap: badMarketMap,
				}

				json, err := json.Marshal(resp)
				require.NoError(t, err)

				return testutils.CreateResponseFromJSON(string(json))
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: fmt.Errorf("invalid market map response"),
				},
			},
		},
		{
			name:   "returns a market map that does not match the chain id",
			chains: chains[:1],
			resp: func() *http.Response {
				resp := mmtypes.GetMarketMapResponse{
					MarketMap: goodMarketMap,
					ChainId:   "invalid",
				}

				json, err := json.Marshal(resp)
				require.NoError(t, err)

				return testutils.CreateResponseFromJSON(string(json))
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: fmt.Errorf("expected chain id dYdX, got invalid"),
				},
			},
		},
		{
			name:   "returns a resolved market map",
			chains: chains[:1],
			resp: func() *http.Response {
				resp := mmtypes.GetMarketMapResponse{
					MarketMap: goodMarketMap,
					ChainId:   chains[0].ChainID,
				}

				json, err := json.Marshal(resp)
				require.NoError(t, err)

				return testutils.CreateResponseFromJSON(string(json))
			},
			expected: types.MarketMapResponse{
				Resolved: types.ResolvedMarketMap{
					chains[0]: types.MarketMapResult{
						Value: &mmtypes.GetMarketMapResponse{
							MarketMap: goodMarketMap,
							ChainId:   chains[0].ChainID,
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			apiHandler, err := marketmap.NewAPIHandler(marketmap.DefaultAPIConfig)
			require.NoError(t, err)

			resp := apiHandler.ParseResponse(tc.chains, tc.resp())
			fmt.Println(resp.String())
			require.Equal(t, len(resp.Resolved), len(tc.expected.Resolved))
			require.Equal(t, len(resp.UnResolved), len(tc.expected.UnResolved))

			for cp, result := range tc.expected.Resolved {
				require.Contains(t, resp.Resolved, cp)
				r := resp.Resolved[cp]
				require.Equal(t, result.Value, r.Value)
			}

			for cp := range tc.expected.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}
