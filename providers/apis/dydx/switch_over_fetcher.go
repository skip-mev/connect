package dydx

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/connect/v2/oracle/config"
	connectgrpc "github.com/skip-mev/connect/v2/pkg/grpc"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	mmclient "github.com/skip-mev/connect/v2/service/clients/marketmap/types"

	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var _ mmclient.MarketMapFetcher = &SwitchOverFetcher{}

// SwitchOverFetcher is an implementation of a RestAPIFetcher that wraps a
// dydx x/prices market map fetcher and a x/marketmap fetcher. The fetcher
// operates by first fetching the market map from the x/prices API and then
// fetching the market map from the x/marketmap API. The fetcher will switch
// over to the x/marketmap API the first time it receives a non-nil market map
// from the x/marketmap API.
type SwitchOverFetcher struct {
	logger *zap.Logger

	// dydxPricesFetcher is the fetcher for the dydx x/prices market map.
	pricesFetcher mmclient.MarketMapFetcher
	// marketmapFetcher is the fetcher for the x/marketmap market map.
	marketmapFetcher mmclient.MarketMapFetcher
	// switched is true if the fetcher has switched over to the x/marketmap API.
	switched bool
	metrics  metrics.APIMetrics
}

// NewDefaultSwitchOverMarketMapFetcher returns a new SwitchOverProvider with the default
// dYdX x/prices and x/marketmap fetchers.
func NewDefaultSwitchOverMarketMapFetcher(
	logger *zap.Logger,
	api config.APIConfig,
	rh apihandlers.RequestHandler,
	metrics metrics.APIMetrics,
) (mmclient.MarketMapFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}
	if api.Name != SwitchOverAPIHandlerName {
		return nil, fmt.Errorf("expected api name %s, got %s", SwitchOverAPIHandlerName, api.Name)
	}
	if len(api.Endpoints) != 2 {
		return nil, fmt.Errorf(
			"expected two endpoints, got %d",
			len(api.Endpoints),
		)
	}
	if rh == nil {
		return nil, fmt.Errorf("request handler is nil")
	}
	if metrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}

	// Construct the dYdX x/prices API handler.
	pricesAPIHandler, err := NewAPIHandler(logger, api)
	if err != nil {
		return nil, err
	}
	pricesFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		pricesAPIHandler,
		metrics,
		api,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// Construct the dYdX x/marketmap API handler.
	conn, err := connectgrpc.NewClient(
		api.Endpoints[1].URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
	)
	if err != nil {
		return nil, err
	}
	marketmapClient, err := marketmap.NewGRPCClientWithConn(conn, api, metrics)
	if err != nil {
		return nil, err
	}
	marketmapFetcher, err := marketmap.NewMarketMapFetcherWithClient(logger, marketmapClient)
	if err != nil {
		return nil, err
	}

	return &SwitchOverFetcher{
		logger:           logger.With(zap.String("api", api.Name)),
		pricesFetcher:    pricesFetcher,
		marketmapFetcher: marketmapFetcher,
		metrics:          metrics,
	}, nil
}

// NewSwitchOverFetcher returns a new SwitchOverFetcher with the provided
// fetchers.
func NewSwitchOverFetcher(
	logger *zap.Logger,
	pricesFetcher mmclient.MarketMapFetcher,
	marketmapFetcher mmclient.MarketMapFetcher,
	metrics metrics.APIMetrics,
) (mmclient.MarketMapFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if pricesFetcher == nil {
		return nil, fmt.Errorf("prices fetcher is nil")
	}

	if marketmapFetcher == nil {
		return nil, fmt.Errorf("marketmap fetcher is nil")
	}

	return &SwitchOverFetcher{
		logger:           logger,
		pricesFetcher:    pricesFetcher,
		marketmapFetcher: marketmapFetcher,
		metrics:          metrics,
	}, nil
}

// Fetch fetches the market map from the x/prices API and then the x/marketmap
// API. The fetcher will switch over to the x/marketmap API the first time it
// receives a non-nil market map from the x/marketmap API.
func (f *SwitchOverFetcher) Fetch(
	ctx context.Context,
	chains []mmclient.Chain,
) mmclient.MarketMapResponse {
	var resp mmclient.MarketMapResponse

	// we need to emit a ProviderResponse (if any of the underlying fetches were successful)
	// to indicate which provider is currently in use
	defer func() {
		provider := Name
		if f.switched {
			provider = marketmap.Name
		}

		for id := range resp.Resolved {
			f.metrics.AddProviderResponse(provider, strings.ToLower(id.String()), providertypes.OK)
		}

		for id, result := range resp.UnResolved {
			f.metrics.AddProviderResponse(provider, strings.ToLower(id.String()), result.Code())
		}
	}()

	if f.switched {
		resp = f.marketmapFetcher.Fetch(ctx, chains)
		return resp
	}

	f.logger.Info("fetching marketmap from x/marketmap")
	resp = f.marketmapFetcher.Fetch(ctx, chains)
	if len(resp.Resolved) > 0 {
		f.logger.Info("got response from x/marketmap; switching over to x/marketmap")
		f.switched = true
		return resp
	}

	f.logger.Info("x/marketmap did not return a marketmap response; fetching marketmap from x/prices")
	resp = f.pricesFetcher.Fetch(ctx, chains)
	return resp
}
