# Configurations

> This readme overviews how to configure the oracle side-car as well as how to hook it up to your application. To see an example of a properly configured oracle-side car, please visit the [local config](./../../config/local) files - `oracle.json` and `market.json`. To see an example of a properly configured application, please visit the test application's [app.toml](./../../tests/simapp/slinkyd/testappd/root.go) generation code. Otherwise, please read on to learn how to configure the oracle side-car and application.
>
> To generate custom side-car configs, please use the `slinky-config` binary after running `make build`. Run `./build/slinky-config --help` to determine the relevant fields that can be set.
> 
> Validator's running on a network that support's Slinky **must** run the oracle side-car and configure it into their application. Non-validator's can configure their oracle config's to be disabled, and the oracle side-car will not be run.
>
> <div align="center">
> 
> | Type | Oracle/Market Config | Oracle Metrics | App Metrics |
> |----------|:--------:|---------:|--------:|
> | Validator     | **Required**     |    **Recommended**   | **Recommended** |
> | Non-Validator | **Optional**     |    **Optional**   | **Optional** |
> </div>

All oracle configurations are broken down into three files:

1. **Oracle side-car configuration (`oracle.json`):** This contains the data provider's that are utilized, how often they should be polled, and a variety of other configurations for API and web socket providers.
2. **Market side-car configuration (`market.json`):** This contains the desired markets that the side-car will fetch prices for. NOTE: It is recommended that this file is **NOT** modified nor created by validators. This file is typically provided by the chain that the oracle supports.
3. **Oracle configuration in the application (`app.toml`):** A few additional lines of code that must be added to the application's `app.toml` file to configure the oracle side car into the application.

*The focus of this readme is the oracle side-car configuration and the application configuration. The market side-car configuration is typically provided by the chain that the oracle supports.*

# App Configuration

The `app.toml` file is the configuration file that is consumed by the application. This file contains over-arching configurations for your entire Cosmos SDK application, as well as a few new configurations for the oracle. You must use this template to add the oracle configurations to your `app.toml` file:

```toml
# Other configurations

...

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Enabled indicates whether the oracle is enabled.
enabled = "{{ .Oracle.Enabled }}"

# Oracle Address is the URL of the out of process oracle sidecar. This is used to
# connect to the oracle sidecar when the application boots up. Note that the address
# can be modified at any point, but will only take effect after the application is
# restarted. This can be the address of an oracle container running on the same
# machine or a remote machine.
oracle_address = "{{ .Oracle.OracleAddress }}"

# Client Timeout is the time that the client is willing to wait for responses from 
# the side-car before timing out.
client_timeout = "{{ .Oracle.ClientTimeout }}"

# MetricsEnabled determines whether oracle metrics are enabled. Specifically
# this enables instrumentation of the side-car client and the interaction between
# the side-car and the app.
metrics_enabled = "{{ .Oracle.MetricsEnabled }}"

# PrometheusServerAddress is the address of the prometheus server that metrics will be
# exposed to.
prometheus_server_address = "{{ .Oracle.PrometheusServerAddress }}"

...

# More configurations
```

In your `app.toml`, you should see / write something that looks like this.

> Note: This is only required if you are running a validator node. If you are running a non-validator node, you can skip this section.

```toml
...


###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Enabled indicates whether the oracle is enabled.
enabled = "true"

# Oracle Address is the URL of the out of process oracle sidecar. This is used to
# connect to the oracle sidecar when the application boots up. Note that the address
# can be modified at any point, but will only take effect after the application is
# restarted. This can be the address of an oracle container running on the same
# machine or a remote machine.
oracle_address = "0.0.0.0:8080"

# Client Timeout is the time that the client is willing to wait for responses from 
# the oracle before timing out.
client_timeout = "1s"

# MetricsEnabled determines whether oracle metrics are enabled. Specifically
# this enables intsrumentation of the oracle client and the interaction between
# the oracle and the app.
metrics_enabled = "true"

# PrometheusServerAddress is the address of the prometheus server that metrics will be
# exposed to.
prometheus_server_address = "0.0.0.0:8001"

...
```

# Oracle Side-Car Configuration

The `oracle.json` file is the configuration file that is consumed by the oracle side-car. Note that in most cases, this should **NOT** be custom made by validators - unless specified otherwise. A predefined oracle side car configuration should be provided by the chain that the oracle supports. This file contains:

* The desired data providers to be utilized i.e. Coinbase, Binance, etc.
* Metrics instrumentation.
* API & WebSocket configurations for each provider.

