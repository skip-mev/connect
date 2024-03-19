package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block.  It resets the count of
// removed currency pairs.
func (k *Keeper) BeginBlocker(goCtx context.Context) error {
	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	removes, err := k.numRemoves.Get(ctx)
	if err != nil {
		return err
	}

	numCPs, err := k.numCPs.Get(ctx)
	if err != nil {
		return err
	}

	if numCPs < removes {
		return fmt.Errorf("invalid decrement amount - result will be negative")
	}

	err = k.numCPs.Set(ctx, numCPs-removes)
	if err != nil {
		return err
	}

	return k.numRemoves.Set(ctx, 0)
}
