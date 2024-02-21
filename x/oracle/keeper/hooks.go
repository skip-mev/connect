package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
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

func (h Hooks) AfterMarketCreated(_ sdk.Context, ticker marketmaptypes.Ticker) error {
	// TODO finish

	return nil
}

func (h Hooks) AfterMarketUpdated(_ sdk.Context, ticker marketmaptypes.Ticker) error {
	// TODO finish

	return nil
}
