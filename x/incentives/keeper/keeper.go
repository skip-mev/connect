package keeper

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

type (
	// Keeper is the base keeper for the x/incentives module.
	Keeper struct {
		storeKey storetypes.StoreKey

		// incentiveStrategies is a map of incentive types to their corresponding strategy
		// functions.
		incentiveStrategies map[types.Incentive]types.Strategy
	}
)

// NewKeeper constructs a new keeper from a store-key and a given set of
// (incentive, strategies) pairings. Note, if the strategies map is empty,
// then the keeper will not be able to process any incentives. This must be
// set by the application developer. Each incentive type must have a
// corresponding strategy function.
func NewKeeper(
	sk storetypes.StoreKey,
	incentiveStrategies map[types.Incentive]types.Strategy,
) Keeper {
	return Keeper{
		storeKey:            sk,
		incentiveStrategies: incentiveStrategies,
	}
}
