package types

const (
	// Name of module for external use
	ModuleName = "oracle"
	// Top level store key for the oracle module
	StoreKey = ModuleName
)

const (
	keyPrefixCurrencyPairIdx = iota
)

var (
	// Key Prefix under which all CurrencyPairs + QuotePrices will be stored under
	KeyPrefixCurrencyPair = []byte{keyPrefixCurrencyPairIdx}
)

// Get the Prefix for a given QuotePrice for a CurrencyPair
func (cp CurrencyPair) GetStoreKeyForCurrencyPair() []byte {
	return append(KeyPrefixCurrencyPair, []byte(cp.ToString())...)
}

// Get a CurrencyPair from a CurrencyPair store-index. This method errors if the
// CurrencyPair store-index is incorrectly formatted.
func GetCurrencyPairFromKey(bz []byte) (CurrencyPair, error) {
	// chop off prefix
	bz = bz[len(KeyPrefixCurrencyPair):]
	return CurrencyPairFromString(string(bz))
}
