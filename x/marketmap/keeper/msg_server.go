package keeper

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
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
	for _, market := range msg.CreateMarkets {
		err = ms.k.CreateMarket(ctx, market)
		if err != nil {
			return nil, err
		}

		err = ms.k.hooks.AfterMarketCreated(ctx, market)
		if err != nil {
			return nil, fmt.Errorf("unable to run create market hook: %w", err)
		}

		event := sdk.NewEvent(
			types.EventTypeCreateMarket,
			sdk.NewAttribute(types.AttributeKeyCurrencyPair, market.Ticker.String()),
			sdk.NewAttribute(types.AttributeKeyDecimals, strconv.FormatUint(market.Ticker.Decimals, 10)),
			sdk.NewAttribute(types.AttributeKeyMinProviderCount, strconv.FormatUint(market.Ticker.MinProviderCount, 10)),
			sdk.NewAttribute(types.AttributeKeyMetadata, market.Ticker.Metadata_JSON),
		)
		ctx.EventManager().EmitEvent(event)
	}

	// validate that the new state of the marketmap is valid
	err = ms.k.ValidateState(ctx, msg.CreateMarkets)
	if err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgCreateMarketsResponse{}, ms.k.SetLastUpdated(ctx, uint64(ctx.BlockHeight()))
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
			return nil, fmt.Errorf("unable to update market: %w", err)
		}

		err = ms.k.hooks.AfterMarketUpdated(ctx, market)
		if err != nil {
			return nil, fmt.Errorf("unable to run update market hook: %w", err)
		}

		event := sdk.NewEvent(
			types.EventTypeUpdateMarket,
			sdk.NewAttribute(types.AttributeKeyCurrencyPair, market.Ticker.String()),
			sdk.NewAttribute(types.AttributeKeyDecimals, strconv.FormatUint(market.Ticker.Decimals, 10)),
			sdk.NewAttribute(types.AttributeKeyMinProviderCount, strconv.FormatUint(market.Ticker.MinProviderCount, 10)),
			sdk.NewAttribute(types.AttributeKeyMetadata, market.Ticker.Metadata_JSON),
		)
		ctx.EventManager().EmitEvent(event)

	}

	// validate that the new state of the marketmap is valid
	err = ms.k.ValidateState(ctx, msg.UpdateMarkets)
	if err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgUpdateMarketsResponse{}, ms.k.SetLastUpdated(ctx, uint64(ctx.BlockHeight()))
}

// UpdateParams updates the x/marketmap module's Params.
func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgParams) (*types.MsgParamsResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != ms.k.authority.String() {
		return nil, fmt.Errorf("request authority %s does not match module keeper authority %s", msg.Authority, ms.k.authority.String())
	}

	if err := ms.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgParamsResponse{}, nil
}

func (ms msgServer) RemoveMarketAuthorities(goCtx context.Context, msg *types.MsgRemoveMarketAuthorities) (*types.MsgRemoveMarketAuthoritiesResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := ms.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Admin != params.Admin {
		return nil, fmt.Errorf("request admin %s does not match module admin %s", msg.Admin, params.Admin)
	}

	if len(msg.RemoveAddresses) > len(params.MarketAuthorities) {
		return nil, fmt.Errorf("remove addresses must be a subset of the current market authorities")
	}

	removeAddresses := make(map[string]struct{}, len(msg.RemoveAddresses))
	for _, remove := range msg.RemoveAddresses {
		removeAddresses[remove] = struct{}{}
	}

	for i, address := range params.MarketAuthorities {
		if _, found := removeAddresses[address]; found {
			params.MarketAuthorities = slices.Delete(params.MarketAuthorities, i, i+1)
		}
	}

	if err := ms.k.SetParams(ctx, params); err != nil {
		return nil, err
	}

	return &types.MsgRemoveMarketAuthoritiesResponse{}, nil
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
