package keeper

import (
	"github.com/skip-mev/slinky/x/marketmap/types"
)

// Hooks gets the hooks for x/marketmap keeper.
func (k *Keeper) Hooks() types.MarketMapHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiMarketMapHooks{}
	}

	return k.hooks
}

// SetHooks sets the x/marketmap hooks.
func (k *Keeper) SetHooks(mmh types.MarketMapHooks) {
	if k.hooks != nil {
		panic("cannot set marketmap hooks twice")
	}

	k.hooks = mmh
}
