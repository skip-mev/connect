package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	mm2types "github.com/skip-mev/slinky/x/mm2/types"
)

// Hooks is a wrapper struct around Keeper.
type Hooks struct {
	k *Keeper
}

var (
	_ marketmaptypes.MarketMapHooks = Hooks{}
	_ mm2types.MarketMapHooks       = Hooks{}
)

// Hooks returns registered hooks for x/oracle.
func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

// LegacyAfterMarketCreated is the marketmap hook for x/oracle that is run after a market is created in
// the marketmap.  After the market is created, a currency pair and its state are initialized in the
// oracle module.
func (h Hooks) LegacyAfterMarketCreated(ctx sdk.Context, ticker marketmaptypes.Ticker) error {
	return h.k.CreateCurrencyPair(ctx, ticker.CurrencyPair)
}

func (h Hooks) LegacyAfterMarketUpdated(_ sdk.Context, _ marketmaptypes.Ticker) error {
	// TODO finish

	return nil
}

// LegacyAfterMarketGenesis verifies that all tickers set in the x/marketmap genesis are registered in
// the x/oracle module.
func (h Hooks) LegacyAfterMarketGenesis(ctx sdk.Context, tickers map[string]marketmaptypes.Ticker) error {
	for _, ticker := range tickers {
		if !h.k.HasCurrencyPair(ctx, ticker.CurrencyPair) {
			return fmt.Errorf("currency pair %s is registered in x/marketmap but not in x/oracle", ticker.String())
		}
	}

	return nil
}

// AfterMarketCreated is the marketmap hook for x/oracle that is run after a market is created in
// the marketmap.  After the market is created, a currency pair and its state are initialized in the
// oracle module.
func (h Hooks) AfterMarketCreated(ctx sdk.Context, market mm2types.Market) error {
	return h.k.CreateCurrencyPair(ctx, market.Ticker.CurrencyPair)
}

// AfterMarketUpdated is the marketmap hook for x/oracle that is run after a market is updated in
// the marketmap.
func (h Hooks) AfterMarketUpdated(_ sdk.Context, _ mm2types.Market) error {
	// TODO
	return nil
}

// AfterMarketGenesis verifies that all markets set in the x/marketmap genesis are registered in
// the x/oracle module.
func (h Hooks) AfterMarketGenesis(ctx sdk.Context, markets map[string]mm2types.Market) error {
	for _, market := range markets {
		if !h.k.HasCurrencyPair(ctx, market.Ticker.CurrencyPair) {
			return fmt.Errorf("currency pair %s is registered in x/marketmap but not in x/oracle", market.Ticker.String())
		}
	}

	return nil
}
