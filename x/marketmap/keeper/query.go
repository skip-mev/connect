package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

type queryServerImpl struct {
	k *Keeper
}

// NewQueryServer returns an implementation of the x/marketmap QueryServer.
func NewQueryServer(k *Keeper) types.QueryServer {
	return &queryServerImpl{k}
}

// MarketMap returns the full MarketMap and associated information stored in the x/marketmap module.
func (q queryServerImpl) MarketMap(goCtx context.Context, req *types.MarketMapRequest) (*types.MarketMapResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	markets, err := q.k.GetAllMarkets(ctx)
	if err != nil {
		return nil, err
	}

	lastUpdated, err := q.k.GetLastUpdated(ctx)
	return &types.MarketMapResponse{
			MarketMap: types.MarketMap{
				Markets: markets,
			},

			LastUpdated: lastUpdated,
			ChainId:     ctx.ChainID(),
		},
		err
}

// Market returns the requested market stored in the x/marketmap module.
func (q queryServerImpl) Market(goCtx context.Context, req *types.MarketRequest) (*types.MarketResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if err := req.CurrencyPair.ValidateBasic(); err != nil {
		return nil, err
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := q.k.GetMarket(ctx, req.CurrencyPair.String())
	if err != nil {
		return nil, err
	}

	return &types.MarketResponse{Market: market}, nil
}

// LastUpdated returns the last height the marketmap was updated in the x/marketmap module.
func (q queryServerImpl) LastUpdated(goCtx context.Context, req *types.LastUpdatedRequest) (*types.LastUpdatedResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	lastUpdated, err := q.k.lastUpdated.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.LastUpdatedResponse{LastUpdated: lastUpdated}, nil
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
