package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

// InitGenesis initializes the module state from a GenesisState object. Specifically, this method sets the
// params to state and adds all alerts from the genesis state to state.
func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	// validate the genesis state
	if err := gs.ValidateBasic(); err != nil {
		panic(err)
	}

	// set the params
	if err := k.SetParams(ctx, gs.Params); err != nil {
		panic(err)
	}

	// add all alerts
	for _, alert := range gs.Alerts {
		if err := k.SetAlert(ctx, alert); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState object containing the current module state. Specifically, this method
// returns the current params and all alerts in state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// get the params
	params := k.GetParams(ctx)

	// get all alerts
	alerts, err := k.GetAllAlerts(ctx)
	if err != nil {
		panic(err)
	}

	gs := types.NewGenesisState(params, alerts)
	if err := gs.ValidateBasic(); err != nil {
		panic(err)
	}

	return &gs
}
