# API Metrics

## Overview

The API Metrics package provides a set of metrics that will be implemented by default for all providers that inherit from the Base Provider and implement an API-based provider. These metrics are intended to be used by the provider to track the usage of the provider's APIs and the resources it manages.

## API Metrics

The following metrics are provided by the API Metrics package for implementations that use the API Query Handler:

```golang
// APIMetrics is an interface that defines the API for metrics collection for providers
// that implement the APIQueryHandler.
type APIMetrics interface {
	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	// This increments the number of responses by provider, id (i.e. currency pair), and status.
	AddProviderResponse(providerName, id string, status Status)

	// ObserveProviderResponseTime records the time it took for a provider to respond for
	// within a single interval. Note that if the provider is not atomic, this will be the
	// time it took for all of the requests to complete.
	ObserveProviderResponseLatency(providerName string, duration time.Duration)
}
```

### AddProviderResponse

The `AddProviderResponse` metric is used to track the number of ticks with a fully successful provider update. Specifically, this tracks how often providers return good responses within the configured interval.

### ObserveProviderResponseTime

The `ObserveProviderResponseTime` metric is used to track the time it took for a provider to respond. Specifically, provider's must return a response within the configured interval. If the response time is very close to the configured interval, this could indicate that the provider is taking too long to respond, may be timing out, and consuming more resources than necessary.

## Usage

Below we overview some of the more useful prometheus queries that can be used to get insight into the health of a provider.

### Total number of responses

> ```promql
> sum(increase(oracle_api_response_status_per_provider[1h]))
> ```

This will return the total number of provider responses over the last hour (failures and successes across all providers).

### Total number of responses by status

> ```promql
> sum by (status) (increase(oracle_api_response_status_per_provider[1h]))
> ```

This will return the total number of provider responses by status over the last hour.

### Total number of responses by ID (i.e. price feed) and status

> ```promql
> sum by (id, status) (increase(oracle_api_response_status_per_provider[1h]))
> ```

This will return the total number of provider responses by ID and status over the last hour.

### Average number of responses by provider

> ```promql
> sum by (provider) (increase(oracle_api_response_time_per_provider_sum[1h])) / (60 * 60)
> ```

This will return the average number of responses by provider over the last hour.
