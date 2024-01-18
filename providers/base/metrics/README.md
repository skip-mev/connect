# Base Provider Metrics

## Overview

The Base Provider Metrics package provides a set of metrics that will be implemented by default for all providers that inherit from the Base Provider. These metrics are intended to be used to get insight into how often providers are successfully updating the data they are responsible for managing.

## Provider Metrics

The following metrics are provided by the Base Provider Metrics package:

```golang
// ProviderMetrics is an interface that defines the API for metrics collection for providers. The
// base provider utilizes this interface to collect metrics, whether the underlying implementation
// is API or web socket based.
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
