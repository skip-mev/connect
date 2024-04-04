package mm2

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	marketmapmodulev1 "github.com/skip-mev/slinky/api/slinky/mm2/module/v1"
	"github.com/skip-mev/slinky/x/mm2/client/cli"
	"github.com/skip-mev/slinky/x/mm2/keeper"
	"github.com/skip-mev/slinky/x/mm2/types"
)

// ConsensusVersion is the x/marketmap module's current version, as modules integrate and updates are made, this value determines what
// version of the module is being run by the chain.
const ConsensusVersion = 1

var (
	_ module.HasName        = AppModule{}
	_ module.HasGenesis     = AppModule{}
	_ module.AppModuleBasic = AppModule{}
	_ module.HasServices    = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
	_ appmodule.HasEndBlocker   = AppModule{}
)

// AppModuleBasic is the base struct for the x/marketmap module. It implements the module.AppModuleBasic interface.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the canonical name of the module.
func (amb AppModuleBasic) Name() string {
	return types.ModuleName
}

// GetTxCmd is a no-op, as no txs are registered for submission (apart from messages that can only be executed by governance).
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the x/marketmap module base query cli-command.
func (amb AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterLegacyAminoCodec registers the necessary types from the x/marketmap module for amino serialization.
func (amb AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the necessary implementations / interfaces in the x/marketmap module w/ the interface-registry ir.
func (amb AppModuleBasic) RegisterInterfaces(ir codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(ir)
}

// RegisterGRPCGatewayRoutes registers the necessary REST routes for the GRPC-gateway to the x/marketmap module QueryService on mux.
// This method panics on failure.
func (amb AppModuleBasic) RegisterGRPCGatewayRoutes(cliCtx client.Context, mux *runtime.ServeMux) {
	// register the gate-way routes w/ the provided mux
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(cliCtx)); err != nil {
		panic(err)
	}
}

// DefaultGenesis returns default genesis state as raw bytes for the marketmap
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	gs := types.DefaultGenesisState()
	return cdc.MustMarshalJSON(gs)
}

// ValidateGenesis performs genesis state validation for the marketmap module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}

	return gs.ValidateBasic()
}

// RegisterInvariants registers the invariants of the marketmap module. If an invariant
// deviates from its predicted value, the InvariantRegistry triggers appropriate
// logic (most often the chain will be halted). No invariants exist for the marketmap module.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// AppModule is the actual app module for x/marketmap.
type AppModule struct {
	AppModuleBasic
	k *keeper.Keeper
}

// BeginBlock is a no-op for x/marketmap.
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock is a no-op for x/marketmap.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// InitGenesis performs the genesis initialization for the x/marketmap module. It determines the
// genesis state to initialize from via a json-encoded genesis-state. This method returns no validator set updates.
// This method panics on any errors.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	// unmarshal genesis-state (panic on errors)
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(bz, &gs)

	// initialize genesis
	am.k.InitGenesis(ctx, gs)
}

// ExportGenesis returns the oracle module's exported genesis state as raw
// JSON bytes. This method panics on any error.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.k.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// RegisterServices registers the module's services with the app's module configurator.
func (am AppModule) RegisterServices(cfc module.Configurator) {
	// register MsgServer
	types.RegisterMsgServer(cfc.MsgServer(), keeper.NewMsgServer(am.k))

	// register Query Service
	types.RegisterQueryServer(cfc.QueryServer(), keeper.NewQueryServer(am.k))
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface. It is a no-op.
func (am AppModule) IsAppModule() {}

// NewAppModule constructs a new application module for the x/marketmap module.
func NewAppModule(cdc codec.Codec, k *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}

/*

	Dep-Inject

*/

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
