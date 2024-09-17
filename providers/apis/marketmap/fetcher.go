package marketmap

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// MarketMapFetcher is the x/marketmap fetcher. This fetcher is responsible for querying the
// x/marketmap module and returning the market map data. The fetcher utilizes the QueryClient
// to query the x/marketmap module.
type MarketMapFetcher struct { //nolint
	logger *zap.Logger

	// client is the QueryClient implementation. This is used to interact with the x/marketmap
	// module.
	client mmtypes.QueryClient
}

// NewMarketMapFetcher returns a new MarketMap fetcher with the standard grpc client.
func NewMarketMapFetcher(
	logger *zap.Logger,
	api config.APIConfig,
	metrics metrics.APIMetrics,
) (*MarketMapFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid api name; expected %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics is required")
	}

	client, err := NewGRPCClient(api, metrics)
	if err != nil {
		return nil, err
	}

	return &MarketMapFetcher{
		logger: logger.With(zap.String("fetcher", Name)),
		client: client,
	}, nil
}

// NewMarketMapFetcherWithClient returns a new MarketMap fetcher.
func NewMarketMapFetcherWithClient(
	logger *zap.Logger,
	client mmtypes.QueryClient,
) (*MarketMapFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if client == nil {
		return nil, fmt.Errorf("client is required")
	}

	return &MarketMapFetcher{
		logger: logger.With(zap.String("fetcher", Name)),
		client: client,
	}, nil
}

// Fetch returns the latest market map data from the x/marketmap module. It expects only a single
// chain ID since the current implementation assumes a single connection to one chain.
func (f *MarketMapFetcher) Fetch(
	ctx context.Context,
	chains []types.Chain,
) types.MarketMapResponse {
	if len(chains) != 1 {
		f.logger.Info("expected one chain, got multiple chains", zap.Any("chains", chains))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected one chain, got %d", len(chains)),
				providertypes.ErrorInvalidAPIChains,
			),
		)
	}

	// Query the x/marketmap module for the market map data.
	resp, err := f.client.MarketMap(ctx, &mmtypes.MarketMapRequest{})
	if err != nil {
		f.logger.Error("failed to query market map module on node", zap.Error(err))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to query market map: %w", err),
				providertypes.ErrorGRPCGeneral,
			),
		)
	}

	if resp == nil {
		f.logger.Info("nil response from market map module query")
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("nil response from market map query"),
				providertypes.ErrorGRPCGeneral,
			),
		)
	}

	resolved := make(types.ResolvedMarketMap)
	resolved[chains[0]] = types.NewMarketMapResult(resp, time.Now())

	f.logger.Info("successfully fetched market map data from module; checking if market map has changed")
	return types.NewMarketMapResponse(resolved, nil)
}
