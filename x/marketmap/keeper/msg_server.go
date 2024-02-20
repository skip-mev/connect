package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

// msgServer is the default implementation of the x/marketmap MsgService.
type msgServer struct {
	k Keeper
}

// NewMsgServer returns the default implementation of the x/marketmap message service.
// NewMsgServer returns the default implementation of the x/marketmap message service.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = (*msgServer)(nil)

// CreateMarket creates a market from the given message.
func (ms msgServer) CreateMarket(goCtx context.Context, msg *types.MsgCreateMarket) (*types.MsgCreateMarketResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	// Update the module's parameters.
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: add check when params are added
	// params, err := ms.k.GetParams(ctx)
	// if err != nil {
	//	return nil, fmt.Errorf("unable to get marketmap params: %w", err)
	// }
	//
	// if msg.Signer != params.MarketAuthority {
	//	return nil, fmt.Errorf("request signer %s does not match module market authority %s", msg.Signer, params.MarketAuthority)
	// }

	err := ms.k.CreateMarket(ctx, msg.Ticker, msg.Paths, msg.Providers)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMarketResponse{}, nil
}
