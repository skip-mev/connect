package marketmap_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/marketmap/types/mocks"
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

	btcusd = connecttypes.NewCurrencyPair("BTC", "USD")

	goodMarketMap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			btcusd.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     btcusd,
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
			btcusd.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     btcusd,
					Decimals:         8,
					MinProviderCount: 3,
				},
			},
		},
	}

	logger = zap.NewExample()
)

func TestFetch(t *testing.T) {
	cases := []struct {
		name     string
		chains   []types.Chain
		client   func() mmtypes.QueryClient
		expected types.MarketMapResponse
	}{
		{
			name:   "errors when too many chains are inputted",
			chains: chains,
			client: func() mmtypes.QueryClient {
				return mocks.NewQueryClient(t)
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("expected one chain, got 2"), providertypes.ErrorInvalidAPIChains),
					},
					chains[1]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("expected one chain, got 2"), providertypes.ErrorInvalidAPIChains),
					},
				},
			},
		},
		{
			name:   "errors when the response is nil",
			chains: chains[:1],
			client: func() mmtypes.QueryClient {
				c := mocks.NewQueryClient(t)
				c.On("MarketMap", mock.Anything, mock.Anything).Return(nil, nil)
				return c
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("nil response"), providertypes.ErrorGRPCGeneral),
					},
				},
			},
		},
		{
			name:   "errors when the request cannot be made",
			chains: chains[:1],
			client: func() mmtypes.QueryClient {
				c := mocks.NewQueryClient(t)
				c.On("MarketMap", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("could not make request"))
				return c
			},
			expected: types.MarketMapResponse{
				UnResolved: types.UnResolvedMarketMap{
					chains[0]: providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("could not make request"), providertypes.ErrorGRPCGeneral),
					},
				},
			},
		},
		{
			name:   "does not error when the market map response is invalid",
			chains: chains[:1],
			client: func() mmtypes.QueryClient {
				c := mocks.NewQueryClient(t)
				c.On("MarketMap", mock.Anything, mock.Anything).Return(
					&mmtypes.MarketMapResponse{
						MarketMap:   badMarketMap,
						ChainId:     chains[0].ChainID,
						LastUpdated: 11,
					},
					nil,
				)
				return c
			},
			expected: types.MarketMapResponse{
				Resolved: types.ResolvedMarketMap{
					chains[0]: types.MarketMapResult{
						Value: &mmtypes.MarketMapResponse{
							MarketMap:   badMarketMap,
							ChainId:     chains[0].ChainID,
							LastUpdated: 11,
						},
					},
				},
			},
		},
		{
			name:   "returns a resolved market map",
			chains: chains[:1],
			client: func() mmtypes.QueryClient {
				c := mocks.NewQueryClient(t)
				c.On("MarketMap", mock.Anything, mock.Anything).Return(
					&mmtypes.MarketMapResponse{
						MarketMap:   goodMarketMap,
						ChainId:     chains[0].ChainID,
						LastUpdated: 10,
					},
					nil,
				)
				return c
			},
			expected: types.MarketMapResponse{
				Resolved: types.ResolvedMarketMap{
					chains[0]: types.MarketMapResult{
						Value: &mmtypes.MarketMapResponse{
							MarketMap:   goodMarketMap,
							ChainId:     chains[0].ChainID,
							LastUpdated: 10,
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fetcher, err := marketmap.NewMarketMapFetcherWithClient(logger, tc.client())
			require.NoError(t, err)

			resp := fetcher.Fetch(context.TODO(), tc.chains)
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
