package types

import (
	"cosmossdk.io/collections"
)

const (
	// Name of module for external use
	ModuleName = "oracle"
	// Top level store key for the oracle module
	StoreKey = ModuleName
)

var (
	// CurrencyPairKeyPrefix is the key-prefix under which currency-pair state is stored
	CurrencyPairKeyPrefix = collections.NewPrefix(0)

	// CurrencyPairIDKeyPrefix is the key-prefix under which the next currency-pairID is stored
	CurrencyPairIDKeyPrefix = collections.NewPrefix(1)

	// UniqueIndexCurrencyPairKeyPrefix is the key-prifix under which the unique index on
	// currency-pairs is stored
	UniqueIndexCurrencyPairKeyPrefix = collections.NewPrefix(2)

	// IDIndexCurrencyPairKeyPrefix
	IDIndexCurrencyPairKeyPrefix = collections.NewPrefix(3)
)
