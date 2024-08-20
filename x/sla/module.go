package sla

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

	slamodulev1 "github.com/skip-mev/connect/v2/api/slinky/sla/module/v1"
	"github.com/skip-mev/connect/v2/x/sla/client/cli"
	"github.com/skip-mev/connect/v2/x/sla/keeper"
	"github.com/skip-mev/connect/v2/x/sla/types"
)

// ConsensusVersion is the x/sla module's consensus version identifier.
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

// AppModuleBasic defines the base interface that the x/sla module exposes to the application.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the name of x/sla module.
func (amb AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the necessary types from the x/sla module for amino
// serialization.
func (amb AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the necessary implementations / interfaces in the x/sla
// module w/ the interface-registry.
func (amb AppModuleBasic) RegisterInterfaces(ir codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(ir)
}

// RegisterGRPCGatewayRoutes registers the necessary REST routes for the GRPC-gateway to
// the x/sla module QueryService on mux. This method panics on failure.
func (amb AppModuleBasic) RegisterGRPCGatewayRoutes(cliCtx client.Context, mux *runtime.ServeMux) {
	// Register the gate-way routes w/ the provided mux.
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(cliCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd is a no-op, as no txs are registered for submission (apart from messages that
// can only be executed by governance).
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the x/sla module base query cli-command.
func (amb AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule represents an application module for the x/sla module.
type AppModule struct {
	AppModuleBasic

	k keeper.Keeper
}

// NewAppModule returns an application module for the x/sla module.
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}

// BeginBlock returns a beginblocker for the x/sla module.
func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return am.k.BeginBlocker(ctx)
}

// EndBlock is a no-op for x/sla.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterServices registers the module's services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.k))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.k))
}

// DefaultGenesis returns default genesis state as raw bytes for the sla
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.NewDefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the sla module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}

	return gs.ValidateBasic()
}

// InitGenesis performs the genesis initialization for the x/sla module. This method returns
// no validator set updates. This method panics on any errors.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(bz, &gs)

	am.k.InitGenesis(ctx, gs)
}

// ExportGenesis returns the sla module's exported genesis state as raw
// JSON bytes. This method panics on any error.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.k.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

func init() {
	appmodule.Register(
		&slamodulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type Inputs struct {
	depinject.In

	Config       *slamodulev1.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService

	StakingKeeper  types.StakingKeeper
	SlashingKeeper types.SlashingKeeper
}

type Outputs struct {
	depinject.Out

	SLAKeeper keeper.Keeper
	Module    appmodule.AppModule
}

func ProvideModule(in Inputs) Outputs {
	var (
		authority sdk.AccAddress
		err       error
	)
	if in.Config.Authority != "" {
		authority, err = sdk.AccAddressFromBech32(in.Config.Authority)
		if err != nil {
			panic(err)
		}
	} else {
		authority = authtypes.NewModuleAddress(govtypes.ModuleName)
	}

	slaKeeper := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		authority,
		in.StakingKeeper,
		in.SlashingKeeper,
	)

	m := NewAppModule(in.Cdc, *slaKeeper)

	return Outputs{SLAKeeper: *slaKeeper, Module: m}
}
