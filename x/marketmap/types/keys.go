package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

const (
	// ModuleName defines the canonical name identifying the module.
	ModuleName = "marketmap"
	// StoreKey holds the unique key used to access the module keeper's KVStore.
	StoreKey = ModuleName
)

var (
	// LastUpdatedPrefix is the key prefix for the lastUpdated height.
	LastUpdatedPrefix = collections.NewPrefix(1)

	// MarketsPrefix is the key prefix for Markets.
	MarketsPrefix = collections.NewPrefix(2)

	// ParamsPrefix is the key prefix of the module Params.
	ParamsPrefix = collections.NewPrefix(3)

	// TickersCodec is the collections.KeyCodec value used for the markets map.
	TickersCodec = codec.NewStringKeyCodec[TickerString]()

	// LastUpdatedCodec is the collections.KeyCodec value used for the lastUpdated value.
	LastUpdatedCodec = codec.KeyToValueCodec[uint64](codec.NewUint64Key[uint64]())
)

// TickerString is the key used to identify unique pairs of Base/Quote with corresponding PathsConfig objects--or in other words AggregationConfigs.
// The TickerString is identical to Connect's CurrencyPair.String() output in that it is `Base` and `Quote` joined by `/` i.e. `$BASE/$QUOTE`.
type TickerString string
