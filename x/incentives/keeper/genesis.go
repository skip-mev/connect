package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// InitGenesis initializes the store state from a genesis state. Note, that
// all of the incentive types (e.g. badprice, goodprice) must be registered
// with the keeper in order for this to execute successfully.
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	// Validate the genesis state.
	if err := gs.ValidateBasic(); err != nil {
		panic(err)
	}

	// Create a reverse map of the incentives.
	reverseMap := make(map[string]types.Incentive)
	for incentive := range k.incentiveStrategies {
		reverseMap[incentive.Type()] = incentive
	}

	// Add each incentive to the store.
	for _, entry := range gs.Registry {
		name, incentives := entry.IncentiveType, entry.Entries

		// Get the incentive type.
		incentiveType, ok := reverseMap[name]
		if !ok {
			panic("unknown incentive type: " + name)
		}

		// Unmarshal each incentive with the correspond type.
		unmarshalledIncentives := make([]types.Incentive, len(incentives))
		for i, bz := range incentives {
			// Attempt to unmarshal the incentive.
			if err := incentiveType.Unmarshal(bz); err != nil {
				panic(err)
			}

			unmarshalledIncentives[i] = incentiveType.Copy()
		}

		// Add the incentives to the store.
		if err := k.AddIncentives(ctx, unmarshalledIncentives); err != nil {
			panic(err)
		}

		// Remove the incentive type from the reverse map since we've already
		// processed it.
		delete(reverseMap, name)
	}
}

// ExportGenesis returns the current store state as a genesis state. Note, that
// if any of the incentive types have no entries in the store, they will not
// be included in the genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	incentiveGenesis := make([]types.IncentivesByType, 0)

	// Get all of the incentive types and sort them by name.
	sortedIncentives := types.SortIncentivesStrategiesMap(k.incentiveStrategies)

	// Iterate over each incentive type.
	for _, incentiveType := range sortedIncentives {
		// Get the incentives for the current type.
		incentives, err := k.GetIncentivesByType(ctx, incentiveType)
		if err != nil {
			panic(err)
		}

		if len(incentives) == 0 {
			continue
		}

		// Marshal each incentive.
		marshalledIncentives, err := types.IncentivesToBytes(incentives...)
		if err != nil {
			panic(err)
		}

		// Add the incentives to the genesis state.
		incentiveGenesis = append(incentiveGenesis, types.NewIncentives(incentiveType.Type(), marshalledIncentives))
	}

	return types.NewGenesisState(incentiveGenesis)
}
