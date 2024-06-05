package aggregator

import (
	"fmt"
)

// CommitPricesError is an error that is returned when there is a failure in committing the prices to state.
type CommitPricesError struct {
	Err error
}

func (e CommitPricesError) Error() string {
	return fmt.Sprintf("commit prices error: %s", e.Err.Error())
}

func (e CommitPricesError) Label() string {
	return "CommitPricesError"
}

// PriceAggregationError is an error that is returned when there is a failure in aggregating the prices.
type PriceAggregationError struct {
	Err error
}

func (e PriceAggregationError) Error() string {
	return fmt.Sprintf("price aggregation error: %s", e.Err.Error())
}

func (e PriceAggregationError) Label() string {
	return "PriceAggregationError"
}
