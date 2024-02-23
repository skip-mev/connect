package oracle

import (
	"net/http"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/base"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

// DefaultMarketMapProvider returns a sample implementation of the market map provider. This provider
// is responsible for fetching updates to the cannonical market map from the x/marketmap module.
func DefaultMarketMapProvider(
	logger *zap.Logger,
	providerMetrics providermetrics.ProviderMetrics,
	apiMetrics apimetrics.APIMetrics,
	cfg config.ProviderConfig,
) (types.MarketMapProvider, error) {
	apiDataHandler, err := marketmap.NewAPIHandler(cfg.API)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: cfg.API.MaxQueries},
		Timeout:   cfg.API.Timeout,
	}
	requestHandler, err := apihandlers.NewRequestHandlerImpl(client)
	if err != nil {
		return nil, err
	}

	queryHandler, err := types.NewMarketMapAPIQueryHandler(
		logger,
		cfg.API,
		requestHandler,
		apiDataHandler,
		apiMetrics,
	)
	if err != nil {
		return nil, err
	}

	return types.NewMarketMapProvider(
		base.WithName[types.Chain, *mmtypes.GetMarketMapResponse](cfg.Name),
		base.WithLogger[types.Chain, *mmtypes.GetMarketMapResponse](logger),
		base.WithAPIQueryHandler[types.Chain, *mmtypes.GetMarketMapResponse](queryHandler),
		base.WithAPIConfig[types.Chain, *mmtypes.GetMarketMapResponse](cfg.API),
		base.WithMetrics[types.Chain, *mmtypes.GetMarketMapResponse](providerMetrics),
	)
}
