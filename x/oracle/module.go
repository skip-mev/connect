package oracle

import (
	"context"
	"encoding/json"

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

	oraclemodulev1 "github.com/skip-mev/connect/v2/api/connect/oracle/module/v2"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/oracle/client/cli"
	"github.com/skip-mev/connect/v2/x/oracle/keeper"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

// ConsensusVersion is the x/oracle module's current version, as modules integrate and updates are made, this value determines what
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

// AppModuleBasic defines the base interface that the x/oracle module exposes to the application.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the name of this module.
func (amb AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the necessary types from the x/oracle module for amino serialization.
func (amb AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the necessary implementations / interfaces in the x/oracle module w/ the interface-registry ir.
func (amb AppModuleBasic) RegisterInterfaces(ir codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(ir)
}

// RegisterGRPCGatewayRoutes registers the necessary REST routes for the GRPC-gateway to the x/oracle module QueryService on mux. This method
// panics on failure.
func (amb AppModuleBasic) RegisterGRPCGatewayRoutes(cliCtx client.Context, mux *runtime.ServeMux) {
	// register the gate-way routes w/ the provided mux
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(cliCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd is a no-op, as no txs are registered for submission (apart from messages that can only be executed by governance).
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the x/oracle module base query cli-command.
func (amb AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule represents an application module for the x/oracle module.
type AppModule struct {
	AppModuleBasic

	k keeper.Keeper
}

// BeginBlock calls the x/oracle keeper's BeginBlocker function.
func (am AppModule) BeginBlock(goCtx context.Context) error {
	return am.k.BeginBlocker(goCtx)
}

// EndBlock  is a no-op for x/oracle.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// NewAppModule returns an application module for the x/oracle module.
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterServices registers the module's services with the app's module configurator.
func (am AppModule) RegisterServices(cfc module.Configurator) {
	// register MsgServer
	types.RegisterMsgServer(cfc.MsgServer(), keeper.NewMsgServer(am.k))
	// register Query Service
	types.RegisterQueryServer(cfc.QueryServer(), keeper.NewQueryServer(am.k))
}

// DefaultGenesis returns default genesis state as raw bytes for the oracle
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	// by default no CurrencyPairs will be added to state initially
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the oracle module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	// unmarshal genesis-state
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}

	// validate
	return gs.Validate()
}

// RegisterInvariants registers the invariants of the oracle module. If an invariant
// deviates from its predicted value, the InvariantRegistry triggers appropriate
// logic (most often the chain will be halted). No invariants exist for the oracle module.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the genesis initialization for the x/oracle module. It determines the
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

func init() {
	appmodule.Register(
		&oraclemodulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type Inputs struct {
	depinject.In

	// keepers
	types.MarketMapKeeper

	// module-dependencies
	Config       *oraclemodulev1.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService
}

type Outputs struct {
	depinject.Out

	OracleKeeper *keeper.Keeper
	Module       appmodule.AppModule
	Hooks        marketmaptypes.MarketMapHooksWrapper
}

func ProvideModule(in Inputs) Outputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	oracleKeeper := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.MarketMapKeeper,
		authority,
	)

	m := NewAppModule(in.Cdc, oracleKeeper)

	return Outputs{
		OracleKeeper: &oracleKeeper,
		Module:       m,
		Hooks:        marketmaptypes.MarketMapHooksWrapper{MarketMapHooks: oracleKeeper.Hooks()},
	}
}
