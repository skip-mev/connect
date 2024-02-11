package types

import "fmt"

type AggregationConfigAlreadyExistsError struct {
	ticker string
}

func NewAggregationConfigAlreadyExistsError(ticker string) AggregationConfigAlreadyExistsError {
	return AggregationConfigAlreadyExistsError{ticker: ticker}
}

func (e AggregationConfigAlreadyExistsError) Error() string {
	return fmt.Sprintf("aggregation config already exists for ticker %s", e.ticker)
}

type MarketConfigAlreadyExistsError struct {
	provider string
}

func NewMarketConfigAlreadyExistsError(key string) MarketConfigAlreadyExistsError {
	return MarketConfigAlreadyExistsError{provider: key}
}

func (e MarketConfigAlreadyExistsError) Error() string {
	return fmt.Sprintf("market config already exists for provider %s", e.provider)
}
