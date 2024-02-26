package module

import (
	"fmt"
	"sort"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"golang.org/x/exp/maps"

	marketmapmodulev1 "github.com/skip-mev/slinky/api/slinky/marketmap/module/v1"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

var _ depinject.OnePerModuleType = AppModule{}

func init() {
	appmodule.Register(
		&marketmapmodulev1.Module{},
		appmodule.Provide(ProvideModule),
		appmodule.Invoke(InvokeSetMarketMapHooks),
	)
}

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

	MarketMapKeeper *keeper.Keeper
	Module          appmodule.AppModule
}

// ProvideModule is the depinject constructor for the module.
func ProvideModule(in Inputs) Outputs {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	marketmapKeeper := keeper.NewKeeper(in.StoreService, in.Cdc, authority)

	m := NewAppModule(in.Cdc, marketmapKeeper)

	return Outputs{
		MarketMapKeeper: marketmapKeeper,
		Module:          m,
	}
}

// InvokeSetMarketMapHooks uses the module config to set the hooks on the module.
func InvokeSetMarketMapHooks(
	config *marketmapmodulev1.Module,
	keeper *keeper.Keeper,
	hooks map[string]types.MarketMapHooksWrapper,
) error {
	// all arguments to invokers are optional
	if keeper == nil || config == nil {
		return nil
	}

	modNames := maps.Keys(hooks)
	order := config.HooksOrder
	if len(order) == 0 {
		order = modNames
		sort.Strings(order)
	}

	if len(order) != len(modNames) {
		return fmt.Errorf("len(hooks_order: %v) != len(hooks modules: %v)", order, modNames)
	}

	if len(modNames) == 0 {
		return nil
	}

	var multiHooks types.MultiMarketMapHooks
	for _, modName := range order {
		hook, ok := hooks[modName]
		if !ok {
			return fmt.Errorf("can't find marketmap hooks for module %s", modName)
		}

		multiHooks = append(multiHooks, hook)
	}

	keeper.SetHooks(multiHooks)
	return nil
}
