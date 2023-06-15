package types

import (
	"fmt"
)

const (
	// Name of module for external use
	ModuleName = "oracle"
	// Top level store key for the oracle module
	StoreKey = ModuleName
)

const (
	keyPrefixCurrencyPairIdx = iota
)

// KeyPrefixCurrencyPair is the key prefix under which all CurrencyPairs + QuotePrices will be stored under
var KeyPrefixCurrencyPair = []byte{keyPrefixCurrencyPairIdx}

// GetStoreKeyForCurrencyPair gets the QuotePrice store-key for a CurrencyPair
func (cp CurrencyPair) GetStoreKeyForCurrencyPair() []byte {
	return append(KeyPrefixCurrencyPair, []byte(cp.ToString())...)
}

// GetCurrencyPairFromKey gets a CurrencyPair from a CurrencyPair store-index. This method errors if the
// CurrencyPair store-index is incorrectly formatted.
func GetCurrencyPairFromKey(bz []byte) (CurrencyPair, error) {
	if len(bz) < len(KeyPrefixCurrencyPair) {
		return CurrencyPair{}, fmt.Errorf("invalid length of key: %v", len(bz))
	}
	// chop off prefix
	bz = bz[len(KeyPrefixCurrencyPair):]
	return CurrencyPairFromString(string(bz))
}
