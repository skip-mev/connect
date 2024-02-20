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

	lastUpdated, err := q.k.GetLastUpdated(ctx)
	return &types.GetMarketMapResponse{
			MarketMap: types.MarketMap{
				Tickers:   tickers,
				Paths:     paths,
				Providers: providers,
			},
			LastUpdated: lastUpdated,
		},
		err
}
