package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/mm2/types"
)

// msgServer is the default implementation of the x/marketmap MsgService.
type msgServer struct {
	k *Keeper
}

// NewMsgServer returns the default implementation of the x/marketmap message service.
func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = (*msgServer)(nil)

// CreateMarkets updates the marketmap by creating markets from the given message.  All updates are made to the market
// map and then the resulting final state is checked to verify that the end state is valid.
func (ms msgServer) CreateMarkets(goCtx context.Context, msg *types.MsgCreateMarkets) (*types.MsgCreateMarketsResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := ms.k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get marketmap params: %w", err)
	}

	found := checkMarketAuthority(msg.Authority, params)
	if !found {
		return nil, fmt.Errorf("request signer %s does not match module market authorities", msg.Authority)
	}

	// create markets
	for _, createMarket := range msg.CreateMarkets {
		err = ms.k.CreateMarket(ctx, createMarket)
		if err != nil {
			return nil, err
		}

		// TODO hooks

		// TODO events
	}

	// validate that the new state of the marketmap is valid
	err = ms.k.ValidateState(ctx, msg.CreateMarkets, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgCreateMarketsResponse{}, nil
}

// UpdateMarkets updates the marketmap by updating markets from the given message.  All updates are made to the market
// map and then the resulting final state is checked to verify that the end state is valid.
func (ms msgServer) UpdateMarkets(goCtx context.Context, msg *types.MsgUpdateMarkets) (*types.MsgUpdateMarketsResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := ms.k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get marketmap params: %w", err)
	}

	found := checkMarketAuthority(msg.Authority, params)
	if !found {
		return nil, fmt.Errorf("request signer %s does not match module market authorities", msg.Authority)
	}

	for _, market := range msg.UpdateMarkets {
		err = ms.k.UpdateMarket(ctx, market)
		if err != nil {
			return nil, err
		}

		// TODO hooks

		// TODO events
	}

	// validate that the new state of the marketmap is valid
	err = ms.k.ValidateState(ctx, nil, msg.UpdateMarkets)
	if err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgUpdateMarketsResponse{}, nil
}

// Params updates the x/marketmap module's Params.
func (ms msgServer) Params(goCtx context.Context, msg *types.MsgParams) (*types.MsgParamsResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	// Update the module's parameters.
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := ms.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Authority != ms.k.authority.String() {
		return nil, fmt.Errorf("request authority %s does not match module keeper authority %s", msg.Authority, ms.k.authority.String())
	}

	if msg.Params.Version < params.Version {
		return nil, fmt.Errorf("request version %d is less than current params version %d", msg.Params.Version, params.Version)
	}

	if err := ms.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgParamsResponse{}, nil
}

// checkMarketAuthority checks if the given authority is the x/marketmap's list of MarketAuthorities.
func checkMarketAuthority(authority string, params types.Params) bool {
	if len(params.MarketAuthorities) == 0 {
		return false
	}

	for _, auth := range params.MarketAuthorities {
		if authority == auth {
			return true
		}
	}

	return false
}
