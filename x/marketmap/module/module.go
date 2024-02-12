package module

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
func (amb AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {
	// todo
}

// AppModule is the actual app module for x/marketmap.
type AppModule struct {
	AppModuleBasic
	k keeper.Keeper
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (a AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface. It is a no-op.
func (a AppModule) IsAppModule() {}

// NewAppModule constructs a new application module for the x/marketmap module.
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		k: k,
	}
}
