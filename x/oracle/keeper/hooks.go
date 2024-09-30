package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// Hooks is a wrapper struct around Keeper.
type Hooks struct {
	k *Keeper
}

var _ marketmaptypes.MarketMapHooks = Hooks{}

// Hooks returns registered hooks for x/oracle.
func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterMarketCreated is the marketmap hook for x/oracle that is run after a market is created in
// the marketmap.  After the market is created, a currency pair and its state are initialized in the
// oracle module.
func (h Hooks) AfterMarketCreated(ctx sdk.Context, market marketmaptypes.Market) error {
	ctx.Logger().Info(fmt.Sprintf("creating x/oracle state for market %s", market.Ticker.String()))
	return h.k.CreateCurrencyPair(ctx, market.Ticker.CurrencyPair)
}

// AfterMarketUpdated is the marketmap hook for x/oracle that is run after a market is updated in
// the marketmap.
func (h Hooks) AfterMarketUpdated(ctx sdk.Context, market marketmaptypes.Market) error {
	ctx.Logger().Info(fmt.Sprintf("market %s updated", market.Ticker.String()))
	return nil
}

// AfterMarketGenesis verifies that all markets set in the x/marketmap genesis are registered in
// the x/oracle module.
func (h Hooks) AfterMarketGenesis(ctx sdk.Context, markets map[string]marketmaptypes.Market) error {
	for _, market := range markets {
		if !h.k.HasCurrencyPair(ctx, market.Ticker.CurrencyPair) {
			return fmt.Errorf("currency pair %s is registered in x/marketmap but not in x/oracle", market.Ticker.String())
		}
	}

	return nil
}

// AfterMarketRemoved is the marketmap hook for x/oracle that is run after a market is removed in
// the marketmap.
func (h Hooks) AfterMarketRemoved(ctx sdk.Context, key string) error {
	ctx.Logger().Info(fmt.Sprintf("market %s removed. retaining x/oracle state if it exists", key))

	return nil
}
