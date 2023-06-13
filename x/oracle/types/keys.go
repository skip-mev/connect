package types

const (
	// Name of module for external use
	ModuleName = "oracle"
	// Top level store key for the oracle module
	StoreKey = ModuleName
)

const (
	keyPrefixTicker = iota
)

func (t Ticker) GetStoreKeyForTicker() []byte {
	return append([]byte{keyPrefixTicker}, []byte(t.String())...)
}
