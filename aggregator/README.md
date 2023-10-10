# Aggregator

## Overview

The aggregator is a aggregation alias that allows developers to plug and play different price aggregation strategies. The aggregator maintains
the latest price for a given asset pair for each provider - whether its a validator, API provider, or other. When a aggregated price is requested,
the price aggregator will utilize its configured strategies to determine the final price to return to the caller.

> **NOTE**: Each strategy must be deterministic if used in a distributed environment. This means that the same inputs must always produce the same outputs.

## Configuration

The aggregator is configured by supplimenting a price aggregator strategy.

```golang
// PriceAggregator is a simple aggregator for provider prices.
// It is thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
type PriceAggregator struct {
	mtx sync.RWMutex

	// aggregateFn is the function used to aggregate prices from each provider.
	aggregateFn AggregateFn **<- This is the strategy**

	// providerPrices is a map of provider -> asset -> QuotePrice
	providerPrices AggregatedProviderPrices

	// prices is the current set of prices aggregated across the providers.
	prices map[types.CurrencyPair]*uint256.Int
}
```

The aggreagtion strategy is defined as follows:

```golang
// AggregateFn is the function used to aggregate prices from each provider. Providers
// should be responsible for aggregating prices using TWAPs, TVWAPs, etc. The oracle
// will then compute the canonical price for a given currency pair by computing the
// median price across all providers.
AggregateFn func(providers AggregatedProviderPrices) map[types.CurrencyPair]*uint256.Int

// AggregateFnFromContext is a function that is used to parametrize an aggregateFn by an sdk.Context. This is used
// to allow the aggregateFn to access the latest state of an application. I.e computing a stake weighted median based
// on the latest validator set.
AggregateFnFromContext func(ctx sdk.Context) AggregateFn
```

`AggregateFn` inputs a set of prices from each provider and outputs a map of currency pairs to prices. If developer's need to
implement a strategy that requires stateful information - such as SDK state - they can utilize the `AggregateFnFromContext` type. 
A currency pair is defined as follows:

```golang
// CurrencyPair is the standard representation of a pair of assets, where one
// (Base) is priced in terms of the other (Quote)
type CurrencyPair struct {
	Base  string `protobuf:"bytes,1,opt,name=Base,proto3" json:"Base,omitempty"`
	Quote string `protobuf:"bytes,2,opt,name=Quote,proto3" json:"Quote,omitempty"`
}
```

Please reference the sample implementation - `ComputeMedian` - for an example of how to implement a strategy.

