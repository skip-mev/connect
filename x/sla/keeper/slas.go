package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// PriceFeedSLACB is a callback function that is executed for each SLA in the
// x/sla module's state.
type PriceFeedSLACB func(sla slatypes.PriceFeedSLA) error

// GetSLA returns the SLA with the given ID from the x/sla module's state.
func (k *Keeper) GetSLA(ctx sdk.Context, slaID string) (slatypes.PriceFeedSLA, error) {
	return k.slas.Get(ctx, slaID)
}

// GetSLAs returns the set of SLAs that are currently in the x/sla module's state.
func (k *Keeper) GetSLAs(ctx sdk.Context) ([]slatypes.PriceFeedSLA, error) {
	var slas []slatypes.PriceFeedSLA
	cb := func(sla slatypes.PriceFeedSLA) error {
		slas = append(slas, sla)
		return nil
	}

	if err := k.iterateSLAs(ctx, cb); err != nil {
		return nil, err
	}

	return slas, nil
}

// AddSLAs adds a set of SLAs to the x/sla module's state. Note, this will
// overwrite any existing SLA with the same ID.
func (k *Keeper) AddSLAs(ctx sdk.Context, slas []slatypes.PriceFeedSLA) error {
	for _, sla := range slas {
		if err := k.SetSLA(ctx, sla); err != nil {
			return err
		}
	}

	return nil
}

// SetSLA sets an SLA to the x/sla module's state. Note, this will overwrite any
// existing SLA with the same ID.
func (k *Keeper) SetSLA(ctx sdk.Context, sla slatypes.PriceFeedSLA) error {
	return k.slas.Set(ctx, sla.ID, sla)
}

// RemoveSLAs removes a set of SLAs from the x/sla module's state.
func (k *Keeper) RemoveSLAs(ctx sdk.Context, slaIDs []string) error {
	for _, id := range slaIDs {
		if err := k.RemoveSLA(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

// RemoveSLA removes an SLA from the x/sla module's state. If the SLA does not
// exist, the function will not error.
func (k *Keeper) RemoveSLA(ctx sdk.Context, slaID string) error {
	return k.slas.Remove(ctx, slaID)
}

// iterateSLAs iterates over the set of SLAs that are currently in the x/sla
// module's state. The function inputs a callback that will be executed for each
// SLA in the state.
func (k *Keeper) iterateSLAs(ctx sdk.Context, cb PriceFeedSLACB) error {
	iterator, err := k.slas.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		sla, err := iterator.Value()
		if err != nil {
			return err
		}

		if err := cb(sla); err != nil {
			return err
		}
	}

	return nil
}
