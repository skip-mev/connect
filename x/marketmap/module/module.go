package module

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

// ConsensusVersion is the x/marketmap module's current version, as modules integrate and updates are made, this value determines what
// version of the module is being run by the chain.
const ConsensusVersion = 1

var (
	_ appmodule.AppModule   = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.HasServices    = AppModule{}
)

// AppModuleBasic is the base struct for the x/marketmap module. It implements the module.AppModuleBasic interface.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the canonical name of the module.
func (amb AppModuleBasic) Name() string {
	return types.ModuleName
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

// AppModule is the actual app module for x/marketmap.
type AppModule struct {
	AppModuleBasic
	k keeper.Keeper
}

// InitGenesis performs the genesis initialization for the x/marketmap module. It determines the
// genesis state to initialize from via a json-encoded genesis-state. This method returns no validator set updates.
// This method panics on any errors.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) []cometabci.ValidatorUpdate {
	// unmarshal genesis-state (panic on errors)
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(bz, &gs)

	// initialize genesis
	am.k.InitGenesis(ctx, gs)

	// return no validator-set updates
	return []cometabci.ValidatorUpdate{}
}

// ExportGenesis returns the oracle module's exported genesis state as raw
// JSON bytes. This method panics on any error.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.k.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// RegisterServices registers the module's services with the app's module configurator.
func (am AppModule) RegisterServices(cfc module.Configurator) {
	// register MsgServer TODO
	// types.RegisterMsgServer(cfc.MsgServer(), keeper.NewMsgServer(am.k))

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
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}
