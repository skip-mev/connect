package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MarketMapHooks is the interface that defines the hooks that can be integrated by other modules.
type MarketMapHooks interface {
	LegacyAfterMarketCreated(ctx sdk.Context, ticker Ticker) error

	LegacyAfterMarketUpdated(ctx sdk.Context, ticker Ticker) error

	// LegacyAfterMarketGenesis is called after x/marketmap init genesis.
	LegacyAfterMarketGenesis(ctx sdk.Context, tickers map[string]Ticker) error
}

var _ MarketMapHooks = &MultiMarketMapHooks{}

// MultiMarketMapHooks defines an array of MarketMapHooks which can be executed in sequence.
type MultiMarketMapHooks []MarketMapHooks

// LegacyAfterMarketCreated calls all AfterMarketCreated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) LegacyAfterMarketCreated(ctx sdk.Context, ticker Ticker) error {
	for i := range mh {
		if err := mh[i].LegacyAfterMarketCreated(ctx, ticker); err != nil {
			return err
		}
	}

	return nil
}

// LegacyAfterMarketUpdated calls all AfterMarketUpdated hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) LegacyAfterMarketUpdated(ctx sdk.Context, ticker Ticker) error {
	for i := range mh {
		if err := mh[i].LegacyAfterMarketUpdated(ctx, ticker); err != nil {
			return err
		}
	}

	return nil
}

// LegacyAfterMarketGenesis calls all AfterMarketGenesis hooks registered to the MultiMarketMapHooks.
func (mh MultiMarketMapHooks) LegacyAfterMarketGenesis(ctx sdk.Context, tickers map[string]Ticker) error {
	for i := range mh {
		if err := mh[i].LegacyAfterMarketGenesis(ctx, tickers); err != nil {
			return err
		}
	}

	return nil
}

// MarketMapHooksWrapper is a wrapper for modules to inject MarketMapHooks using depinject.
type MarketMapHooksWrapper struct{ MarketMapHooks }
