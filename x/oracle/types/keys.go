package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

const (
	// ModuleName is the name of module for external use.
	ModuleName = "oracle"
	// StoreKey is the top-level store key for the oracle module.
	StoreKey = ModuleName
)

var (
	// CurrencyPairKeyPrefix is the key-prefix under which currency-pair state is stored.
	CurrencyPairKeyPrefix = collections.NewPrefix(0)

	// CurrencyPairIDKeyPrefix is the key-prefix under which the next currency-pairID is stored.
	CurrencyPairIDKeyPrefix = collections.NewPrefix(1)

	// UniqueIndexCurrencyPairKeyPrefix is the key-prefix under which the unique index on
	// currency-pairs is stored.
	UniqueIndexCurrencyPairKeyPrefix = collections.NewPrefix(2)

	// IDIndexCurrencyPairKeyPrefix is the key-prefix under which a currency-pair index.
	// is stored.
	IDIndexCurrencyPairKeyPrefix = collections.NewPrefix(3)

	// NumRemovesKeyPrefix is the key-prefix under which the number of removed CPs is stored.
	NumRemovesKeyPrefix = collections.NewPrefix(4)

	// NumCPsKeyPrefix is the key-prefix under which the number CPs is stored.
	NumCPsKeyPrefix = collections.NewPrefix(5)

	// CounterCodec is the collections.KeyCodec value used for the counter values.
	CounterCodec = codec.KeyToValueCodec[uint64](codec.NewUint64Key[uint64]())
)
