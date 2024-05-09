# Slinky Upgrades

The `upgrades` package contains a suite of upgrade handlers that can be integrated into your application
for various use cases with Slinky.

This package is broken down into the following sub-packages:

- [core](./core/README.md)
  - This package contains common upgrades needed for users of the core marketmap.

We export upgrade handlers that can be easily wrap your own upgrade logic.  This way,
all logic can be encapsulated in a single `upgradehandler`.

All Slinky upgrade handlers follow the given interface:

```go
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
```

Developers can feed their own `upgradehandlers` as the `handler` argument.  This will be called after the given
Slinky upgrade is executed.