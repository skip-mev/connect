package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"

	marketmapmodulev1 "github.com/skip-mev/slinky/api/slinky/marketmap/module/v1"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
)

var _ depinject.OnePerModuleType = AppModule{}

// Inputs contains the dependencies required for module construction.
type Inputs struct {
	depinject.In

	// module dependencies
	Config       *marketmapmodulev1.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService
}

// Outputs defines the constructor outputs for the module.
type Outputs struct {
	depinject.Out

	MarketMapKeeper keeper.Keeper
	Module          appmodule.AppModule
}

// ProvideModule is the depinject constructor for the module.
func ProvideModule(in Inputs) Outputs {
	marketmapKeeper := keeper.NewKeeper(in.StoreService, in.Cdc)

	m := NewAppModule(in.Cdc, marketmapKeeper)

	return Outputs{
		MarketMapKeeper: marketmapKeeper,
		Module:          m,
	}
}
