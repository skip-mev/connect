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
	// MarketConfigsPrefix is the key prefix for provider MarketConfigs.
	MarketConfigsPrefix = collections.NewPrefix(0)

	// AggregationConfigsPrefix is the key prefix for PathsConfigs per-Ticker.
	AggregationConfigsPrefix = collections.NewPrefix(1)

	// MarketProviderCodec is the collections.KeyCodec value used for the marketConfigs map.
	MarketProviderCodec = codec.NewStringKeyCodec[MarketProvider]()

	// TickerStringCodec is the collections.KeyCodec value used for the aggregationConfigs map.
	TickerStringCodec = codec.NewStringKeyCodec[TickerString]()
)

// MarketProvider is the unique name used to key the MarketConfigs in the marketmap module.
// It is identical to the MarketConfig.Name property which is stored as the value in the Keeper.marketConfigs map.
type MarketProvider string

// TickerString is the key used to identify unique pairs of Base/Quote with corresponding PathsConfig objects--or in other words AggregationConfigs.
// The TickerString is identical to Slinky's CurrencyPair.String() output in that it is `Base` and `Quote` joined by `/` i.e. `$BASE/$QUOTE`.
type TickerString string
