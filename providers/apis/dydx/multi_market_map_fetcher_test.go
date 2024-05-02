package dydx_test

import (
	"context"
	"testing"
	"time"

	"github.com/skip-mev/slinky/providers/apis/dydx"
	apihandlermocks "github.com/skip-mev/slinky/providers/base/api/handlers/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmclient "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestDYDXMultiMarketMapFetcher(t *testing.T) {
	dydxMainnetMMFetcher := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)
	dydxResearchMMFetcher := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)

	fetcher := dydx.NewDYDXResearchMarketMapFetcher(dydxMainnetMMFetcher, dydxResearchMMFetcher)

	t.Run("test that if the mainnet api-price fetcher response is unresolved, we return it", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			UnResolved: map[mmclient.Chain]providertypes.UnresolvedResult{
				dydx.DYDXChain: providertypes.UnresolvedResult{},
			},
		}, nil)
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{}, nil)

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test that if the dydx-research response is unresolved, we return that", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{}, time.Now()),
			},
		}, nil)
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			UnResolved: map[mmclient.Chain]providertypes.UnresolvedResult{
				dydx.DYDXChain: providertypes.UnresolvedResult{},
			},
		}, nil)

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test if both responses are resolved, the tickers are appended to each other", func(t *testing.T) {
		ctx := context.Background()
		dydxMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"BTC/USD": mmtypes.Market{},
						},
					},
				}, time.Now()),
			},
		}, nil)
		dydxResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{dydx.DYDXChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				dydx.DYDXChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"ETH/USD": mmtypes.Market{},
						},
					},
				}, time.Now()),
			},
		}, nil)

		response := fetcher.Fetch(ctx, []mmclient.Chain{dydx.DYDXChain})
		require.Len(t, response.Resolved, 1)

		marketMap := response.Resolved[dydx.DYDXChain].Value.MarketMap

		require.Len(t, marketMap.Markets, 2)
		require.Contains(t, marketMap.Markets, "BTC/USD")
		require.Contains(t, marketMap.Markets, "ETH/USD")
	})
}