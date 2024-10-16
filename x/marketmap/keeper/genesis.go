package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

// InitGenesis initializes the genesis state. Panics if there is an error.
// Any modules that integrate with x/marketmap must set their InitGenesis to occur before the x/marketmap
// module's InitGenesis.  This is so that logic any consuming modules may want to implement in AfterMarketGenesis
// will be run properly.
func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	// validate the genesis
	if err := gs.ValidateBasic(); err != nil {
		panic(err)
	}

	for _, market := range gs.MarketMap.Markets {
		if err := k.CreateMarket(ctx, market); err != nil {
			panic(err)
		}
	}

	if err := k.SetLastUpdated(ctx, gs.LastUpdated); err != nil {
		panic(err)
	}

	if err := k.SetParams(ctx, gs.Params); err != nil {
		panic(err)
	}

	if k.hooks != nil {
		if err := k.hooks.AfterMarketGenesis(ctx, gs.MarketMap.Markets); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis retrieves the genesis from state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	markets, err := k.GetAllMarkets(ctx)
	if err != nil {
		panic(err)
	}

	lastUpdated, err := k.GetLastUpdated(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		MarketMap: types.MarketMap{
			Markets: markets,
		},
		LastUpdated: lastUpdated,
		Params:      params,
	}
}

// InitializeForGenesis is a no-op.
func (k *Keeper) InitializeForGenesis(_ sdk.Context) {}