In some cases, validators must configure the market map provider into their `oracle.json`. The market map provider is a special provider that provides the desired markets that the oracle should fetch prices for. This is particularly useful for chains that have a large number of markets that are constantly changing. The market map provider allows the oracle to be updated with new markets without needing to restart the side-car. **Please check the relevant chain's documentation & channels to determine if you need to configure the market map provider.**


## Oracle Configuration

The main oracle configuration object is located in [oracle.go](oracle.go). This is utilized to set up the oracle and to configure the providers that the oracle will use. The object is defined as follows:

```go
type OracleConfig struct {
	UpdateInterval time.Duration    `json:"updateInterval"`
	MaxPriceAge    time.Duration    `json:"maxPriceAge"`
	Providers      []ProviderConfig `json:"providers"`
	Production     bool             `json:"production"`
	Metrics        MetricsConfig    `json:"metrics"`
	Host           string           `json:"host"`
	Port           string           `json:"port"`
}
```

## UpdateInterval

This field is utilized to set the interval at which the oracle will aggregate price feeds from price providers.

## MaxPriceAge

This field is utilized to set the maximum age of a price that the oracle will consider when aggregating prices. If a price is older than this value, the oracle will not consider it when aggregating prices.

## Providers

This field is utilized to set the list of providers that the oracle will fetch prices from. A given provider's configuration is composed of:

* An API configuration that defines the various API configurations that the oracle will use to fetch prices from the provider.
* A WebSocket configuration that defines the various WebSocket configurations that the oracle will use to fetch prices from the provider.
* A type that defines the type of provider - specifically this is currently either a price or market map provider. Price providers supply price feeds for a given set of markets, while a market map provider supplies the desired markets that need to be fetched. Market map providers allow the side-car to be updated with new markets without needing to restart the side-car.

> Note: Typically only one of either the API or websocket config is required. However, some providers may require both. Please read the provider's documentation to learn more about how to configure the provider. Each provider provides sensible defaults for the API and WebSocket configurations that should be used for most cases. This should be modified with caution.

```go
type ProviderConfig struct {
	Name      string          `json:"name"`
	API       APIConfig       `json:"api"`
	WebSocket WebSocketConfig `json:"webSocket"`
	Type      string          `json:"type"`
}
```

### Name

This field is utilized to set the name of the provider. This name is used to identify the provider in the oracle's logs as well as in the oracle's metrics.

### API

This field is utilized to set the various API configurations that are specific to the provider.

```go
type APIConfig struct {
	Enabled          bool          `json:"enabled"`
	Timeout          time.Duration `json:"timeout"`
	Interval         time.Duration `json:"interval"`
	ReconnectTimeout time.Duration `json:"reconnectTimeout"`
	MaxQueries       int           `json:"maxQueries"`
	Atomic           bool          `json:"atomic"`
	URL              string        `json:"url"`
	Name             string        `json:"name"`
}
```

#### Enabled (API)

This field is utilized to set whether the provider is API based. If the provider is not API based, this field should be set to `false`.

#### Timeout

This field is utilized to set the amount of time the provider should wait for a response from its API before timing out.

#### Interval

This field is utilized to set the interval at which the provider should update the prices. Note that provider's may rate limit based on this interval so it is recommended to tune this value as necessary.

#### MaxQueries

This field is utilized to set the maximum number of queries that the provider will make within the interval.

#### Atomic

This field is utilized to set whether the provider can fulfill its queries in a single request. If the provider can fulfill its queries in a single request, this field should be set to `true`. Otherwise, this field should be set to `false`. In the case where all requests can be fulfilled atomically, the oracle will make a single request to the provider to fetch prices for all currency pairs once every interval. 

#### URL

This field is utilized to set the URL that is used to fetch data from the API.

#### Name (Should be the same as the provider's name)

This field is utilized to set the name of the provider. Mostly used as a sanity check to ensure the API configurations correctly correspond to the provider.

### WebSocket

This field is utilized to set the various WebSocket configurations that are specific to the provider.

```go
type WebSocketConfig struct {
	Enabled                       bool          `json:"enabled"`
	MaxBufferSize                 int           `json:"maxBufferSize"`
	ReconnectionTimeout           time.Duration `json:"reconnectionTimeout"`
	WSS                           string        `json:"wss"`
	Name                          string        `json:"name"`
	ReadBufferSize                int           `json:"readBufferSize"`
	WriteBufferSize               int           `json:"writeBufferSize"`
	HandshakeTimeout              time.Duration `json:"handshakeTimeout"`
	EnableCompression             bool          `json:"enableCompression"`
	ReadTimeout                   time.Duration `json:"readTimeout"`
	WriteTimeout                  time.Duration `json:"writeTimeout"`
	PingInterval                  time.Duration `json:"pingInterval"`
	MaxReadErrorCount             int           `json:"maxReadErrorCount"`
	MaxSubscriptionsPerConnection int           `json:"maxSubscriptionsPerConnection"`
}
```

