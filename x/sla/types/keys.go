package types

const (
	// ModuleName is the name of the module.
	ModuleName = "sla"
	// StoreKey is the store key string for the sla module.
	StoreKey = ModuleName
)

const (
	// keyPrefixSLA is the key prefix under which all SLAs are stored.
	keyPrefixSLA = iota
	// keyPrefixParams is the key prefix used to index the SLA params.
	keyPrefixParams
	// keyPrefixPriceFeeds is the key prefix used to index the price feed incentives.
	keyPrefixPriceFeeds
	// keyPrefixCurrencyPairs is the key prefix used to index the currency pairs.
	keyPrefixCurrencyPairs
)

var (
	// KeyPrefixSLA is the root key prefix under which all SLAs are stored.
	KeyPrefixSLA = []byte{keyPrefixSLA}
	// KeyPrefixParams is the root key prefix used to index the SLA params.
	KeyPrefixParams = []byte{keyPrefixParams}
	// KeyPrefixPriceFeeds is the root key prefix used to index the price feed incentives.
	KeyPrefixPriceFeeds = []byte{keyPrefixPriceFeeds}
	// KeyPrefixCurrencyPairs is the root key prefix used to index the currency pairs.
	KeyPrefixCurrencyPairs = []byte{keyPrefixCurrencyPairs}
)
