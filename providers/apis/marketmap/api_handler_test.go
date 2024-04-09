package marketmap_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
		Markets: map[string]mmtypes.Market{
			constants.BITCOIN_USD.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     constants.BITCOIN_USD,
					Decimals:         8,
					MinProviderCount: 1,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "BTC-USD",
					},
				},
			},
		},
	}

	badMarketMap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			constants.BITCOIN_USD.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     constants.BITCOIN_USD,
					Decimals:         8,
					MinProviderCount: 3,
				},
			},
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
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("expected one chain, got 2"), providertypes.ErrorAPIGeneral),
					},
					chains[1]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("expected one chain, got 2"), providertypes.ErrorAPIGeneral),
					},
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
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("nil response"), providertypes.ErrorAPIGeneral),
					},
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
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("failed to parse market map response"), providertypes.ErrorAPIGeneral),
					},
				},
			},
		},
		{
			name:   "errors when the market map response is invalid",
			chains: chains[:1],
			resp: func() *http.Response {
				resp := mmtypes.MarketMapResponse{
					MarketMap: badMarketMap,
				}

				json, err := json.Marshal(resp)
				require.NoError(t, err)

				return testutils.CreateResponseFromJSON(string(json))
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("invalid market map response"), providertypes.ErrorAPIGeneral),
					},
				},
			},
		},
		{
			name:   "returns a market map that does not match the chain id",
			chains: chains[:1],
			resp: func() *http.Response {
				resp := mmtypes.MarketMapResponse{
					MarketMap: goodMarketMap,
					ChainId:   "invalid",
				}

				json, err := json.Marshal(resp)
				require.NoError(t, err)

				return testutils.CreateResponseFromJSON(string(json))
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("expected chain id dYdX, got invalid"), providertypes.ErrorAPIGeneral),
					},
				},
			},
		},
		{
			name:   "returns a resolved market map",
			chains: chains[:1],
			resp: func() *http.Response {
				resp := mmtypes.MarketMapResponse{
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
						Value: &mmtypes.MarketMapResponse{
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
