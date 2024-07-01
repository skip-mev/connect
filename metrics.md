# Slinky Metrics

This document describes the various instrumentation points that are available in the side-car. This should be utilized to monitor the health of Slinky and the services it is proxying.

If there are any additional metrics that you would like to see, please open an issue in the [GitHub repository](https://github.com/skip-mev/slinky). Note that this document is not fully comprehensive and may be updated in the future. However, it should provide a good starting point for monitoring Slinky. It may be useful to read the [configuration overview](./oracle/config/README.md) to understand how the side-car is configured. However, this is NOT necessary.


# Table of Contents

* [Dashboard](#dashboard)
    * [Set Up](#set-up)
* [Metrics](#metrics)
    * [Health Metrics](#health-metrics)
    * [Prices Metrics](#prices-metrics)
        * [Price Feed Metrics](#price-feed-metrics)
        * [Aggregated Price Metrics](#aggregated-price-metrics)
    * [HTTP Metrics](#http-metrics)
    * [WebSocket Metrics](#websocket-metrics)

# Dashboard

The slinky repo contains a sidecar Grafana dashboard that can be utilized across deployments. The dashboard is designed to provide a high-level overview of the sidecar's health and performance. The dashboard includes various panels that display metrics related to the sidecar's health, price feeds, HTTP requests, and WebSocket connections.

![Dashboard Overview](./assets/dashboard.png)

## Set Up

To set up the sidecar Grafana dashboard, you will need to import the [`side-car-dashboard.json`](grafana/provisioning/dashboards/side-car-dashboard.json) file. This file contains the JSON representation of the dashboard. You can import the dashboard by following these steps:

1. Add the side car prometheus data source to your Grafana instance. This can be done by adding a new data source in your datasources directory. An example can be found [here](grafana/provisioning/datasources/prometheus.yml). Note that the `url` field should be the URL of your sidecar Prometheus instance.
2. Log in to your Grafana instance.
3. Click on the "+" icon in the sidebar and select "Import".
4. Click on "Upload JSON file" and select the `side-car-dashboard.json` file or copy paste the file.

The sidecar dashboard should now be available in your Grafana instance. You can access the dashboard by clicking on the "Dashboards" tab in the sidebar.


# Metrics

> **Definitions**:
>
> * **Market**: A market is a pair of assets that are traded against each other. For example, the BTC-USD market is the market where Bitcoin is traded against the US Dollar.
> * **Price Feed**: A price feed is indexed by a price provider and a market. For example, the Coinbase API provides a price feed for the BTC-USD market.
> * **Price Provider**: A price provider is a service that provides price data for a given market. For example, the Coinbase API is a price provider for the Coinbase markets.
> * **Market Map Provider**: A market map provider is a service that supplies the markets that the side-car needs to fetch data for.

Slinky is composed of several price providers and usually one market map provider. Slinky maintains one configuration file (`oracle.json`) that contains information about how often Slinky should poll the providers and a separate configuration file (`market.json`) that contains information about the markets that Slinky should fetch data for. As mentioned above, price providers provide price feeds, while the market map provider provides the markets that Slinky should fetch data for.

Slinky exposes metrics on the `/metrics` endpoint on port `8002` by default. These metrics are in the Prometheus format and can be scraped by Prometheus or any other monitoring system that supports Prometheus format. The configuration we will be focusing on is the `UpdateInterval` found in the `oracle.json` file. This configuration specifies how often Slinky should aggregate the price feeds from the price providers. The examples below are using an `UpdateInterval` of 500 milliseconds (0.5 seconds). For simplicity, this document will primarily focus on the `rate()` function in prometheus. The `rate()` function calculates the per-second average rate of increase of the time series in the range vector.

## Health Metrics

There are three primary health metrics that are exposed by Slinky:

* (RECOMMENDED) [`side_car_health_check_system_updates_total`](#side_car_health_check_system_updates_total): This metric is a counter that increments every time the side-car updates its internal state. This is a good indicator of the side-car's overall health.
* (RECOMMENDED) [`side_car_health_check_ticker_updates_total`](#side_car_health_check_ticker_updates_total): This metric is a counter that increments every time the side-car updates the price of a given market. This is a good indicator of the overall health of a given market.
* (OPTIONAL) [`side_car_health_check_provider_updates_total`](#side_car_health_check_provider_updates_total): This metric is a counter that increments every time the side-car utilizes a given providers market data. This is a good indicator of the health of a given provider. Note that providers may not be responsible for every market. However, the side-car correctly tracks the number of expected updates for each provider.

Given the additional context above, we should expect the health metrics to be increasing at a rate of 2.0 updates/sec (every 500 milliseconds) for each market and provider. If the rate of updates is lower than expected, this may indicate an issue with the side-car.

### `side_car_health_check_system_updates_total`

This metric should be increasing. Specifically, the rate of this metric should be inversely correlated to the configured `UpdateInterval` in the oracle side-car configuration (`oracle.json`). To check this, you can run the following query in Prometheus:

```promql
rate(side_car_health_check_system_updates_total[5m])
```

![Architecture Overview](./assets/side_car_health_check_system_updates_total_rate.png)

As we can see in the graph above, the rate of updates is close to 2.0 updates/sec (every 500 milliseconds) for each market. This indicates that the side-car is updating its internal state as expected.


### `side_car_health_check_ticker_updates_total`

This should be a monotonically increasing counter for each market. Each market's counter should be relatively close to the `side_car_health_check_system_updates_total` counter.

To verify that the rate of updates for each market is as expected, you can run the following query in Prometheus:

```promql
rate(side_car_health_check_ticker_updates_total[5m])
```

![Architecture Overview](./assets/side_car_health_check_ticker_updates_total_rate.png)

As we can see in the graph above, the rate of updates is close to 2.0 updates/sec (every 500 milliseconds) for each market. This indicates that the side-car is updating the price of each market as expected. 

### `side_car_health_check_provider_updates_total`

This metric should be increasing for each (provider, market) pair. Specifically, the metric includes a `success` label that increments the counter every time a price that was needed by the oracle was available. If `success="true"` this means that at the time the oracle needed a price from a provider (i.e. Binance) it was able to provide it. To verify that the rate of updates for each provider is as expected, you can run the following query in Prometheus:

```promql
rate(side_car_health_check_provider_updates_total{provider="coinbase_api", success="true"}[5m])
```

![Architecture Overview](./assets/side_car_health_check_provider_updates_total_rate.png)

### Health Metrics Summary

In summary, the health metrics should be monitored to ensure that the side-car is updating its internal state, updating the price of each market, and fetching data from the price providers as expected. The rate of updates for each of these metrics should be inversely correlated with the `UpdateInterval` in the oracle side-car configuration. 

For example, if the `UpdateInterval` is set to 500 milliseconds, we should expect to see an update twice every second for each market and provider (rate of 2.0). If the rate of updates for the market or provider is lower than expected, this may indicate an issue with the side-car. 

## Prices Metrics

Slinky exposes various metrics related to market prices. These metrics are useful for monitoring the health of the price feeds and the aggregation process. The rate of updates for these metrics is likely going to be much greater than the health metrics. This is expected as we poll the providers much more frequently than we update all of the prices! As such, the metrics below are not expected to be updated at the same rate as the health metrics. However, they are greatly useful to quickly identify issues with the price feeds.

As to be expected, price feeds for the same market should be relatively close to each other. If there is a large discrepancy between the price feeds, this may indicate an issue with the underlying price providers.

### Price Feed Metrics

The following price feed metrics are available to operators:

* [`side_car_provider_price`](#side_car_provider_price): The last recorded price for a given price feed.
* [`side_car_provider_last_updated_id`](#side_car_provider_last_updated_id): The last UNIX timestamp for a given price feed.

#### `side_car_provider_price`

This metric represents the last recorded price for a given price feed. The metric is indexed by the provider and market (id). For example, if we want to check the last recorded price of the BTC-USD market from the Coinbase API, we can run the following query in Prometheus:

```promql
side_car_provider_price{provider="coinbase_api", id="btc/usd"}
```

![Architecture Overview](./assets/side_car_provider_price_coinbase.png)

Alternatively, if we wanted to check that last recorded prices of the BTC-USD market across all price providers, we can run the following query in Prometheus:

```promql
side_car_provider_price{id="btc/usd"}
```

![Architecture Overview](./assets/side_car_provider_price.png)

#### `side_car_provider_last_updated_id`

This metric represents the last recorded timestamp for a given price feed. The metric is indexed by the provider and market. All prices are UNIX timestamped. For example, if we want to check the last recorded timestamp of the BTC-USD market from the Coinbase API, we can run the following query in Prometheus:

```promql
side_car_provider_last_updated_id{provider="coinbase_api", id="btc/usd"}
```

![Architecture Overview](./assets/side_car_provider_last_updated_id_coinbase.png)

Alerts can be configured based on the age of the last recorded price. For example, if the last recorded price is older than a certain threshold, an alert can be triggered. We recommend a threshold of 5 minutes for most use cases.

### Aggregated Price Metrics

The following aggregated price metrics are available to operators:

* [`side_car_aggregated_price`](#side_car_aggregated_price): The aggregated price for a given market. This price is the result of a median aggregation of all available price feeds for a given market. This is the price clients will see when querying the side-car.

#### `side_car_aggregated_price`

This metric represents the aggregated price for a given market. Prices are aggregated across all available price feeds for a given market. The metric includes the number of decimal places for the price - which can be used to quickly identify if the price is being aggregated correctly. For example, if we want to check the aggregated price of the BTC-USD market, we can run the following query in Prometheus:

```promql
side_car_aggregated_price{id="btc/usd"}
```

![Architecture Overview](./assets/side_car_aggregated_price.png)

As we can see, the price is `6961690515` which after normalizing with the `5` decimals we get `6961690515` / `10` ^ `5` = ~`$69,616.91`. This can also be graphed for a given market to visualize the price over time.

![Architecture Overview](./assets/side_car_aggregated_price_graph.png)

### Prices Metrics Summary

In summary, the price feed metrics should be monitored to ensure that prices look reasonable and are being updated as expected. The 
`side_car_provider_price` metrics can be used to check that the `side_car_aggregated_price` is being calculated correctly. Additionally, alerts can be set up based on the age of the last recorded price to ensure that prices are being updated in a timely manner.

## HTTP Metrics

Slinky exposes various metrics related to HTTP requests made by the side-car - including the number of requests, the response time, and the status code. These metrics can be used to monitor the health of the side-car's HTTP endpoints.

The following HTTP metrics are available to operators:

* [`side_car_api_http_status_code`](#side_car_api_http_status_code): The status codes of the HTTP response made by the side-car.
* [`side_car_api_response_latency_bucket`](#side_car_api_response_latency_bucket): The response latency of the HTTP requests made by the side-car.

### `side_car_api_http_status_code`

This metric represents the status codes of the HTTP responses made by the side-car. For example, if we want to check the status codes of the HTTP responses made by the side-car for Coinbase, we can run the following query in Prometheus:

```promql
side_car_api_http_status_code{provider="coinbase_api"}
```

![Architecture Overview](./assets/side_car_api_http_status_code.png)

The status codes are grouped into `2XX`, `3XX`, `4XX`, and `5XX`. Simple queries and alerts can be configured based on the status codes to ensure that the side-car is responding as expected. In particular, each price and market map API provider configures a `Interval` -  which is how frequently the side-car will poll the provider - and a `MaxQueries` field - which is the maximum number of queries the side-car will make to the provider in a given interval. These configurations can be used to set up alerts based on the number of queries made to the provider.

### `side_car_api_response_latency_bucket`

This metric represents the response latency of the HTTP requests made by the side-car. The metric is indexed by the provider. For example, if we want to check the response latency of the HTTP requests made by the side-car for Coinbase, we can run the following query in Prometheus:

```promql
side_car_api_response_latency_bucket{provider="coinbase_api"}
```

![Architecture Overview](./assets/side_car_api_response_latency_bucket.png)

This can be used to monitor the response time of the side-car's HTTP endpoints and set up alerts based on the response time. In particular, each provider configures a `Timeout` - which is the maximum amount of time the side-car will wait for a response from the provider. This configuration can be used to set up alerts based on the response time of the HTTP requests. If the timeout is consistently exceeded, it may indicate that it should be increased.

### HTTP Metrics Summary

In summary, the HTTP metrics should be monitored to ensure that the side-car's HTTP endpoints are responding as expected. The `side_car_api_http_status_code` metrics can be used to check the status codes of the HTTP responses, and the `side_car_api_response_latency_bucket` metrics can be used to monitor the response time of the HTTP requests. If you are seeing several `4XX` or `5XX` status codes, this may indicate an issue with the side-car or the price provider (may require a URL change). If the response time exceeds the timeout, this may indicate that the timeout should be increased.

## WebSocket Metrics

Slinky exposes various metrics related to WebSocket connections made by the side-car. These metrics can be used to monitor the health of the side-car's WebSocket connections. The following WebSocket metrics are available to operators:

* [`side_car_web_socket_connection_status`](#side_car_web_socket_connection_status): This includes various metrics related to the WebSocket connections made by the side-car.
* [`side_car_web_socket_data_handler_status`](#side_car_web_socket_data_handler_status): This includes various metrics related to whether WebSocket messages are being correctly handled by the side-car.
* [`side_car_web_socket_response_time_bucket`](#side_car_web_socket_response_time_bucket): This includes the response time of the WebSocket messages received by the side-car.

### `side_car_web_socket_connection_status`

This metric includes various metrics related to the WebSocket connections made by the side-car. Specifically, this includes the number of reads, writes, and dials for each connection. For example, if we wanted to check these metrics for the Coinbase WebSocket connection, we can run the following query in Prometheus:

```promql
side_car_web_socket_connection_status{provider="coinbase_ws"}
```

![Architecture Overview](./assets/side_car_web_socket_connection_status.png)

The most important statuses to monitor here are `healthy`, `read_success`, `dial_success`, and `write_success`. The `healthy` metric in particular increments every time the side-car establishes and maintains a healthy connection. If the connection is ever unhealthy, you should see an increase in the `unhealthy` label.

### `side_car_web_socket_data_handler_status`

This metric includes various metrics related to whether WebSocket messages are being correctly handled by the side-car. Specifically, this includes the number of messages that were correctly handled, how many heartbeats were sent, and more. For example, if we wanted to check these metrics for the Kucoin WebSocket connection, we can run the following query in Prometheus:

```promql
side_car_web_socket_data_handler_status{provider="kucoin_ws"}
```

![Architecture Overview](./assets/side_car_web_socket_data_handler_status.png)

The most important statuses to monitor here are `handle_message_success` and `heart_beat_success`. This metrics should be increasing over time.

### `side_car_web_socket_response_time_bucket`

This metric includes the response time of the WebSocket messages received by the side-car. Specifically, this includes the time it took to receive a new message and process it. For example, if we wanted to check the response time for the Kucoin WebSocket connection, we can run the following query in Prometheus:

```promql
side_car_web_socket_response_time_bucket{provider="kucoin_ws"}
```

![Architecture Overview](./assets/side_car_web_socket_response_time_bucket.png)

This can be used to monitor the response time of the WebSocket messages received by the side-car and set up alerts based on the response time. We recommend alerts be set up if the response time exceeds a threshold of 5 minutes.

### WebSocket Metrics Summary

In summary, the WebSocket metrics should be monitored to ensure that the side-car's WebSocket connections are functioning as expected. The `side_car_web_socket_connection_status` metrics can be used to check the number of read, write, and dial errors, the `side_car_web_socket_data_handler_status` metrics can be used to check that messages are being correctly handled, and the `side_car_web_socket_response_time` metrics can be used to monitor the response time of the WebSocket messages.

# Conclusion

This document has provided an overview of the various metrics that are available in the side-car. These metrics can be used to monitor the health of the side-car and the services it is proxying. By monitoring these metrics, operators can ensure that the side-car is functioning as expected and take action if any issues arise.


