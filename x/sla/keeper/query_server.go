package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

var _ slatypes.QueryServer = (*QueryServer)(nil)

// QueryServer defines the gRPC server for the x/sla module.
type QueryServer struct {
	k Keeper
}

// NewQueryServer creates a new instance of the x/sla QueryServer type.
func NewQueryServer(keeper Keeper) slatypes.QueryServer {
	return &QueryServer{k: keeper}
}

// GetAllSLAs defines a method that returns all SLAs in the store.
func (s *QueryServer) GetAllSLAs(goCtx context.Context, _ *slatypes.GetAllSLAsRequest) (*slatypes.GetAllSLAsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	slas, err := s.k.GetSLAs(ctx)
	if err != nil {
		return nil, err
	}

	return &slatypes.GetAllSLAsResponse{SLAs: slas}, nil
}

// GetPriceFeeds defines a method that returns all price feeds in the store with the given SLA ID.
func (s *QueryServer) GetPriceFeeds(goCtx context.Context, req *slatypes.GetPriceFeedsRequest) (*slatypes.GetPriceFeedsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeds, err := s.k.GetAllPriceFeeds(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &slatypes.GetPriceFeedsResponse{PriceFeeds: feeds}, nil
}

// Params defines a method that returns the current SLA parameters.
func (s *QueryServer) Params(goCtx context.Context, _ *slatypes.ParamsRequest) (*slatypes.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := s.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &slatypes.ParamsResponse{Params: params}, nil
}
