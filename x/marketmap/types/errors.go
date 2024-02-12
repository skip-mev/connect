package types

import "fmt"

type AggregationConfigAlreadyExistsError struct {
	ticker TickerString
}

func NewAggregationConfigAlreadyExistsError(ticker TickerString) AggregationConfigAlreadyExistsError {
	return AggregationConfigAlreadyExistsError{ticker: ticker}
}

func (e AggregationConfigAlreadyExistsError) Error() string {
	return fmt.Sprintf("aggregation config already exists for ticker %s", e.ticker)
}

type MarketConfigAlreadyExistsError struct {
	provider MarketProvider
}

func NewMarketConfigAlreadyExistsError(key MarketProvider) MarketConfigAlreadyExistsError {
	return MarketConfigAlreadyExistsError{provider: key}
}

func (e MarketConfigAlreadyExistsError) Error() string {
	return fmt.Sprintf("market config already exists for provider %s", e.provider)
}
