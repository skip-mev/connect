package keeper

import "github.com/skip-mev/connect/v2/x/marketmap/types"

// Option is a type that modifies a keeper during instantiation.  These can be passed variadically into NewKeeper
// to specify keeper behavior.
type Option func(*Keeper)

// WithHooks sets the keeper hooks to the given hooks.
func WithHooks(hooks types.MarketMapHooks) Option {
	return func(k *Keeper) {
		k.hooks = hooks
	}
}

// WithDeleteValidationHooks sets the keeper deleteMarketValidationHooks to the given hooks.
func WithDeleteValidationHooks(hooks []types.MarketValidationHook) Option {
	return func(k *Keeper) {
		k.deleteMarketValidationHooks = hooks
	}
}
