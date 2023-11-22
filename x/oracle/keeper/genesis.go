package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/oracle/types"
)

// InitGenesis initializes the set of CurrencyPairs + their genesis prices (if any) for the x/oracle module.
// this function panics on any errors, i.e if the genesis state is invalid, or any state-modifications fail.
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	// validate the genesis
	if err := gs.Validate(); err != nil {
		panic(err)
	}

	// initialize all CurrencyPairs + genesis prices
	for _, cpg := range gs.CurrencyPairGenesis {
		// Only set the CurrencyPair price to state if there is a non-empty price for the pair
		if cpg.CurrencyPairPrice != nil {
			qp := *cpg.CurrencyPairPrice

			// set to state, panic on errors
			if err := k.SetPriceForCurrencyPair(ctx, cpg.CurrencyPair, qp); err != nil {
				panic(err)
			}
		}

		// set the nonce to state
		k.setNonceForCurrencyPair(ctx, cpg.CurrencyPair, cpg.Nonce)

		// set the ID to state
		k.setIDForCurrencyPair(ctx, cpg.CurrencyPair, cpg.Id)
	}

	// set the next ID to state
	k.setNextID(ctx, gs.NextId)
}

// ExportGenesis, retrieve all CurrencyPairs + QuotePrices set for the module, and return them as a genesis state.
// This module panics on any errors encountered in execution.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// instantiate genesis-state w/ empty array
	gs := &types.GenesisState{
		CurrencyPairGenesis: make([]types.CurrencyPairGenesis, 0),
		NextId:              k.getNextID(ctx),
	}

	// next, iterate over NonceKey to retrieve any CurrencyPairs that have not yet been traversed (CurrencyPairs w/ no Price info)
	if err := k.IterateNonces(ctx, func(cp types.CurrencyPair, nonce uint64) {
		// get the id for the currency-pair
		id, ok := k.GetIDForCurrencyPair(ctx, cp)
		if !ok {
			panic(fmt.Errorf("currency pair %s has no id", cp.ToString()))
		}

		cpg := types.CurrencyPairGenesis{
			CurrencyPair: cp,
			Id:           id,
			Nonce:        nonce,
		}

		// get the price for the currency-pair
		qp, err := k.GetPriceForCurrencyPair(ctx, cp)
		if err == nil {
			// if there is a price, set the price to the genesis
			cpg.CurrencyPairPrice = &qp
		}

		// otherwise, aggregate the CurrencyPair and set the Price as nil + nonce as 0
		gs.CurrencyPairGenesis = append(gs.CurrencyPairGenesis, cpg)
	}); err != nil {
		panic(err)
	}

	return gs
}
