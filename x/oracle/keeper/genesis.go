package keeper

import (
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
	}
}

// ExportGenesis, retrieve all CurrencyPairs + QuotePrices set for the module, and return them as a genesis state.
// This module panics on any errors encountered in execution.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// instantiate genesis-state w/ empty array
	gs := &types.GenesisState{
		CurrencyPairGenesis: make([]types.CurrencyPairGenesis, 0),
	}

	cpCache := make(map[string]struct{})

	// populate genesis-state w/ CurrencyPairs that have valid QuotePrices first, and cache the CurrencyPairs that have
	// already been traversed
	if err := k.IterateQuotePrices(ctx, func(cp types.CurrencyPair, qp types.QuotePrice) error {
		// get the nonce for the currency pair
		nonce, err := k.GetNonceForCurrencyPair(ctx, cp)
		if err != nil {
			return err
		}

		// aggregate
		gs.CurrencyPairGenesis = append(gs.CurrencyPairGenesis, types.CurrencyPairGenesis{
			CurrencyPair:      cp,
			CurrencyPairPrice: &qp,
			Nonce:             nonce,
		})

		// cache cp as already traversed
		cpCache[cp.ToString()] = struct{}{}

		return nil
	}); err != nil {
		panic(err)
	}

	// next, iterate over NonceKey to retrieve any CurrencyPairs that have not yet been traversed (CurrencyPairs w/ no Price info)
	if err := k.IterateNonces(ctx, func(cp types.CurrencyPair) {
		// check to see if this CurrencyPair has already been traversed, if so, skip
		if _, ok := cpCache[cp.ToString()]; ok {
			return
		}

		// otherwise, aggregate the CurrencyPair and set the Price as nil + nonce as 0
		gs.CurrencyPairGenesis = append(gs.CurrencyPairGenesis, types.CurrencyPairGenesis{
			CurrencyPair: cp,
		})
	}); err != nil {
		panic(err)
	}

	return gs
}
