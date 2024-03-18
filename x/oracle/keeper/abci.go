package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block.  It flushes the count of
// removed currency pairs to state and rests the counter.
func (k *Keeper) EndBlocker(goCtx context.Context) error {
	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// flush to state
	err := k.numRemoves.Set(ctx, k.removedCPsCounter)
	if err != nil {
		return err
	}

	// reset value
	k.removedCPsCounter = 0

	return nil
}
