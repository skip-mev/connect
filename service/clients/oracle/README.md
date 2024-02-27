# Oracle Client

## Overview

The oracle client is responsible for fetching data from a price oracle that is running externally from the Cosmos SDK application. The client will fetch prices, standardize them according to the preferences of the application and include them in a validator's vote extension.

```golang
// OracleClient defines the interface that will be utilized by the application
// to query the oracle service. This interface is meant to be implemented by
// the gRPC client that connects to the oracle service.
type OracleClient interface {
	// Prices defines a method for fetching the latest prices.
	Prices(ctx context.Context, in *QueryPricesRequest, opts ...grpc.CallOption) (*QueryPricesResponse, error)

	// Start starts the oracle client.
	Start() error

	// Stop stops the oracle client.
	Stop() error
}
```

There are two types of clients that are supported:

* [**Vanilla GRPC oracle client**](./client.go) - This client is responsible for fetching data from an oracle that is aggregating price data. It implements a GRPC client that connects to the oracle service and fetches the latest prices.
* [**Metrics GRPC oracle client**](./client.go) - This client implements the same functionality as the vanilla GRPC oracle client, but also exposes metrics that can be scraped by Prometheus.

To enable the metrics GRPC client, please read over the [oracle configurations](../../../oracle/config/README.md) documentation.
