package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

type queryServerImpl struct {
	k Keeper
}

// NewQueryServer returns an implementation of the x/marketmap QueryServer.
func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServerImpl{k}
}

// GetMarketMap returns the full AggregateMarketConfig stored in the x/marketmap module.
func (q queryServerImpl) GetMarketMap(goCtx context.Context, req *types.GetMarketMapRequest) (*types.GetMarketMapResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	mm, err := q.k.GetMarketMap(ctx)
	if err != nil {
		return nil, err
	}

	params, err := q.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	lastUpdated, err := q.k.GetLastUpdated(ctx)
	return &types.GetMarketMapResponse{
			MarketMap:   *mm,
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

	return &types.ParamsResponse{Params: params}, err
}
