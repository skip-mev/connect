package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MarketMapHooks is the interface that defines the hooks that can be integrated by other modules.
type MarketMapHooks interface {
	// AfterMarketCreated_ is called after CreateMarket is called.
	AfterMarketCreated_(ctx sdk.Context, market Market) error

	// AfterMarketUpdated_ is called after UpdateMarket is called.
	AfterMarketUpdated_(ctx sdk.Context, market Market) error

	// AfterMarketGenesis_ is called after x/marketmap init genesis.
	AfterMarketGenesis_(ctx sdk.Context, tickers map[string]Market) error
}

var _ MarketMapHooks = &MultiMarketMapHooks{}

// MultiMarketMapHooks defines an array of MarketMapHooks which can be executed in sequence.
type MultiMarketMapHooks []MarketMapHooks

// AfterMarketCreated_ calls all AfterMarketCreated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketCreated_(ctx sdk.Context, market Market) error {
	for i := range mh {
		if err := mh[i].AfterMarketCreated_(ctx, market); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketUpdated_ calls all AfterMarketUpdated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketUpdated_(ctx sdk.Context, market Market) error {
	for i := range mh {
		if err := mh[i].AfterMarketUpdated_(ctx, market); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketGenesis_ calls all AfterMarketGenesis hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketGenesis_(ctx sdk.Context, markets map[string]Market) error {
	for i := range mh {
		if err := mh[i].AfterMarketGenesis_(ctx, markets); err != nil {
			return err
		}
	}

	return nil
}

// MarketMapHooksWrapper is a wrapper for modules to inject MarketMapHooks using depinject.
type MarketMapHooksWrapper struct{ MarketMapHooks }
