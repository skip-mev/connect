package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

// InitGenesis initializes the set of CurrencyPairs + their genesis prices (if any) for the x/oracle module.
// this function panics on any errors, i.e. if the genesis state is invalid, or any state-modifications fail.
func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	// validate the genesis
	if err := gs.Validate(); err != nil {
		panic(err)
	}

	// initialize all CurrencyPairs + genesis prices
	for _, cpg := range gs.CurrencyPairGenesis {
		state := types.NewCurrencyPairState(cpg.Id, cpg.Nonce, cpg.CurrencyPairPrice)

		if err := k.currencyPairs.Set(ctx, cpg.CurrencyPair.String(), state); err != nil {
			panic(fmt.Errorf("error in genesis: %w", err))
		}
	}

	// set the next ID to state
	if err := k.nextCurrencyPairID.Set(ctx, gs.NextId); err != nil {
		panic(fmt.Errorf("error in genesis: %w", err))
	}

	if err := k.numCPs.Set(ctx, uint64(len(gs.CurrencyPairGenesis))); err != nil {
		panic(fmt.Errorf("error in genesis: %w", err))
	}

	if err := k.numRemoves.Set(ctx, 0); err != nil {
		panic(fmt.Errorf("error in genesis: %w", err))
	}
}

// ExportGenesis retrieve all CurrencyPairs + QuotePrices set for the module, and return them as a genesis state.
// This module panics on any errors encountered in execution.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// get the current next ID
	id, err := k.nextCurrencyPairID.Peek(ctx)
	if err != nil {
		panic(fmt.Errorf("error in genesis: %w", err))
	}

	// instantiate genesis-state w/ empty array
	gs := &types.GenesisState{
		CurrencyPairGenesis: make([]types.CurrencyPairGenesis, 0),
		NextId:              id,
	}

	// next, iterate over NonceKey to retrieve any CurrencyPairs that have not yet been traversed (CurrencyPairs w/ no Price info)
	err = k.IterateCurrencyPairs(ctx, func(cp connecttypes.CurrencyPair, cps types.CurrencyPairState) {
		// append the currency pair + state to the genesis state
		gs.CurrencyPairGenesis = append(gs.CurrencyPairGenesis, types.CurrencyPairGenesis{
			CurrencyPair:      cp,
			Id:                cps.Id,
			Nonce:             cps.Nonce,
			CurrencyPairPrice: cps.Price,
		})
	})
	if err != nil {
		panic(err)
	}

	return gs
}
