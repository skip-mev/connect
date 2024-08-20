package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

var _ types.QueryServer = QueryServer{}

// QueryServer is the default implementation of the x/incentives QueryServer.
type QueryServer struct {
	k Keeper
}

// NewQueryServer returns an implementation of the x/incentives QueryServer.
func NewQueryServer(k Keeper) QueryServer {
	return QueryServer{
		k,
	}
}

// GetIncentivesByType returns all incentives of a given type currently stored in the
// incentives module. If the type is not registered with the module, an error is returned.
func (q QueryServer) GetIncentivesByType(
	ctx context.Context,
	req *types.GetIncentivesByTypeRequest,
) (*types.GetIncentivesByTypeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	var incentive types.Incentive
	for i := range q.k.incentiveStrategies {
		if i.Type() == req.IncentiveType {
			incentive = i
			break
		}
	}

	if incentive == nil {
		return nil, fmt.Errorf("unknown incentive type: %s", req.IncentiveType)
	}

	incentives, err := q.k.GetIncentivesByType(sdk.UnwrapSDKContext(ctx), incentive)
	if err != nil {
		return nil, err
	}

	incentiveBytes, err := types.IncentivesToBytes(incentives...)
	if err != nil {
		return nil, err
	}

	resp := types.GetIncentivesByTypeResponse{
		Entries: incentiveBytes,
	}

	return &resp, nil
}

// GetAllIncentives returns all incentives currently stored in the module.
func (q QueryServer) GetAllIncentives(
	ctx context.Context,
	_ *types.GetAllIncentivesRequest,
) (*types.GetAllIncentivesResponse, error) {
	incentives := make([]types.IncentivesByType, 0)

	// Get all of the incentive types and sort them by name.
	sortedIncentives := types.SortIncentivesStrategiesMap(q.k.incentiveStrategies)

	for _, incentive := range sortedIncentives {
		incentivesByType, err := q.k.GetIncentivesByType(sdk.UnwrapSDKContext(ctx), incentive)
		if err != nil {
			return nil, err
		}

		if len(incentivesByType) == 0 {
			continue
		}

		incentiveBytes, err := types.IncentivesToBytes(incentivesByType...)
		if err != nil {
			return nil, err
		}

		incentives = append(incentives, types.NewIncentives(incentive.Type(), incentiveBytes))
	}

	resp := types.GetAllIncentivesResponse{
		Registry: incentives,
	}

	return &resp, nil
}
