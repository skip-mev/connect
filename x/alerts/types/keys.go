package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName is the name of the module.
	ModuleName = "alerts"

	// StoreKey is the default store key for alerts.
	StoreKey = ModuleName
)

var (
	// AlertStoreKeyPrefix is the prefix for the alert store key.
	AlertStoreKeyPrefix = collections.NewPrefix(0)
	// ParamsStoreKeyPrefix is the prefix for the params store key.
	ParamsStoreKeyPrefix = collections.NewPrefix(1)
)
