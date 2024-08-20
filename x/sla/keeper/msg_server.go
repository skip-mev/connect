package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

var _ slatypes.MsgServer = (*MsgServer)(nil)

// Request is a basic interface for all messages.
type Request interface {
	ValidateBasic() error
	GetAuthority() string
}

// MsgServer is the server API for x/sla Msg service.
type MsgServer struct {
	k Keeper
}

// NewMsgServer returns the MsgServer implementation.
func NewMsgServer(k Keeper) slatypes.MsgServer {
	return &MsgServer{k}
}

// AddSLAs defines a method that adds a set of SLAs to the store. The SLAs provided must not already
// exist in the store and must be valid. The signer of the message must also be the module authority.
func (m *MsgServer) AddSLAs(goCtx context.Context, req *slatypes.MsgAddSLAs) (*slatypes.MsgAddSLAsResponse, error) {
	// Add all SLAs in message to state.
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.AddSLAs(ctx, req.SLAs); err != nil {
		return nil, err
	}

	return &slatypes.MsgAddSLAsResponse{}, nil
}

// RemoveSLAs defines a method that removes a set of SLAs from the store. The SLAs provided must
// exist in the store, and the signer of the message must be the module authority.
func (m *MsgServer) RemoveSLAs(goCtx context.Context, req *slatypes.MsgRemoveSLAs) (*slatypes.MsgRemoveSLAsResponse, error) {
	// Remove the SLAs in message from state.
	ctx := sdk.UnwrapSDKContext(goCtx)
	for _, id := range req.IDs {
		if err := m.k.RemovePriceFeedsBySLA(ctx, id); err != nil {
			return nil, err
		}

		if err := m.k.RemoveSLA(ctx, id); err != nil {
			return nil, err
		}
	}

	return &slatypes.MsgRemoveSLAsResponse{}, nil
}

// Params defines a method that updates the module's parameters. The signer of the message must
// be the module authority.
func (m *MsgServer) Params(goCtx context.Context, req *slatypes.MsgParams) (*slatypes.MsgParamsResponse, error) {
	// Update the module's parameters.
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &slatypes.MsgParamsResponse{}, nil
}
