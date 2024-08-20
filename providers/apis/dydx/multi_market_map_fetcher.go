package dydx

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/cmd/constants/marketmaps"
	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmclient "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	_         mmclient.MarketMapFetcher = &MultiMarketMapRestAPIFetcher{}
	DYDXChain                           = mmclient.Chain{
		ChainID: ChainID,
	}
)

// NewDYDXResearchMarketMapFetcher returns a MultiMarketMapFetcher composed of dydx mainnet + research
// apiDataHandlers.
func DefaultDYDXResearchMarketMapFetcher(
	rh apihandlers.RequestHandler,
	metrics metrics.APIMetrics,
	api config.APIConfig,
	logger *zap.Logger,
) (*MultiMarketMapRestAPIFetcher, error) {
	if rh == nil {
		return nil, fmt.Errorf("request handler is nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, err
	}

	if len(api.Endpoints) != 2 {
		return nil, fmt.Errorf("expected two endpoint, got %d", len(api.Endpoints))
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	// make a dydx research api-handler
	researchAPIDataHandler, err := NewResearchAPIHandler(logger, api)
	if err != nil {
		return nil, err
	}

	mainnetAPIDataHandler := &APIHandler{
		logger: logger,
		api:    api,
	}

	mainnetFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		mainnetAPIDataHandler,
		metrics,
		api,
		logger,
	)
	if err != nil {
		return nil, err
	}

	researchFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		researchAPIDataHandler,
		metrics,
		api,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return NewDYDXResearchMarketMapFetcher(
		mainnetFetcher,
		researchFetcher,
		logger,
		api.Name == ResearchCMCAPIHandlerName,
	), nil
}

// MultiMarketMapRestAPIFetcher is an implementation of a RestAPIFetcher that wraps
// two underlying Fetchers for fetching the market-map according to dydx mainnet and
// the additional markets that can be added according to the dydx research json.
type MultiMarketMapRestAPIFetcher struct {
	// dydx mainnet fetcher is the api-fetcher for the dydx mainnet market-map
	dydxMainnetFetcher mmclient.MarketMapFetcher

	// dydx research fetcher is the api-fetcher for the dydx research market-map
	dydxResearchFetcher mmclient.MarketMapFetcher

	// logger is the logger for the fetcher
	logger *zap.Logger

	// isCMCOnly is a flag that indicates whether the fetcher should only return CoinMarketCap markets.
	isCMCOnly bool
}

// NewDYDXResearchMarketMapFetcher returns an aggregated market-map among the dydx mainnet and the dydx research json.
func NewDYDXResearchMarketMapFetcher(
	mainnetFetcher, researchFetcher mmclient.MarketMapFetcher,
	logger *zap.Logger,
	isCMCOnly bool,
) *MultiMarketMapRestAPIFetcher {
	return &MultiMarketMapRestAPIFetcher{
		dydxMainnetFetcher:  mainnetFetcher,
		dydxResearchFetcher: researchFetcher,
		logger:              logger.With(zap.String("module", "dydx-research-market-map-fetcher")),
		isCMCOnly:           isCMCOnly,
	}
}

// Fetch fetches the market map from the underlying fetchers and combines the results. If any of the underlying
// fetchers fetch for a chain that is different from the chain that the fetcher is initialized with, those responses
// will be ignored.
func (f *MultiMarketMapRestAPIFetcher) Fetch(ctx context.Context, chains []mmclient.Chain) mmclient.MarketMapResponse {
	// call the underlying fetchers + await their responses
	// channel to aggregate responses
	dydxMainnetResponseChan := make(chan mmclient.MarketMapResponse, 1) // buffer so that sends / receives are non-blocking
	dydxResearchResponseChan := make(chan mmclient.MarketMapResponse, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	// fetch dydx mainnet
	go func() {
		defer wg.Done()
		dydxMainnetResponseChan <- f.dydxMainnetFetcher.Fetch(ctx, chains)
		f.logger.Debug("fetched valid market-map from dydx mainnet")
	}()

	// fetch dydx research
	go func() {
		defer wg.Done()
		dydxResearchResponseChan <- f.dydxResearchFetcher.Fetch(ctx, chains)
		f.logger.Debug("fetched valid market-map from dydx research")
	}()

	// wait for both fetchers to finish
	wg.Wait()

	dydxMainnetMarketMapResponse := <-dydxMainnetResponseChan
	dydxResearchMarketMapResponse := <-dydxResearchResponseChan

	// if the dydx mainnet market-map response failed, return the dydx mainnet failed response
	if _, ok := dydxMainnetMarketMapResponse.UnResolved[DYDXChain]; ok {
		f.logger.Error("dydx mainnet market-map fetch failed", zap.Any("response", dydxMainnetMarketMapResponse))
		return dydxMainnetMarketMapResponse
	}

	// if the dydx research market-map response failed, return the dydx research failed response
	if _, ok := dydxResearchMarketMapResponse.UnResolved[DYDXChain]; ok {
		f.logger.Error("dydx research market-map fetch failed", zap.Any("response", dydxResearchMarketMapResponse))
		return dydxResearchMarketMapResponse
	}

	// otherwise, add all markets from dydx research
	dydxMainnetMarketMap := dydxMainnetMarketMapResponse.Resolved[DYDXChain].Value.MarketMap

	resolved, ok := dydxResearchMarketMapResponse.Resolved[DYDXChain]
	if ok {
		for ticker, market := range resolved.Value.MarketMap.Markets {
			// if the market is not already in the dydx mainnet market-map, add it
			if _, ok := dydxMainnetMarketMap.Markets[ticker]; !ok {
				f.logger.Debug("adding market from dydx research", zap.String("ticker", ticker))
				dydxMainnetMarketMap.Markets[ticker] = market
			}
		}
	}

	// if the fetcher is only for CoinMarketCap markets, filter out all non-CMC markets
	if f.isCMCOnly {
		for ticker, market := range dydxMainnetMarketMap.Markets {
			market.Ticker.MinProviderCount = 1
			dydxMainnetMarketMap.Markets[ticker] = market

			var (
				seenCMC     = false
				cmcProvider mmtypes.ProviderConfig
			)

			for _, provider := range market.ProviderConfigs {
				if provider.Name == coinmarketcap.Name {
					seenCMC = true
					cmcProvider = provider
				}
			}

			// if we saw a CMC provider, add it to the market
			if seenCMC {
				market.ProviderConfigs = []mmtypes.ProviderConfig{cmcProvider}
				dydxMainnetMarketMap.Markets[ticker] = market
				continue
			}

			// If we did not see a CMC provider, we can attempt to add it using the CMC marketmap
			cmcMarket, ok := marketmaps.CoinMarketCapMarketMap.Markets[ticker]
			if !ok {
				f.logger.Info("did not find CMC market for ticker", zap.String("ticker", ticker))
				delete(dydxMainnetMarketMap.Markets, ticker)
				continue
			}

			// add the CMC provider to the market
			market.ProviderConfigs = cmcMarket.ProviderConfigs
			dydxMainnetMarketMap.Markets[ticker] = market
		}
	}

	// validate the combined market-map
	if err := dydxMainnetMarketMap.ValidateBasic(); err != nil {
		f.logger.Error("combined market-map failed validation", zap.Error(err))

		return mmclient.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("combined market-map failed validation: %w", err),
				providertypes.ErrorUnknown,
			),
		)
	}

	dydxMainnetMarketMapResponse.Resolved[DYDXChain].Value.MarketMap = dydxMainnetMarketMap

	return dydxMainnetMarketMapResponse
}
