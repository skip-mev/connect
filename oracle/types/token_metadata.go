package types

type TokenMetadata struct {
	// Symbol is the provider-specific token identifier. This can be a name, ticker, contract address, etc.
	Symbol string `mapstructure:"symbol" toml:"symbol"`

	// Decimals is the number of decimal places the token has on chain.
	Decimals uint64 `mapstructure:"decimals" toml:"decimals"`

	// IsTWAP indicates whether this token's price is a time weighted average.
	IsTWAP bool `mapstructure:"is_twap" toml:"is_twap"`
}
