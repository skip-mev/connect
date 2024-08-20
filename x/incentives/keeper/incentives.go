package keeper

import (
	storetypes "cosmossdk.io/store/types"
	db "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// GetIncentivesByType returns all incentives of a given type.
func (k Keeper) GetIncentivesByType(ctx sdk.Context, incentive types.Incentive) ([]types.Incentive, error) {
	key := types.GetIncentiveKey(incentive)

	// Create a callback to unmarshal the incentives.
	var incentives []types.Incentive
	cb := func(it db.Iterator) error {
		if err := incentive.Unmarshal(it.Value()); err != nil {
			return err
		}

		// Copy the incentive, and append it to the list of incentives.
		incentives = append(incentives, incentive.Copy())
		return nil
	}

	// Iterate through all incentives of the given type, unmashalling them,
	// and appending them to the list of incentives.
	if err := k.iteratorFunc(ctx, key, cb); err != nil {
		return nil, err
	}

	return incentives, nil
}

// AddIncentives adds a set of incentives to the module's state.
func (k Keeper) AddIncentives(ctx sdk.Context, incentives []types.Incentive) error {
	for _, incentive := range incentives {
		if err := k.addIncentive(ctx, incentive); err != nil {
			return err
		}
	}

	return nil
}

// addIncentive adds a single incentive to the module's state.
func (k Keeper) addIncentive(ctx sdk.Context, incentive types.Incentive) error {
	// Get the next incentive index.
	index, err := k.getIncentiveCount(ctx, incentive)
	if err != nil {
		return err
	}

	// Marshal the incentive into
	bz, err := incentive.Marshal()
	if err != nil {
		return err
	}

	// set the incentive in the store
	store := ctx.KVStore(k.storeKey)
	key := types.GetIncentiveKeyWithIndex(incentive, index+1)
	store.Set(key, bz)

	// increment the incentive count
	k.setIncentiveCount(ctx, incentive, index+1)

	return nil
}

// RemoveIncentivesByType removes all incentives of a given type from the module's state.
func (k Keeper) RemoveIncentivesByType(ctx sdk.Context, incentive types.Incentive) error {
	key := types.GetIncentiveKey(incentive)

	// Create a callback to delete the incentives.
	cb := func(it db.Iterator) error {
		store := ctx.KVStore(k.storeKey)
		store.Delete(it.Key())
		return nil
	}

	// Iterate through all incentives of the given type, deleting them.
	return k.iteratorFunc(ctx, key, cb)
}

// getIncentiveCount returns the number of incentives of a given type. Note that this
// is the number of incentives that have been added to the module's state, not the
// number of incentives that are currently active.
func (k Keeper) getIncentiveCount(ctx sdk.Context, incentive types.Incentive) (uint64, error) {
	key := types.GetIncentiveCountKey(incentive)

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return 0, nil
	}

	return sdk.BigEndianToUint64(bz), nil
}

// setIncentiveCount updates the number of incentives of a given type. Note that this
// is the number of incentives that have been added to the module's state, not the
// number of incentives that are currently active.
func (k Keeper) setIncentiveCount(ctx sdk.Context, incentive types.Incentive, count uint64) {
	key := types.GetIncentiveCountKey(incentive)

	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(key, bz)
}

// iteratorFunc is a helper function that will create an iterator for a given
// store, and execute a call-back for each key/value pair.
func (k Keeper) iteratorFunc(ctx sdk.Context, prefix []byte, f func(db.Iterator) error) error {
	// get iterator for store w/ prefix
	store := ctx.KVStore(k.storeKey)
	it := storetypes.KVStorePrefixIterator(store, prefix)

	// close the iterator
	defer it.Close()
	for ; it.Valid(); it.Next() {
		// execute call-back, and return error if necessary
		if err := f(it); err != nil {
			return err
		}
	}
	return nil
}
