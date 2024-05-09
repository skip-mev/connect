package upgrades

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	marketmapkeeper "github.com/skip-mev/slinky/x/marketmap/keeper"
	oraclekeeper "github.com/skip-mev/slinky/x/oracle/keeper"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade interface {
	// CreateUpgradeHandler defines the function that creates an upgrade handler that wraps the provided handler.
	CreateUpgradeHandler(
		mm *module.Manager,
		configurator module.Configurator,
		oracleKeeper *oraclekeeper.Keeper,
		marketMapKeeper *marketmapkeeper.Keeper,
		cdc codec.Codec,
		handler upgradetypes.UpgradeHandler,
	) upgradetypes.UpgradeHandler
}
