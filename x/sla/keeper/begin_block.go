package keeper

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the start of every block. This will fetch
// all SLAs from state and execute them against the current set
// of price feeds the network is maintaining.
func (k *Keeper) BeginBlocker(ctx sdk.Context) ([]cometabci.ValidatorUpdate, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	if !params.Enabled {
		return nil, nil
	}

	slas, err := k.GetSLAs(ctx)
	if err != nil {
		return nil, err
	}

	for _, sla := range slas {
		if err := k.ExecSLA(ctx, sla); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
