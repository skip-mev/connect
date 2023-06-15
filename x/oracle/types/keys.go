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
	keyPrefixCurrencyPairNonceIdx
)

var (
	// KeyPrefixQuotePrice is the key prefix under which all CurrencyPairs + QuotePrices will be stored under
	KeyPrefixQuotePrice = []byte{keyPrefixCurrencyPairIdx}
	// KeyPrefixNonce is the key prefix under which all CurrencyPairs + nonces are stored
	KeyPrefixNonce = []byte{keyPrefixCurrencyPairNonceIdx}
)

// GetStoreKeyForQuotePrice gets the QuotePrice store-prefix for a CurrencyPair
func (cp CurrencyPair) GetStoreKeyForQuotePrice() []byte {
	return append(KeyPrefixQuotePrice, []byte(cp.ToString())...)
}

// GetStoreKeyForNonce gets the store-prefix for nonces from the CurrencyPair
func (cp CurrencyPair) GetStoreKeyForNonce() []byte {
	return append(KeyPrefixNonce, []byte(cp.ToString())...)
}

// GetCurrencyPairFromQuotePriceKey gets a CurrencyPair from a CurrencyPair store-index. This method errors if the
// CurrencyPair store-index is incorrectly formatted.
func GetCurrencyPairFromQuotePriceKey(bz []byte) (CurrencyPair, error) {
	if len(bz) < len(KeyPrefixQuotePrice) {
		return CurrencyPair{}, fmt.Errorf("invalid length of key: %v", len(bz))
	}
	// chop off prefix
	bz = bz[len(KeyPrefixQuotePrice):]
	return CurrencyPairFromString(string(bz))
}
