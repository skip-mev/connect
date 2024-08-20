package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// ExecuteByIncentiveTypeCB is a callback function that can utilized to update all
// incentives of a given type. This is useful for having incentive/strategy pairs
// that are meant to last several blocks. This function should return the updated
// incentive, or nil if the incentive should be deleted.
type ExecuteByIncentiveTypeCB func(incentive types.Incentive) (types.Incentive, error)

// ExecuteStrategies executes all of the strategies with the stored incentives.
func (k Keeper) ExecuteStrategies(ctx sdk.Context) error {
	for incentive, strategy := range k.incentiveStrategies {
		if err := k.ExecuteIncentiveStrategy(ctx, incentive, strategy); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteIncentiveStrategy executes a given strategy for all incentives of a given type.
// Note that the strategy may mutate the incentive, and return a new incentive to be
// stored. Strategies must return nil if the incentive should be deleted. Otherwise, the
// incentive will be updated.
func (k Keeper) ExecuteIncentiveStrategy(
	ctx sdk.Context,
	incentive types.Incentive,
	strategy types.Strategy,
) error {
	cb := func(incentive types.Incentive) (types.Incentive, error) {
		return strategy(ctx, incentive)
	}

	return k.ExecuteByIncentiveType(ctx, incentive, cb)
}

// ExecuteByIncentiveType updates all incentives of a given type.
func (k Keeper) ExecuteByIncentiveType(
	ctx sdk.Context,
	incentive types.Incentive,
	cb ExecuteByIncentiveTypeCB,
) error {
	// get iterator for store w/ prefix
	store := ctx.KVStore(k.storeKey)
	key := types.GetIncentiveKey(incentive)
	it := storetypes.KVStorePrefixIterator(store, key)

	// close the iterator
	defer it.Close()
	for ; it.Valid(); it.Next() {
		// Unmarshal the incentive.
		if err := incentive.Unmarshal(it.Value()); err != nil {
			return err
		}

		update, err := cb(incentive)
		if err != nil {
			return err
		}

		// If the callback returns nil, then delete the incentive.
		if update == nil {
			store.Delete(it.Key())
		} else {
			updateBz, err := update.Marshal()
			if err != nil {
				return err
			}

			store.Set(it.Key(), updateBz)
		}
	}

	return nil
}
