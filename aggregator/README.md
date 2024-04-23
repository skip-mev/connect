# Data Aggregator

## Overview

The aggregator is a aggregation alias that allows developers to plug and play different data aggregation strategies based on the type of data they are working with. The aggregator maintains the latest data i.e. price for a given asset pair for each provider - whether its a validator, API provider, or other. When a aggregated information is requested, the price aggregator will utilize its configured strategies to determine the final result to return to the caller. In the case of price aggregation, the aggregator may return something like the median price across all providers.

> **NOTE**: Each strategy must be deterministic if used in a distributed environment. This means that the same inputs must always produce the same outputs.

## Configuration

The aggregator is configured by supplementing a data aggregator strategy.

```golang
// DataAggregator is a simple aggregator for provider data. It is thread-safe since
// it is assumed to be called concurrently in data fetching goroutines.
type DataAggregator[K comparable, V any] struct {
	mtx sync.RWMutex

	// aggregateFn is the function used to aggregate data from each provider.
	aggregateFn AggregateFn[K, V]

	// providerData is a map of provider -> value (i.e. prices).
	providerData AggregatedProviderData[K, V]

	// aggregatedData is the current set of aggregated data across the providers.
	aggregatedData V
}
```

The aggregation strategy is defined as follows:

```golang
// AggregateFn is the function used to aggregate data from each provider. Given a
// map of provider -> values, the aggregate function should return a final
// value.
AggregateFn[K comparable, V any] func(providers AggregatedProviderData[K, V]) V

// AggregateFnFromContext is a function that is used to parametrize an aggregateFn
// by an sdk.Context. This is used to allow the aggregateFn to access the latest state
// of an application i.e computing a stake weighted median based on the latest validator set.
AggregateFnFromContext[K comparable, V any] func(ctx sdk.Context) AggregateFn[K, V]
```

`AggregateFn` inputs data from each provider and outputs a final value. If developer's need to implement a strategy that requires stateful information - such as SDK state - they can utilize the `AggregateFnFromContext` type.

Please reference the sample implementation - [`ComputeMedian`](math.go) - for an example of how to implement a strategy.