#### Enabled (Websocket)

This field is utilized to set whether the provider is WebSocket based. If the provider is not WebSocket based, this field should be set to `false`.

#### MaxBufferSize

This field is utilized to set the maximum number of messages that the provider will buffer at any given time. If the provider receives more messages than this, it will block receiving messages until the buffer is cleared.

#### ReconnectionTimeout

This field is utilized to set the timeout for the provider to attempt to reconnect to the websocket endpoint. In the case when the connection is corrupted, the provider will wait the `ReconnectionTimeout` before attempting to reconnect.

#### WSS

This field is utilized to set the websocket endpoint for the provider.

#### Name (Should match the provider's name)

This field is utilized to set the name of the provider. Mostly used as a sanity check to ensure the WebSocket configurations correctly correspond to the provider.

#### ReadBufferSize

This field is utilized to set the I/O read buffer size. If a buffer size of 0 is specified, then a default buffer size is used.

#### WriteBufferSize

This field is utilized to set the I/O write buffer size. If a buffer size of 0 is specified, then a default buffer size is used.

#### HandshakeTimeout

This field is utilized to set the duration for the handshake to complete.

#### EnableCompression

This field is utilized to set whether the client should attempt to negotiate per message compression (RFC 7692). Setting this value to true does not guarantee that compression will be supported. Note that enabling compression may increase latency.

#### ReadTimeout

This field is utilized to set the read deadline on the underlying network connection. After a read has timed out, the websocket connection state is corrupt and all future reads will return an error. A zero value for t means reads will not time out.

#### WriteTimeout

This field is utilized to set the write deadline on the underlying network connection. After a write has timed out, the websocket state is corrupt and all future writes will return an error. A zero value for t means writes will not time out.

#### PingInterval

This field is utilized to set the interval to ping the server. Note that a ping interval of 0 disables pings. This is utilized to send heartbeat messages to the server to ensure that the connection is still alive.

#### MaxReadErrorCount

This field is utilized to set the maximum number of read errors that the provider will tolerate before closing the connection and attempting to reconnect.

#### MaxSubscriptionsPerConnection

This field is utilized to set the maximum number of subscriptions that the provider will allow per connection. By default, this value is set to 0, which means that there is no limit to the number of subscriptions that can be made per connection.

## Production

This field is utilized to set whether the oracle is running in production mode. This is used to determine whether the oracle should be run in debug mode or not. This particularly helpful for logging purposes.

## Metrics

This field is utilized to set the metrics configurations for the oracle. To read more about the various metrics that are collected and corresponding queries, please read the [Readme](../../README.md).

```go
type MetricsConfig struct {
	PrometheusServerAddress string `json:"prometheusServerAddress"`
	Enabled                 bool   `json:"enabled"`
}
```

### PrometheusServerAddress

This field is utilized to set the address of the prometheus server that the oracle will expose metrics to.

### Enabled

This field is utilized to set whether metrics should be enabled.

Sample configuration:

