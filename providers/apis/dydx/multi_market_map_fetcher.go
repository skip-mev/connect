package dydx

import (
	"context"
	"fmt"
	"sync"

	"github.com/skip-mev/slinky/oracle/config"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmclient "github.com/skip-mev/slinky/service/clients/marketmap/types"
	"go.uber.org/zap"
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
	cfg config.APIConfig,
	logger *zap.Logger,
) (*MultiMarketMapRestAPIFetcher, error) {
	// make a dydx research api-handler
	researchAPIDataHandler, err := NewResearchAPIHandler(logger, cfg)
	if err != nil {
		return nil, err
	}

	// construct a dydx mainnet api-handler
	if len(cfg.URL) == 0 {
		return nil, fmt.Errorf("no URL provided for dydx mainnet")
	}

	mainnetAPIDataHandler := &APIHandler{
		logger: logger,
		api:    cfg,
	}

	mainnetFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		mainnetAPIDataHandler,
		metrics,
		cfg,
		logger,
	)
	if err != nil {
		return nil, err
	}

	researchFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		researchAPIDataHandler,
		metrics,
		cfg,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return NewDYDXResearchMarketMapFetcher(mainnetFetcher, researchFetcher, logger), nil
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
}

// NewDYDXResearchMarketMapFetcher returns an aggregated market-map among the dydx mainnet and the dydx research json.
func NewDYDXResearchMarketMapFetcher(mainnetFetcher, researchFetcher mmclient.MarketMapFetcher, logger *zap.Logger) *MultiMarketMapRestAPIFetcher {
	return &MultiMarketMapRestAPIFetcher{
		dydxMainnetFetcher:  mainnetFetcher,
		dydxResearchFetcher: researchFetcher,
		logger:              logger.With(zap.String("module", "dydx-research-market-map-fetcher")),
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

	// combine the two market maps
	// if the dydx mainnet market-map response failed, return the dydx mainnet failed response
	if _, ok := dydxMainnetMarketMapResponse.UnResolved[DYDXChain]; ok {
		f.logger.Error("dydx mainnet market-map fetch failed", zap.Any("response", dydxMainnetMarketMapResponse))
		return dydxMainnetMarketMapResponse
	}

	// otherwise, add all markets from dydx research
	dydxMainnetMarketMap := dydxMainnetMarketMapResponse.Resolved[DYDXChain].Value.MarketMap

	if resolved, ok := dydxResearchMarketMapResponse.Resolved[DYDXChain]; ok {
		for ticker, market := range resolved.Value.MarketMap.Markets {
			// if the market is not already in the dydx mainnet market-map, add it
			if _, ok := dydxMainnetMarketMap.Markets[ticker]; !ok {
				f.logger.Debug("adding market from dydx research", zap.String("ticker", ticker))
				dydxMainnetMarketMap.Markets[ticker] = market
			}
		}
	} else {
		return dydxResearchMarketMapResponse
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
