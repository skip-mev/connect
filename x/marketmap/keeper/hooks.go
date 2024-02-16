package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.MarketMapHooks = Hooks{}

// Hooks gets the hooks for x/marketmap keeper {
func (k *Keeper) Hooks() types.MarketMapHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiMarketMapHooks{}
	}

	return k.hooks
}

// SetHooks sets the x/marketmap hooks.  In contrast to other receivers, this method must take a pointer due to nature
// of the hooks interface and SDK start up sequence.
func (k *Keeper) SetHooks(mmh types.MarketMapHooks) {
	if k.hooks != nil {
		panic("cannot set marketmap hooks twice")
	}

	k.hooks = mmh
}

// Hooks wrapper struct for x/marketmap keeper
type Hooks struct {
	k *Keeper
}

func (h Hooks) AfterMarketCreated(_ sdk.Context) error {
	return nil
}

func (h Hooks) AfterMarketUpdated(_ sdk.Context) error {
	return nil
}
