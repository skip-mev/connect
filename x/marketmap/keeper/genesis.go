package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

// InitGenesis initializes the genesis state. Panics if there is an error.
func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	// validate the genesis
	if err := gs.ValidateBasic(); err != nil {
		panic(err)
	}

	for _, ticker := range gs.MarketMap.Tickers {
		paths, ok := gs.MarketMap.Paths[ticker.String()]
		if !ok {
			panic(fmt.Errorf("paths for ticker %s not found", ticker.String()))
		}

		providers, ok := gs.MarketMap.Providers[ticker.String()]
		if !ok {
			panic(fmt.Errorf("providers for ticker %s not found", ticker.String()))
		}

		if err := k.CreateMarket(ctx, ticker, paths, providers); err != nil {
			panic(err)
		}
	}

	if err := k.SetLastUpdated(ctx, gs.LastUpdated); err != nil {
		panic(err)
	}

	if err := k.SetParams(ctx, gs.Params); err != nil {
		panic(err)
	}

	if err := k.hooks.AfterMarketGenesis(ctx, gs.MarketMap.Tickers); err != nil {
		panic(err)
	}
}

// ExportGenesis retrieves the genesis from state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	tickers, err := k.GetAllTickersMap(ctx)
	if err != nil {
		panic(err)
	}

	paths, err := k.GetAllPathsMap(ctx)
	if err != nil {
		panic(err)
	}

	providers, err := k.GetAllProvidersMap(ctx)
	if err != nil {
		panic(err)
	}

	lastUpdated, err := k.GetLastUpdated(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		MarketMap: types.MarketMap{
			Tickers:   tickers,
			Paths:     paths,
			Providers: providers,
		},
		LastUpdated: lastUpdated,
	}
}
