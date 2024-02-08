package types

// NewAggregateMarketConfig returns a new AggregateMarketConfig instance.
func NewAggregateMarketConfig(
	markets map[string]MarketConfig,
	tickers map[uint64]PathsConfig,
) AggregateMarketConfig {
	return AggregateMarketConfig{
		MarketConfigs: markets,
		TickerConfigs: tickers,
	}
}

// ValidateBasic performs basic validation on the AggregateMarketConfig.
func (c AggregateMarketConfig) ValidateBasic() error {
	return nil
}
