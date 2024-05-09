# Slinky Core Upgrades

The following upgrades are provided by this package:

## `InitializeUpgrade`

This upgrade allows developers to initialize the `x/marketmap` and `x/oracle` state when Slinky is first
being enabled in the application.  All that needs to be passed is the set of initial `Params`, and a list 
of `Market`s that are to be instantiated in `x/marketmap` with corresponding price feeds in `x/oracle.
