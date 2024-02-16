package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MarketMapHooks is the interface that defines the hooks that can be integrated by other modules.
type MarketMapHooks interface {
	AfterMarketCreated(ctx sdk.Context) error

	AfterMarketUpdated(ctx sdk.Context) error
}

var _ MarketMapHooks = &MultiMarketMapHooks{}

// MultiMarketMapHooks defines an array of MarketMapHooks which can be executed in sequence.
type MultiMarketMapHooks []MarketMapHooks

// NewMultiMarketMapHooks creates a MultiMarketMapHooks object from a variadic amount of MarketMapHooks.
func NewMultiMarketMapHooks(hooks ...MarketMapHooks) MultiMarketMapHooks {
	return hooks
}

// AfterMarketCreated calls all AfterMarketCreated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketCreated(ctx sdk.Context) error {
	for i := range mh {
		if err := mh[i].AfterMarketCreated(ctx); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketUpdated calls all AfterMarketUpdated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketUpdated(ctx sdk.Context) error {
	for i := range mh {
		if err := mh[i].AfterMarketUpdated(ctx); err != nil {
			return err
		}
	}

	return nil
}

// MarketMapHooksWrapper is a wrapper for modules to inject MarketMapHooks using depinject.
type MarketMapHooksWrapper struct{ MarketMapHooks }
