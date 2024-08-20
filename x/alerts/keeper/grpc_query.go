package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

type queryServerImpl struct {
	k Keeper
}

// NewQueryServer returns an implementation of the x/alerts QueryServer.
func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServerImpl{k}
}

// Alerts returns all Alerts that match the provided alert status, if no alert-status is given, all alerts
// in module state will be returned.
func (q queryServerImpl) Alerts(srvCtx context.Context, req *types.AlertsRequest) (*types.AlertsResponse, error) {
	// if the request is nil, error
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(srvCtx)

	// get all alerts from state
	alerts, err := q.k.GetAllAlertsWithCondition(ctx, func(a types.AlertWithStatus) bool {
		switch req.Status {
		case types.AlertStatusID_CONCLUSION_STATUS_CONCLUDED:
			return a.Status.ConclusionStatus == uint64(types.Concluded)
		case types.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED:
			return a.Status.ConclusionStatus == uint64(types.Unconcluded)
		case types.AlertStatusID_CONCLUSION_STATUS_UNSPECIFIED:
			return true
		default:
			return false
		}
	})
	if err != nil {
		return nil, err
	}

	// convert to alerts
	res := &types.AlertsResponse{}
	res.Alerts = alertWithStatusToAlerts(alerts)
	return res, nil
}

func alertWithStatusToAlerts(as []types.AlertWithStatus) []types.Alert {
	alerts := make([]types.Alert, len(as))
	for i, a := range as {
		alerts[i] = a.Alert
	}
	return alerts
}

// Params returns the current module params for x/alerts.
func (q queryServerImpl) Params(srvCtx context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	// if the request is nil, error
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(srvCtx)

	// get the params from state
	params := q.k.GetParams(ctx)

	// convert to response
	return &types.ParamsResponse{
		Params: params,
	}, nil
}
