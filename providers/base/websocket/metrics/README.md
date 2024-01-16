# WebSocket Metrics

## Overview

The WebSocket Metrics package provides a set of metrics that will be implemented by default for all providers that inherit from the Base Provider and implement a Web Socket based provider. These metrics are intended to be used by the provider to track the usage of the provider's Web Socket APIs and the resources it manages. Specifically, this package tracks various metrics related to the underlying connection as well as the data handler which is responsible for processing the data received from the Web Socket connection.

## Web Socket Metrics

The following metrics are provided by the Web Socket Metrics package for implementations that use the Web Socket Handler:

```golang
// WebSocketMetrics is an interface that defines the API for metrics collection for providers
// that implement the WebSocketQueryHandler.
type WebSocketMetrics interface {
	// AddWebSocketConnectionStatus adds a method / status response to the metrics collector for the
	// given provider. Specifically, this tracks various connection related errors.
	AddWebSocketConnectionStatus(provider string, status ConnectionStatus)

	// AddWebSocketDataHandlerStatus adds a method / status response to the metrics collector for the
	// given provider. Specifically, this tracks various data handler related errors.
	AddWebSocketDataHandlerStatus(provider string, status HandlerStatus)

	// ObserveWebSocketLatency adds a latency observation to the metrics collector for the
	// given provider.
	ObserveWebSocketLatency(provider string, duration time.Duration)
}
```

### AddWebSocketConnectionStatus

The `AddWebSocketConnectionStatus` metric is used to track the number of connection related errors that occur. Specifically, this tracks the number of connection related errors that occur for a given provider. For example, there may be a connection error when attempting to connect to the Web Socket API or there may be a connection error that occurs after the connection has been established - write errors, read errors, etc.

### AddWebSocketDataHandlerStatus

The `AddWebSocketDataHandlerStatus` metric is used to track the number of data handler related errors that occur. Specifically, this tracks the number of data handler related errors that occur for a given provider. For example, there may be a data handler error when attempting to process the data received from the Web Socket API or constructing an update message to send to the client.

### ObserveWebSocketLatency

The `ObserveWebSocketLatency` metric is used to track the time it took for a provider to respond. Specifically, this tracks how long it takes to succesfully receive and process data from the Web Socket API. If the response time is very large, this could mean that the provider is not sending data frequently enough or that the data handler is taking too long to process the data.
