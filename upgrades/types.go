package upgrades

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"

	marketmapkeeper "github.com/skip-mev/slinky/x/marketmap/keeper"
	oraclekeeper "github.com/skip-mev/slinky/x/oracle/keeper"
)

// Upgrade defines an interface for a Slinky Upgrade.
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

// EmptyUpgrade is a useful alias for an empty upgrade handler you can append as a no-op to wrap.
func EmptyUpgrade(_ context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return fromVM, nil
}
