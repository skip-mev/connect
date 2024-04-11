package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block.  It resets the count of
// removed currency pairs.
func (k *Keeper) BeginBlocker(goCtx context.Context) error {
	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.numRemoves.Set(ctx, 0)
}
