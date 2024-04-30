package oracle

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/base"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				// Force IPv4 by specifying the network type as "tcp4"
				Resolver: &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						return net.Dial("tcp4", address)
					},
				},
			}).DialContext,
		},
		Timeout: cfg.API.Timeout,
	}
	var apiDataHandler types.MarketMapAPIDataHandler
	var requestHandler apihandlers.RequestHandler
	var ids []types.Chain

	switch cfg.Name {
	case dydx.Name:
		apiDataHandler, err = dydx.NewAPIHandler(logger, cfg.API)
		ids = []types.Chain{{ChainID: dydx.ChainID}}
	default:
		apiDataHandler, err = marketmap.NewAPIHandler(cfg.API)
	}
	if err != nil {
		return nil, err
	}
	requestHandler, err = apihandlers.NewRequestHandlerImpl(client)
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
		base.WithName[types.Chain, *mmtypes.MarketMapResponse](cfg.Name),
		base.WithLogger[types.Chain, *mmtypes.MarketMapResponse](logger),
		base.WithAPIQueryHandler(queryHandler),
		base.WithAPIConfig[types.Chain, *mmtypes.MarketMapResponse](cfg.API),
		base.WithMetrics[types.Chain, *mmtypes.MarketMapResponse](providerMetrics),
		base.WithIDs[types.Chain, *mmtypes.MarketMapResponse](ids),
	)
}
