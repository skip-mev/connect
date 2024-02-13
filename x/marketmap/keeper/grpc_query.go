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
	return &types.GetMarketMapResponse{
			MarketMap:   *mm,
			LastUpdated: 0,
		},
		err
}
