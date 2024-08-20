package incentives

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	incentivesmodulev1 "github.com/skip-mev/connect/v2/api/slinky/incentives/module/v1"
	"github.com/skip-mev/connect/v2/x/incentives/client/cli"
	"github.com/skip-mev/connect/v2/x/incentives/keeper"
	"github.com/skip-mev/connect/v2/x/incentives/types"
)

// ConsensusVersion is the x/incentives module's current version, as modules integrate and
// updates are made, this value determines what version of the module is being run by the chain.
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

// AppModuleBasic defines the base interface that the x/incentives module exposes to the
// application.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the name of this module.
func (AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the necessary types from the x/incentives module
// for amino serialization.
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// RegisterInterfaces registers the necessary implementations / interfaces in the
// x/incentives module w/ the interface-registry.
func (AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

// RegisterGRPCGatewayRoutes registers the necessary REST routes for the GRPC-gateway to
// the x/incentives module QueryService on mux. This method panics on failure.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(cliCtx client.Context, mux *runtime.ServeMux) {
	// Register the gate-way routes w/ the provided mux.
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(cliCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd is a no-op, as no txs are registered for submission.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the x/incentives module base query cli-command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// DefaultGenesis returns default genesis state as raw bytes for the incentives
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.NewDefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the x/incentives module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}

	return gs.ValidateBasic()
}

// AppModule represents an application module for the x/incentives module.
type AppModule struct {
	AppModuleBasic

	k keeper.Keeper
}

// NewAppModule returns an application module for the x/incentives module.
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}

// BeginBlock returns the beginblocker for the x/incentives module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.k.ExecuteStrategies(sdk.UnwrapSDKContext(ctx))
}

// EndBlock is a no-op for x/incentives.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterServices registers the module's services with the app's module configurator.
func (am AppModule) RegisterServices(cfc module.Configurator) {
	// Register the query service.
	types.RegisterQueryServer(cfc.QueryServer(), keeper.NewQueryServer(am.k))
}

// RegisterInvariants registers the invariants of the incentives module. If an invariant
// deviates from its predicted value, the InvariantRegistry triggers appropriate
// logic (most often the chain will be halted). No invariants exist for the incentives module.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the genesis initialization for the x/incentives module. It determines the
// genesis state to initialize from via a json-encoded genesis-state. This method returns no validator set updates.
// This method panics on any errors.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(bz, &gs)

	am.k.InitGenesis(ctx, gs)
}

// ExportGenesis returns the incentives module's exported genesis state as raw
// JSON bytes. This method panics on any error.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.k.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

func init() {
	appmodule.Register(
		&incentivesmodulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type Inputs struct {
	depinject.In

	Config *incentivesmodulev1.Module
	Cdc    codec.Codec
	Key    *storetypes.KVStoreKey

	// IncentiveStrategies
	IncentiveStrategies map[types.Incentive]types.Strategy `optional:"true"`
}

type Outputs struct {
	depinject.Out

	IncentivesKeeper keeper.Keeper
	Module           appmodule.AppModule
}

func ProvideModule(in Inputs) Outputs {
	incentivesKeeper := keeper.NewKeeper(
		in.Key,
		in.IncentiveStrategies,
	)

	m := NewAppModule(in.Cdc, incentivesKeeper)

	return Outputs{IncentivesKeeper: incentivesKeeper, Module: m}
}
