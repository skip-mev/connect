# Oracle Metrics

## Overview

The Oracle Metrics package provides a set of metrics that will be implemented by default for the oracle. These metrics are intended to be used by developers to monitor and track prices and the health of the oracle.

## Oracle Metrics

The following metrics are provided by the API Metrics package for implementations that use the API Query Handler:

```golang
// Metrics is an interface that defines the API for oracle metrics.
type Metrics interface {
	// AddTick increments the number of ticks, this can represent a liveness counter. This
	// is incremented once every interval (which is defined by the oracle config).
	AddTick()

	// UpdatePrice price updates the price for the given pairID for the provider.
	UpdatePrice(name, handlerType, pairID string, price float64)

	// UpdateAggregatePrice rice updates the aggregated price for the given pairID.
	UpdateAggregatePrice(pairID string, price float64)
}
```

### AddTick

The `AddTick` metric is used to track the number of ticks with a fully successful provider update. Everytime a new set of aggregated prices is calculated, this metric is incremented. This can be used to track the liveness of the oracle.

### UpdatePrice

The `UpdatePrice` metric is used to track the price updates for a given provider.

### UpdateAggregatePrice

The `UpdateAggregatePrice` metric is used to track the aggregated price updates for a given pair.

## Usage

Below we overview some of the more useful prometheus queries that can be used to get insight into the oracle overall.

### Graph of the price for a given pair

> ```promql
> oracle_aggregate_price{pair="bitcoin/usd"} # Replace with the pair you want to graph
> ```

This will graph the price for a given pair over time.

### Graph of the price for a given provider

> ```promql
> oracle_price{provider="binance", pair="bitcoin/usd"} # Replace with the provider and pair you want to graph
> ```

This will graph the price for a given provider and pair over time.

### Graph of the aggregated price for a given pair

> ```promql
> oracle_aggregate_price{pair="bitcoin/usd"} # Replace with the pair you want to graph
> ```

This will graph the aggregated price for a given pair over time.

### Number of oracle ticks

> ```promql
> oracle_ticks
> ```

This will return the number of ticks that have occurred. This can be used to track the liveness of the oracle.

### Average number of go-routines over a time window

> ```promql
> avg (go_goroutines)
> ```

This will return the average number of active go-routines over a given time window. This can be used to ensure that there are no runaway go-routines.
