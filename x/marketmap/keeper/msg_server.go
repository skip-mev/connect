package keeper

import (
	"context"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

// msgServer is the default implementation of the x/marketmap MsgService.
type msgServer struct {
	k Keeper
}

// NewMsgServer returns the default implementation of the x/marketmap message service.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = (*msgServer)(nil)

// CreateMarket creates a market from the given message.
func (ms msgServer) CreateMarket(_ context.Context, _ *types.MsgCreateMarket) (*types.MsgCreateMarketResponse, error) {
	// TODO finish

	return nil, nil
}

// Params updates the x/marketmap module's Params.
func (ms msgServer) Params(_ context.Context, _ *types.MsgParams) (*types.MsgParamsResponse, error) {
	// TODO finish

	return nil, nil
}
