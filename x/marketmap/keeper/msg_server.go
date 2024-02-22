package keeper

import (
	"context"
	"fmt"
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

// UpdateMarketMap updates the marketmap from the given message.  All updates are made to the market map and then
// the resulting final state is checked to verify that the end state is valid.
func (ms msgServer) UpdateMarketMap(goCtx context.Context, msg *types.MsgUpdateMarketMap) (*types.MsgUpdateMarketMapResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

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

	// create markets
	for _, market := range msg.CreateMarkets {
		err := ms.k.CreateMarket(ctx, market.Ticker, market.Paths, market.Providers)
		if err != nil {
			return nil, err
		}

		// TODO: call creation hooks

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeCreateMarket,
				sdk.NewAttribute(types.AttributeKeyCurrencyPair, market.Ticker.String()),
				sdk.NewAttribute(types.AttributeKeyDecimals, strconv.FormatUint(market.Ticker.Decimals, 10)),
				sdk.NewAttribute(types.AttributeKeyMinProviderCount, strconv.FormatUint(market.Ticker.MinProviderCount, 10)),
				sdk.NewAttribute(types.AttributeKeyMetadata, market.Ticker.Metadata_JSON),
				sdk.NewAttribute(types.AttributeKeyPaths, market.Paths.String()),
				sdk.NewAttribute(types.AttributeKeyProviders, market.Providers.String()),
			),
		)
	}

	// update markets
	// TODO

	// validate that the new state of the marketmap is valid
	err := ms.k.ValidateState(ctx, msg.CreateMarkets)
	if err != nil {
		return nil, fmt.Errorf("invalid state resulting from update: %w", err)
	}

	return &types.MsgUpdateMarketMapResponse{}, nil
}