```json
{
  "updateInterval": 500000000,
  "maxPriceAge": 120000000000,
  "providers": [
    {
      "name": "binance_api",
      "api": {
        "enabled": true,
        "timeout": 1000000000,
        "interval": 400000000,
        "reconnectTimeout": 2000000000,
        "maxQueries": 1,
        "atomic": true,
        "url": "https://api.binance.com/api/v3/ticker/price?symbols=%s%s%s",
        "name": "binance_api"
      },
      "webSocket": {
        "enabled": false,
        "maxBufferSize": 0,
        "reconnectionTimeout": 0,
        "wss": "",
        "name": "",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 0,
        "enableCompression": false,
        "readTimeout": 0,
        "writeTimeout": 0,
        "pingInterval": 0,
        "maxReadErrorCount": 0,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "coinbase_api",
      "api": {
        "enabled": true,
        "timeout": 500000000,
        "interval": 100000000,
        "reconnectTimeout": 2000000000,
        "maxQueries": 5,
        "atomic": false,
        "url": "https://api.coinbase.com/v2/prices/%s/spot",
        "name": "coinbase_api"
      },
      "webSocket": {
        "enabled": false,
        "maxBufferSize": 0,
        "reconnectionTimeout": 0,
        "wss": "",
        "name": "",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 0,
        "enableCompression": false,
        "readTimeout": 0,
        "writeTimeout": 0,
        "pingInterval": 0,
        "maxReadErrorCount": 0,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "kraken_api",
      "api": {
        "enabled": true,
        "timeout": 500000000,
        "interval": 400000000,
        "reconnectTimeout": 2000000000,
        "maxQueries": 1,
        "atomic": true,
        "url": "https://api.kraken.com/0/public/Ticker?pair=%s",
        "name": "kraken_api"
      },
      "webSocket": {
        "enabled": false,
        "maxBufferSize": 0,
        "reconnectionTimeout": 0,
        "wss": "",
        "name": "",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 0,
        "enableCompression": false,
        "readTimeout": 0,
        "writeTimeout": 0,
        "pingInterval": 0,
        "maxReadErrorCount": 0,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "bitfinex_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1000,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://api-pub.bitfinex.com/ws/2",
        "name": "bitfinex_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 0,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "bitstamp_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1024,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://ws.bitstamp.net",
        "name": "bitstamp_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 10000000000,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "bybit_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1000,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://stream.bybit.com/v5/public/spot",
        "name": "bybit_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 15000000000,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "coinbase_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1024,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://ws-feed.exchange.coinbase.com",
        "name": "coinbase_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 5000000000,
        "pingInterval": 0,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "crypto_dot_com_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1024,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://stream.crypto.com/exchange/v1/market",
        "name": "crypto_dot_com_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 0,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "gate_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1000,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://api.gateio.ws/ws/v4/",
        "name": "gate_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 0,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "huobi_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1000,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://api.huobi.pro/ws",
        "name": "huobi_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 0,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "kucoin_ws",
      "api": {
        "enabled": false,
        "timeout": 5000000000,
        "interval": 60000000000,
        "reconnectTimeout": 0,
        "maxQueries": 1,
        "atomic": false,
        "url": "https://api.kucoin.com",
        "name": "kucoin_ws"
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1024,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://ws-api-spot.kucoin.com/",
        "name": "kucoin_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 10000000000,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "mexc_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1000,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://wbs.mexc.com/ws",
        "name": "mexc_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 20000000000,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 20
      },
      "type": "price_provider"
    },
    {
      "name": "okx_ws",
      "api": {
        "enabled": false,
        "timeout": 0,
        "interval": 0,
        "reconnectTimeout": 0,
        "maxQueries": 0,
        "atomic": false,
        "url": "",
        "name": ""
      },
      "webSocket": {
        "enabled": true,
        "maxBufferSize": 1000,
        "reconnectionTimeout": 10000000000,
        "wss": "wss://ws.okx.com:8443/ws/v5/public",
        "name": "okx_ws",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 45000000000,
        "enableCompression": false,
        "readTimeout": 45000000000,
        "writeTimeout": 45000000000,
        "pingInterval": 0,
        "maxReadErrorCount": 100,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "price_provider"
    },
    {
      "name": "dydx_api",
      "api": {
        "enabled": true,
        "timeout": 20000000000,
        "interval": 10000000000,
        "reconnectTimeout": 2000000000,
        "maxQueries": 1,
        "atomic": true,
        "url": "localhost:1317",
        "name": "dydx_api"
      },
      "webSocket": {
        "enabled": false,
        "maxBufferSize": 0,
        "reconnectionTimeout": 0,
        "wss": "",
        "name": "",
        "readBufferSize": 0,
        "writeBufferSize": 0,
        "handshakeTimeout": 0,
        "enableCompression": false,
        "readTimeout": 0,
        "writeTimeout": 0,
        "pingInterval": 0,
        "maxReadErrorCount": 0,
        "maxSubscriptionsPerConnection": 0
      },
      "type": "market_map_provider"
    }
  ],
  "production": true,
  "metrics": {
    "prometheusServerAddress": "0.0.0.0:8002",
    "enabled": true
  },
  "host": "0.0.0.0",
  "port": "8080"
}
```

# Conclusion

This readme has provided an overview of how to configure the oracle side-car and application. It has also provided a brief overview of the oracle side-car configuration and the application configuration. To see an example of a properly configured oracle side car, please visit the [local config](./../../config/local) files - `oracle.json` and `market.json`. 

In general, it is best to consult the chain's documentation and channels to determine the correct configurations for the oracle side-car. If you have any questions, please feel free to reach out to the Skip team on the [Skip Discord](https://discord.com/invite/hFeHVAE26P). 
