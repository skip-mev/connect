package types

const (
	// Name of module for external use
	ModuleName = "oracle"
	// Top level store key for the oracle module
	StoreKey = ModuleName
)

const (
	keyPrefixCurrencyPair = iota
)

func (cp CurrencyPair) GetStoreKeyForCurrencyPair() []byte {
	return append([]byte{keyPrefixCurrencyPair}, []byte(cp.ToString())...)
}

func GetCurrencyPairFromKey(bz []byte) CurrencyPair {
	// chop off prefix
	bz = bz[1:]
	return CurrencyPairFromString(string(bz))
}