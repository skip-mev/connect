package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MarketMapHooks is the interface that defines the hooks that can be integrated by other modules.
type MarketMapHooks interface {
	AfterMarketCreated(ctx sdk.Context, ticker Ticker) error

	AfterMarketUpdated(ctx sdk.Context, ticker Ticker) error

	// AfterMarketGenesis is called after x/marketmap init genesis.
	AfterMarketGenesis(ctx sdk.Context, tickers []Tickers) error
}

var _ MarketMapHooks = &MultiMarketMapHooks{}

// MultiMarketMapHooks defines an array of MarketMapHooks which can be executed in sequence.
type MultiMarketMapHooks []MarketMapHooks

// AfterMarketCreated calls all AfterMarketCreated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketCreated(ctx sdk.Context, ticker Ticker) error {
	for i := range mh {
		if err := mh[i].AfterMarketCreated(ctx, ticker); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketUpdated calls all AfterMarketUpdated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketUpdated(ctx sdk.Context, ticker Ticker) error {
	for i := range mh {
		if err := mh[i].AfterMarketUpdated(ctx, ticker); err != nil {
			return err
		}
	}

	return nil
}

// AfterMarketGenesis calls all AfterMarketGenesis hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) AfterMarketGenesis(ctx sdk.Context, tickers []Ticker) error {
	for i := range mh {
		if err := mh[i].AfterMarketGenesis(ctx, tickers); err != nil {
			return err
		}
	}

	return nil
}

// MarketMapHooksWrapper is a wrapper for modules to inject MarketMapHooks using depinject.
type MarketMapHooksWrapper struct{ MarketMapHooks }
