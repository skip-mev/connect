package alerts

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"

	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	alertsmodulev1 "github.com/skip-mev/connect/v2/api/slinky/alerts/module/v1"
	alertclient "github.com/skip-mev/connect/v2/x/alerts/client"
	"github.com/skip-mev/connect/v2/x/alerts/keeper"
	"github.com/skip-mev/connect/v2/x/alerts/types"
	"github.com/skip-mev/connect/v2/x/alerts/types/strategies"
)

// ConsensusVersion is the x/alerts module's current version, as modules integrate and
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

// AppModuleBasic defines the base interface that the x/alerts module exposes to the
// application.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the name of this module.
func (AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the necessary types from the x/alerts module
// for amino serialization.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// register strategies legacy amino codec
	strategies.RegisterLegacyAminoCodec(cdc)

	// register alerts legacy amino codec
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the necessary implementations / interfaces in the
// x/alerts module w/ the interface-registry.
func (AppModuleBasic) RegisterInterfaces(ir codectypes.InterfaceRegistry) {
	// register the msgs / interfaces for the alerts module
	types.RegisterInterfaces(ir)

	// register the strategies
	strategies.RegisterInterfaces(ir)
}

// RegisterGRPCGatewayRoutes registers the necessary REST routes for the GRPC-gateway to
// the x/alerts module QueryService on mux. This method panics on failure.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(cliCtx client.Context, mux *runtime.ServeMux) {
	// Register the gate-way routes w/ the provided mux.
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(cliCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd is a no-op, as no txs are registered for submission.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return alertclient.GetTxCmd()
}

// GetQueryCmd returns the x/alerts module base query cli-command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return alertclient.GetQueryCmd()
}

// DefaultGenesis returns default genesis state as raw bytes for the alerts
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	gs := types.DefaultGenesisState()
	return cdc.MustMarshalJSON(&gs)
}

// ValidateGenesis performs genesis state validation for the x/alerts module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}

	return gs.ValidateBasic()
}

// AppModule represents an application module for the x/alerts module.
type AppModule struct {
	AppModuleBasic

	k keeper.Keeper
}

// BeginBlock is a no-op for x/alerts.
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// NewAppModule returns an application module for the x/alerts module.
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) error {
	return am.k.EndBlocker(ctx)
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

	// Register the message service.
	types.RegisterMsgServer(cfc.MsgServer(), keeper.NewMsgServer(am.k))
}

// RegisterInvariants registers the invariants of the alerts module. If an invariant
// deviates from its predicted value, the InvariantRegistry triggers appropriate
// logic (most often the chain will be halted). No invariants exist for the alerts module.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the genesis initialization for the x/alerts module. It determines the
// genesis state to initialize from via a json-encoded genesis-state. This method returns no validator set updates.
// This method panics on any errors.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(bz, &gs)

	am.k.InitGenesis(ctx, gs)
}

// ExportGenesis returns the alerts module's exported genesis state as raw
// JSON bytes. This method panics on any error.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.k.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

func init() {
	appmodule.Register(
		&alertsmodulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type Inputs struct {
	depinject.In

	// module-dependencies
	Config       *alertsmodulev1.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService

	// Keepers
	IncentiveKeeper types.IncentiveKeeper
	OracleKeeper    types.OracleKeeper
	BankKeeper      types.BankKeeper

	// HandleValidatorIncentive function
	ValidatorIncentiveHandler strategies.ValidatorIncentiveHandler `optional:"true"`
}

type Outputs struct {
	depinject.Out

	AlertsKeeper keeper.Keeper
	Module       appmodule.AppModule
}

func ProvideModule(in Inputs) Outputs {
	var authority sdk.AccAddress

	// if an authority is given, attempt to parse it, and panic if this fails
	if in.Config.Authority != "" {
		var err error
		authority, err = sdk.AccAddressFromBech32(in.Config.Authority)
		if err != nil {
			panic(err)
		}
	} else {
		// otherwise, default to the governance module account
		authority = authtypes.NewModuleAddress(govtypes.ModuleName)
	}

	if in.ValidatorIncentiveHandler == nil {
		in.ValidatorIncentiveHandler = strategies.DefaultHandleValidatorIncentive()
	}

	alertsKeeper := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.OracleKeeper,
		in.BankKeeper,
		in.IncentiveKeeper,
		in.ValidatorIncentiveHandler,
		authority,
	)

	m := NewAppModule(in.Cdc, *alertsKeeper)

	return Outputs{AlertsKeeper: *alertsKeeper, Module: m}
}
