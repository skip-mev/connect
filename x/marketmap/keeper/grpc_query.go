package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

type queryServerImpl struct {
	k *Keeper
}

// NewQueryServer returns an implementation of the x/marketmap QueryServer.
func NewQueryServer(k *Keeper) types.QueryServer {
	return &queryServerImpl{k}
}

// MarketMap returns the full MarketMap and associated information stored in the x/marketmap module.
func (q queryServerImpl) MarketMap(goCtx context.Context, req *types.GetMarketMapRequest) (*types.GetMarketMapResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	tickers, err := q.k.GetAllTickersMap(ctx)
	if err != nil {
		return nil, err
	}

	paths, err := q.k.GetAllPathsMap(ctx)
	if err != nil {
		return nil, err
	}

	providers, err := q.k.GetAllProvidersMap(ctx)
	if err != nil {
		return nil, err
	}

	params, err := q.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	lastUpdated, err := q.k.GetLastUpdated(ctx)
	return &types.GetMarketMapResponse{
			MarketMap: types.MarketMap{
				Tickers:   tickers,
				Paths:     paths,
				Providers: providers,
			},
			LastUpdated: lastUpdated,
			Version:     params.Version,
		},
		err
}

// Params returns the parameters stored in the x/marketmap module.
func (q queryServerImpl) Params(goCtx context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := q.k.params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.ParamsResponse{Params: params}, nil
}

// LastUpdated returns the last height the marketmap was updated in the x/marketmap module.
func (q queryServerImpl) LastUpdated(goCtx context.Context, req *types.GetLastUpdatedRequest) (*types.GetLastUpdatedResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	lastUpdated, err := q.k.lastUpdated.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.GetLastUpdatedResponse{LastUpdated: lastUpdated}, nil
}
