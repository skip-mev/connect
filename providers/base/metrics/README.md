# Base Provider Metrics

## Overview

The Base Provider Metrics package provides a set of metrics that will be implemented by default for all providers that inherit from the Base Provider. These metrics are intended to be used to get insight into how often providers are successfully updating the data they are responsible for managing.

## Provider Metrics

The following metrics are provided by the Base Provider Metrics package:

```golang
// ProviderMetrics is an interface that defines the API for metrics collection for providers. The
// base provider utilizes this interface to collect metrics, whether the underlying implementation
// is API or websocket based.
type ProviderMetrics interface {
	// AddProviderResponseByID increments the number of ticks with a fully successful provider update
	// for a given provider and ID (i.e. currency pair).
	AddProviderResponseByID(providerName, id string, status Status)

	// AddProviderResponse increments the number of ticks with a fully successful provider update.
	AddProviderResponse(providerName string, status Status)

	// LastUpdated updates the last time a given ID (i.e. currency pair) was updated.
	LastUpdated(providerName, id string)
}
```

### AddProviderResponseByID

The `AddProviderResponseByID` metric is used to track the number of ticks with a fully successful provider update for a given provider and ID (i.e. currency pair). This metric provides direct introspection into every data source (i.e. price feed) that the provider is responsible for managing.

### AddProviderResponse

The `AddProviderResponse` metric is used to track the number of ticks with a fully successful provider update. This metric provides a high level view of the provider's overall health.

### LastUpdated

The `LastUpdated` metric is used to track the last time a given ID (i.e. currency pair) was updated. This metric provides direct introspection into every data source (i.e. price feed) that the provider is responsible for managing.

## Usage

Below we overview some of the more useful prometheus queries that can be used to get insight into the health of a provider.

### Total number of responses within a time window

> ```promql
> sum(increase(oracle_provider_status_responses[1h]))
> ```

This will return the total number of provider responses over the last hour (failures and successes across all providers).


### Number of Provider Responses by status (success/failure) within a time window

> ```promql
> sum by (status) (increase(oracle_provider_status_responses[1h]))
> ```

This will return the total number of provider responses by status (i.e. success or failure) over the last hour. This provides introspection into how often providers are successfully updating their data.

### Number of Provider Response by ID (i.e. price feed) and status (success/failure) within a time window

> ```promql
> sum by (status, id) (increase(oracle_provider_status_responses_per_id[1h]))
> ```

This will return the number of provider responses by ID (i.e. price feed) and status (i.e. success or failure) over the last hour. This provides introspection into how often each price feed is being updated successfully. If a provider is not returning an expected amount of successful responses for a given price feed, this could be an indication that the provider is not properly configured.

### Average number of responses within a time window

> ```promql
> sum(increase(oracle_provider_status_responses[1h])) / (60 * 60)
> ```

### Average number of responses per provider and status (success/failure) within a time window

> ```promql
> sum by (provider, status) (increase(oracle_provider_status_responses[1h])) / (60 * 60)
> ```

This will return the average number of responses per provider over the last hour. This provides introspection into how often each provider is successfully updating their data. If a provider is not returning an expected amount of successful responses, this could be an indication that the provider is not properly configured.

### Average number of responses per ID (i.e. price feed) and status (success/failure) within a time window

> ```promql
> sum by (id, status) (increase(oracle_provider_status_responses_per_id[1h])) / (60 * 60)
> ```

This will give the average number of responses per ID (i.e. price feed) over the last hour. This provides introspection into how often each price feed is being updated. If certain price feeds are not being updated as often as others, this could be an indication that the provider is not properly configured to manage that price feed.

### Average number of responses by provider type (API/Websocket)

> ```promql
> avg by (type) (increase(oracle_provider_status_responses[24h])) / (24 * 60 * 60)
> ```

This will give the average number of responses by provider type (i.e. API or websocket) over the last 24 hours. This provides introspection into how much more performant websockets are compared to API based providers. 




