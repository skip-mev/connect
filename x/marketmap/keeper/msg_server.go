package keeper

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
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

// UpsertMarkets wraps both Create / Update markets into a single message. Specifically
// if a market does not exist it will be created, otherwise it will be updated. The response
// will be a map between ticker -> updated.
func (ms msgServer) UpsertMarkets(goCtx context.Context, msg *types.MsgUpsertMarkets) (*types.MsgUpsertMarketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// perform basic msg validity checks
	if err := ms.verifyMarketAuthorities(ctx, msg); err != nil {
		return nil, fmt.Errorf("unable to verify market authorities: %w", err)
	}

	// iterate over all markets and either create them (if no market exists), or update them
	for _, market := range msg.Markets {
		// check if market exists
		exists, err := ms.k.HasMarket(ctx, market.Ticker.String())
		if err != nil {
			return nil, err
		}

		var eventType string
		// if market does not exist, create it
		if !exists {
			err = ms.k.CreateMarket(ctx, market)
			if err != nil {
				return nil, err
			}

			// run hooks
			if err = ms.k.hooks.AfterMarketCreated(ctx, market); err != nil {
				return nil, err
			}

			eventType = types.EventTypeCreateMarket
		} else {
			err = ms.k.UpdateMarket(ctx, market)
			if err != nil {
				return nil, err
			}

			// run hooks
			if err = ms.k.hooks.AfterMarketUpdated(ctx, market); err != nil {
				return nil, err
			}

			eventType = types.EventTypeUpdateMarket
		}

		event := sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyCurrencyPair, market.Ticker.String()),
			sdk.NewAttribute(types.AttributeKeyDecimals, strconv.FormatUint(market.Ticker.Decimals, 10)),
			sdk.NewAttribute(types.AttributeKeyMinProviderCount, strconv.FormatUint(market.Ticker.MinProviderCount, 10)),
			sdk.NewAttribute(types.AttributeKeyMetadata, market.Ticker.Metadata_JSON),
		)
		ctx.EventManager().EmitEvent(event)
	}

	// validate that the new state of the marketmap is valid
	if err := ms.k.ValidateState(ctx, msg.Markets); err != nil {
		return nil, err
	}

	return &types.MsgUpsertMarketsResponse{}, ms.k.SetLastUpdated(ctx, uint64(ctx.BlockHeight())) //nolint:gosec
}

// CreateMarkets updates the marketmap by creating markets from the given message.  All updates are made to the market
// map and then the resulting final state is checked to verify that the end state is valid.
func (ms msgServer) CreateMarkets(goCtx context.Context, msg *types.MsgCreateMarkets) (*types.MsgCreateMarketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// perform basic msg validity checks
	if err := ms.verifyMarketAuthorities(ctx, msg); err != nil {
		return nil, fmt.Errorf("unable to verify market authorities: %w", err)
	}

	// create markets
	for _, market := range msg.CreateMarkets {
		err := ms.k.CreateMarket(ctx, market)
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
	err := ms.k.ValidateState(ctx, msg.CreateMarkets)
	if err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgCreateMarketsResponse{}, ms.k.SetLastUpdated(ctx, uint64(ctx.BlockHeight())) //nolint:gosec
}

// UpdateMarkets updates the marketmap by updating markets from the given message.  All updates are made to the market
// map and then the resulting final state is checked to verify that the end state is valid.
func (ms msgServer) UpdateMarkets(goCtx context.Context, msg *types.MsgUpdateMarkets) (*types.MsgUpdateMarketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// perform basic msg validity checks
	if err := ms.verifyMarketAuthorities(ctx, msg); err != nil {
		return nil, fmt.Errorf("unable to verify market authorities: %w", err)
	}

	for _, market := range msg.UpdateMarkets {
		err := ms.k.UpdateMarket(ctx, market)
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
	if err := ms.k.ValidateState(ctx, msg.UpdateMarkets); err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgUpdateMarketsResponse{}, ms.k.SetLastUpdated(ctx, uint64(ctx.BlockHeight())) //nolint:gosec
}

// verifyMarketAuthorities verifies that the msg-submitter is a market-authority
// and returns the context for the msg, this method returns an error if the submitter is not a market
// authority.
func (ms msgServer) verifyMarketAuthorities(ctx sdk.Context, msg interface {
	GetAuthority() string
},
) error {
	if msg == nil {
		return fmt.Errorf("unable to process nil msg")
	}

	params, err := ms.k.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("unable to get marketmap params: %w", err)
	}

	found := checkMarketAuthority(msg.GetAuthority(), params)
	if !found {
		return fmt.Errorf("request signer %s does not match module market authorities", msg.GetAuthority())
	}

	return nil
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
