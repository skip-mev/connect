package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the start of every block. This will fetch
// all SLAs from state and execute them against the current set
// of price feeds the network is maintaining.
func (k *Keeper) BeginBlocker(ctx sdk.Context) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	if !params.Enabled {
		return nil
	}

	slas, err := k.GetSLAs(ctx)
	if err != nil {
		return err
	}

	for _, sla := range slas {
		if err := k.ExecSLA(ctx, sla); err != nil {
			return err
		}
	}

	return nil
}
