# Websocket Metrics

## Overview

The Websocket Metrics package provides a set of metrics that will be implemented by default for all providers that inherit from the Base Provider and implement a websocket-based provider. These metrics are intended to be used by the provider to track the usage of the provider's Websocket APIs and the resources it manages. Specifically, this package tracks various metrics related to the underlying connection as well as the data handler which is responsible for processing the data received from the Websocket connection.

## Websocket Metrics

The following metrics are provided by the Websocket Metrics package for implementations that use the Websocket Handler:

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

The `AddWebSocketConnectionStatus` metric is used to track the number of connection related errors that occur. Specifically, this tracks the number of connection related errors that occur for a given provider. For example, there may be a connection error when attempting to connect to the Websocket API or there may be a connection error that occurs after the connection has been established - write errors, read errors, etc.

### AddWebSocketDataHandlerStatus

The `AddWebSocketDataHandlerStatus` metric is used to track the number of data handler related errors that occur. Specifically, this tracks the number of data handler related errors that occur for a given provider. For example, there may be a data handler error when attempting to process the data received from the Websocket API or constructing an update message to send to the client.

### ObserveWebSocketLatency

The `ObserveWebSocketLatency` metric is used to track the time it took for a provider to respond. Specifically, this tracks how long it takes to successfully receive and process data from the Websocket API. If the response time is very large, this could mean that the provider is not sending data frequently enough or that the data handler is taking too long to process the data.

## Usage

Below we overview some of the more useful prometheus queries that can be used to get insight into the health of a provider.

### Total number of reads, writes, and errors related to the underlying connection

> ```promql
> sum by (status) (increase(oracle_web_socket_connection_status_per_provider[1h]))
> ```

This will return the total number of reads, writes, dials, and health errors pertaining to all of the underlying web socket connections. This provides insight into how reliable the websocket connection is.

### Total number of reads, writes, and errors related to the underlying connection for a given provider

> ```promql
> sum by (status, provider) (increase(oracle_web_socket_connection_status_per_provider[1h]))
> ```

This will return the total number of reads, writes, dials, and health errors pertaining to the underlying web socket connection for a given provider. This provides insight into how reliable the websocket connection is for a given provider.

### Total number of data handler successes and failures

> ```promql
> sum by (status) (increase(oracle_web_socket_data_handler_status_per_provider[1h]))
> ```

This will return the total number of successfully handled messages, heartbeats, subscriptions, and errors pertaining to the data handler. This provides insight into how reliable the data handler is.

### Total number of data handler successes and failures for a given provider

> ```promql
> sum by (status, provider) (increase(oracle_web_socket_data_handler_status_per_provider[1h]))
> ```

This will return the total number of successfully handled messages, heartbeats, subscriptions, and errors pertaining to the data handler for a given provider. This provides insight into how reliable the data handler is for a given provider.
