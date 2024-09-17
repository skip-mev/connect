package dydx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	apihandlermocks "github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmclient "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestDYDXMultiMarketMapFetcher(t *testing.T) {
	dydxMainnetMMFetcher := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)
	dydxResearchMMFetcher := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)

	fetcher := dydx.NewDYDXResearchMarketMapFetcher(dydxMainnetMMFetcher, dydxResearchMMFetcher, zap.NewExample(), false)

	t.Run("test that if the mainnet api-price fetcher response is unresolved, we return it", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			UnResolved: map[mmclient.Chain]providertypes.UnresolvedResult{
				dydx.DYDXChain: {
					ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("error"), providertypes.ErrorAPIGeneral),
				},
			},
		}, nil).Once()
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test that if the dydx-research response is unresolved, we return that", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{}, time.Now()),
			},
		}, nil).Once()
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			UnResolved: map[mmclient.Chain]providertypes.UnresolvedResult{
				dydx.DYDXChain: {},
			},
		}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test if both responses are resolved, the tickers are appended to each other + validation fails", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"BTC/USD": {},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"ETH/USD": {},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test that if both responses are resolved, the responses are aggregated + validation passes", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"BTC/USD": {
								Ticker: mmtypes.Ticker{
									CurrencyPair:     connecttypes.NewCurrencyPair("BTC", "USD"),
									Decimals:         8,
									MinProviderCount: 1,
									Enabled:          true,
								},
								ProviderConfigs: []mmtypes.ProviderConfig{
									{
										Name:           "dydx",
										OffChainTicker: "BTC/USD",
									},
								},
							},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"ETH/USD": {
								Ticker: mmtypes.Ticker{
									CurrencyPair:     connecttypes.NewCurrencyPair("ETH", "USD"),
									Decimals:         8,
									MinProviderCount: 1,
								},
								ProviderConfigs: []mmtypes.ProviderConfig{
									{
										Name:           "dydx",
										OffChainTicker: "BTC/USD",
									},
								},
							},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.Resolved, 1)

		marketMap := response.Resolved[dydx.DYDXChain].Value.MarketMap

		require.Len(t, marketMap.Markets, 2)
		require.Contains(t, marketMap.Markets, "BTC/USD")
		require.Contains(t, marketMap.Markets, "ETH/USD")
	})
}
