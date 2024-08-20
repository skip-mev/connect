package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// InitGenesis initializes the store state from a genesis state.
func (k *Keeper) InitGenesis(ctx sdk.Context, gs slatypes.GenesisState) {
	// Validate the genesis state.
	if err := gs.ValidateBasic(); err != nil {
		panic(err)
	}

	// Add each price feed sla to the store.
	for _, sla := range gs.SLAs {
		if err := k.SetSLA(ctx, sla); err != nil {
			panic(err)
		}
	}

	seenCPs := make(map[slinkytypes.CurrencyPair]struct{})
	for _, feed := range gs.PriceFeeds {
		if err := k.SetPriceFeed(ctx, feed); err != nil {
			panic(err)
		}

		seenCPs[feed.CurrencyPair] = struct{}{}
	}

	// Set the currency pairs.
	if err := k.SetCurrencyPairs(ctx, seenCPs); err != nil {
		panic(err)
	}

	// Set the params.
	if err := k.SetParams(ctx, gs.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the current store state as a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *slatypes.GenesisState {
	// Get the set of SLAs.
	slas, err := k.GetSLAs(ctx)
	if err != nil {
		panic(err)
	}

	// Get all price feeds.
	aggFeeds := make([]slatypes.PriceFeed, 0)
	for _, sla := range slas {
		feeds, err := k.GetAllPriceFeeds(ctx, sla.ID)
		if err != nil {
			panic(err)
		}

		aggFeeds = append(aggFeeds, feeds...)
	}

	// Get the params.
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return slatypes.NewGenesisState(slas, aggFeeds, params)
}
