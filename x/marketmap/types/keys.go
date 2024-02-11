package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the canonical name identifying the module
	ModuleName = "marketmap"
	// StoreKey holds the unique key used to access the module keeper's KVStore
	StoreKey = ModuleName
)

var (
	// MarketConfigsPrefix is the key prefix for provider MarketConfigs
	MarketConfigsPrefix = collections.NewPrefix(0)

	// AggregationConfigsPrefix is the key prefix for PathsConfigs per-Ticker
	AggregationConfigsPrefix = collections.NewPrefix(1)
)
