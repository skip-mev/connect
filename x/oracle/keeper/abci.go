package keeper

import (
	"context"
)

// BeginBlocker is called at the beginning of every block.  It resets the count of
// removed currency pairs.
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	return k.numRemoves.Set(ctx, 0)
}
