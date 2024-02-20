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
	// TickersPrefix is the key prefix for Tickers.
	TickersPrefix = collections.NewPrefix(0)

	// PathsPrefix is the key prefix for Paths.
	PathsPrefix = collections.NewPrefix(1)

	// ProvidersPrefix is the key prefix for Providers.
	ProvidersPrefix = collections.NewPrefix(2)

	// LastUpdatedPrefix is the key prefix for the lastUpdated height.
	LastUpdatedPrefix = collections.NewPrefix(3)

	// TickersCodec is the collections.KeyCodec value used for the markets map.
	TickersCodec = codec.NewStringKeyCodec[TickerString]()

	// LastUpdatedCodec is the collections.KeyCodec value used for the lastUpdated value.
	LastUpdatedCodec = codec.KeyToValueCodec[int64](codec.NewInt64Key[int64]())
)

// TickerString is the key used to identify unique pairs of Base/Quote with corresponding PathsConfig objects--or in other words AggregationConfigs.
// The TickerString is identical to Slinky's CurrencyPair.String() output in that it is `Base` and `Quote` joined by `/` i.e. `$BASE/$QUOTE`.
type TickerString string
