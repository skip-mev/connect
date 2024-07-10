package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MarketMapHooks is the interface that defines the hooks that can be integrated by other modules.
//
//go:generate mockery --name MarketMapHooks
type MarketMapHooks interface {
	// AfterMarketCreated is called after CreateMarket is called.
	AfterMarketCreated(ctx sdk.Context, market Market) error

	// AfterMarketUpdated is called after UpdateMarket is called.
	AfterMarketUpdated(ctx sdk.Context, market Market) error

	// AfterMarketGenesis is called after x/marketmap init genesis.
	AfterMarketGenesis(ctx sdk.Context, tickers map[string]Market) error
}

var _ MarketMapHooks = &MultiMarketMapHooks{}

// MultiMarketMapHooks defines an array of MarketMapHooks which can be executed in sequence.
type MultiMarketMapHooks []MarketMapHooks

// AfterMarketCreated calls all AfterMarketCreated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketCreated(ctx sdk.Context, market Market) error {
	for i := range mh {
		if err := mh[i].AfterMarketCreated(ctx, market); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketUpdated calls all AfterMarketUpdated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketUpdated(ctx sdk.Context, market Market) error {
	for i := range mh {
		if err := mh[i].AfterMarketUpdated(ctx, market); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketGenesis calls all AfterMarketGenesis hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketGenesis(ctx sdk.Context, markets map[string]Market) error {
	for i := range mh {
		if err := mh[i].AfterMarketGenesis(ctx, markets); err != nil {
			return err
		}
	}

	return nil
}

// MarketMapHooksWrapper is a wrapper for modules to inject MarketMapHooks using depinject.
type MarketMapHooksWrapper struct{ MarketMapHooks }

var _ MarketMapHooks = &NoopMarketMapHooks{}

// NoopMarketMapHooks defines market map hooks that are a no-op.
type NoopMarketMapHooks struct{}

func (n *NoopMarketMapHooks) AfterMarketCreated(_ sdk.Context, _ Market) error {
	return nil
}

func (n *NoopMarketMapHooks) AfterMarketUpdated(_ sdk.Context, _ Market) error {
	return nil
}

func (n *NoopMarketMapHooks) AfterMarketGenesis(_ sdk.Context, _ map[string]Market) error {
	return nil
}
