package oracle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	"github.com/skip-mev/connect/v2/providers/base"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	"github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// MarketMapProviderFactory returns a sample implementation of the market map provider. This provider
// is responsible for fetching updates to the canonical market map on the given chain.
func MarketMapProviderFactory(
	logger *zap.Logger,
	providerMetrics providermetrics.ProviderMetrics,
	apiMetrics apimetrics.APIMetrics,
	cfg config.ProviderConfig,
) (*types.MarketMapProvider, error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: cfg.API.MaxQueries,
			Proxy:           http.ProxyFromEnvironment,
		},
		Timeout: cfg.API.Timeout,
	}

	var (
		apiDataHandler   types.MarketMapAPIDataHandler
		ids              []types.Chain
		marketMapFetcher types.MarketMapFetcher
	)

	requestHandler, err := apihandlers.NewRequestHandlerImpl(client)
	if err != nil {
		return nil, err
	}

	switch cfg.Name {
	case dydx.Name:
		apiDataHandler, err = dydx.NewAPIHandler(logger, cfg.API)
		ids = []types.Chain{{ChainID: dydx.ChainID}}
	case dydx.SwitchOverAPIHandlerName:
		marketMapFetcher, err = dydx.NewDefaultSwitchOverMarketMapFetcher(
			logger,
			cfg.API,
			requestHandler,
			apiMetrics,
		)
		ids = []types.Chain{{ChainID: dydx.ChainID}}
	case dydx.ResearchAPIHandlerName, dydx.ResearchCMCAPIHandlerName:
		marketMapFetcher, err = dydx.DefaultDYDXResearchMarketMapFetcher(
			requestHandler,
			apiMetrics,
			cfg.API,
			logger,
		)
		ids = []types.Chain{{ChainID: dydx.ChainID}}
	default:
		marketMapFetcher, err = marketmap.NewMarketMapFetcher(
			logger,
			cfg.API,
			apiMetrics,
		)
		ids = []types.Chain{{ChainID: "local-node"}}
	}
	if err != nil {
		return nil, err
	}

	if marketMapFetcher == nil {
		marketMapFetcher, err = apihandlers.NewRestAPIFetcher(
			requestHandler,
			apiDataHandler,
			apiMetrics,
			cfg.API,
			logger,
		)
		if err != nil {
			return nil, err
		}
	}

	queryHandler, err := types.NewMarketMapAPIQueryHandlerWithMarketMapFetcher(
		logger,
		cfg.API,
		marketMapFetcher,
		apiMetrics,
	)
	if err != nil {
		return nil, err
	}

	return types.NewMarketMapProvider(
		base.WithName[types.Chain, *mmtypes.MarketMapResponse](cfg.Name),
		base.WithLogger[types.Chain, *mmtypes.MarketMapResponse](logger),
		base.WithAPIQueryHandler(queryHandler),
		base.WithAPIConfig[types.Chain, *mmtypes.MarketMapResponse](cfg.API),
		base.WithMetrics[types.Chain, *mmtypes.MarketMapResponse](providerMetrics),
		base.WithIDs[types.Chain, *mmtypes.MarketMapResponse](ids),
	)
}
