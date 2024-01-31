package config

import oracletypes "github.com/skip-mev/slinky/x/oracle/types"

// AggregateMarketConfig represents the configurations for the conversion market. Specifically,
// this allows the oracle to convert prices between different currency pairs to resolve to a
// common currency pair. For example, if the oracle receives a price for BTC/USDT and USDT/USD,
// it can use the conversion market to convert the BTC/USDT price to BTC/USD.
type AggregateMarketConfig struct {
	// CurrencyPairs is the list of currency pairs that the oracle will fetch prices for.
	CurrencyPairs map[string]AggregateCurrencyPairConfig `mapstructure:"currency_pair_config" toml:"currency_pair_config"`
}

// AggregateCurrencyPairConfig represents the configurations for a single currency pair.
type AggregateCurrencyPairConfig struct {
	// CurrencyPair is the currency pair that the oracle will fetch prices for.
	CurrencyPair oracletypes.CurrencyPair `mapstructure:"currency_pair" toml:"currency_pair"`

	// ConvertableMarkets is the list of valid markets (i.e. price feeds) that can be
	// used to convert the price of the currency pair to a common currency pair. For
	// example, if the oracle receives a price for BTC/USDT and USDT/USD, it can use
	// the conversion market to convert the BTC/USDT price to BTC/USD.
	ConvertableMarkets [][]oracletypes.CurrencyPair `mapstructure:"convertable_markets" toml:"convertable_markets"`
}
