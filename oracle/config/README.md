# Configurations

> To read more about how to configure a new provider to your oracle, please read the [Create Your Own Provider](../../providers/README.md) readme. This readme overview's the configuration options for the oracle service and its corresponding metrics intrumentation.
> 
> Validator's running on a network that support's the oracle module (`x/oracle`) **must** run the oracle side-car. Non-validator's can configure their oracle config's to be disabled, and the oracle side-car will not be run.
>
> <div align="center">
> 
> | Type | Oracle Config | Oracle Metrics | App Metrics |
> |----------|:--------:|---------:|--------:|
> | Validator     | **Required**     |    **Recommended**   | **Recommended** |
> | Non-Validator | **Optional**     |    **Optional**   | **Optional** |
> </div>

The oracle configuration file has two parts:

1. The basic oracle configuration paths itself are consumed by the application - `app.toml`. These path's point to the oracle configuration files and are read in by the application at startup. This can accept absolute paths or relative paths.
2. the oracle configuration files that is consumed by the oracle and metrics services - `oracle.toml` / `metrics.toml`. These files contain the configurations for the oracle and metrics services and are read in by the oracle and metrics services at startup.

## `app.toml`

The `app.toml` file is the configuration file that is consumed by the application. This file contains over-arching configurations for your entire Cosmos SDK application, as well as a few new configurations for the oracle. The oracle configurations are as follows:

```toml
# Other configurations

...

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Oracle path is the path for the config file for the oracle.
oracle_path = "{{ .Oracle.OraclePath }}"

# Metrics path is the path for the config file for the metrics.
metrics_path = "{{ .Oracle.MetricsPath }}"

...

# More configurations
```

In your `app.toml`, you should see / implement something that looks like this.

```toml
...

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Oracle path is the path for the config file for the oracle.
oracle_path = "path/to/oracle.toml"

# Metrics path is the path for the config file for the metrics.
metrics_path = "path/to/metrics.toml"

...
```

## `oracle.toml`

The `oracle.toml` file is the configuration file that is consumed by the oracle service. This file contains configurations for the oracle service, as well as configurations for the providers that are registered with the application & oracle. To add your new provider to the oracle configuration, you must add a new entry to the `providers` array. The `name` field must match the name of the provider that you registered with the application. The `path` field must point to the configuration file for the provider. Please keep reading to learn more about the oracle configuration.


# Oracle Configuration

The main oracle configuration object is located in [oracle.go](oracle.go). This is utilized to set up the oracle, whether in process or out of process, and to configure the providers that the oracle will use. The object is defined as follows:

```go
type OracleConfig struct {
	// Enabled specifies whether the side-car oracle needs to be run.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
    
	// InProcess specifies whether the oracle configured, is currently running as a remote grpc-server, or will be run in process
	InProcess bool `mapstructure:"in_process" toml:"in_process"`
	
    // RemoteAddress is the address of the remote oracle server (if it is running out-of-process)
	RemoteAddress string `mapstructure:"remote_address" toml:"remote_address"`
	
    // Timeout is the time that the client is willing to wait for responses from the oracle before timing out.
	Timeout time.Duration `mapstructure:"timeout" toml:"timeout"`
	
    // UpdateInterval is the interval at which the oracle will fetch prices from providers
	UpdateInterval time.Duration `mapstructure:"update_interval" toml:"update_interval"`
	
    // Providers is the list of providers that the oracle will fetch prices from.
	Providers []ProviderConfig `mapstructure:"providers" toml:"providers"`
	
    // CurrencyPairs is the list of currency pairs that the oracle will fetch prices for.
	CurrencyPairs []oracletypes.CurrencyPair `mapstructure:"currency_pairs" toml:"currency_pairs"`
}
```

## Enabled

This flag is utilized to note whether the oracle side car needs to be run. **In the case of staked validator, this must be set to true with correct configurations for all remaining fields.** In the case of a non-validator, this can be set to false, and the oracle side car will not be run.

## InProcess

This flag is utilized to note whether the oracle side car is running in process or out of process. If this flag is set to true, then the oracle side car will be run in process, and the remaining fields will be used to configure the oracle. If this flag is set to false, then the `Providers`, `CurrencyPairs`, and `UpdateInterval` fields will be unused as these *must* be configured on the out-of-process oracle that is running.

## RemoteAddress

This field is utilized to set the address of the remote oracle server. This field is only used if the oracle is running out of process. If the oracle is running in process, then this field is not used. The out-of-process oracle server must be correctly configured with the correct providers and currency pairs.

## Timeout

This field is utilized to set the timeout for the oracle client. This is the time that the oracle client will wait for a response from the oracle server before timing out i.e. how long will [`ExtendVote`](../../abci/ve/vote_extension.go) wait for a response from the oracle server before timing out.

## UpdateInterval

This field is utilized to set the interval at which the oracle will fetch prices from providers. In the case where a provider fails to respond, the oracle will wait until the next update interval to fetch prices from the provider again.

## Providers

This field is utilized to set the list of providers that the oracle will fetch prices from. Please read the [Create Your Own Provider](../../providers/README.md) readme to learn more about how to configure a provider.

## CurrencyPairs

This field is utilized to set the list of currency pairs that the oracle will fetch prices for. These should be the same exact currency pairs that the oracle module (`x/oracle`) is currently configured to accept.

## Example Configuration

```toml
###############################################################################
###                               Oracle                                    ###
###############################################################################
# OracleConfig TOML Configuration

# Enabled specifies whether the oracle side-car needs to be run.
enabled = true

# InProcess specifies whether the oracle is currently running in-process (true) or out-of-process (false).
in_process = true

# Timeout is the time that the client is willing to wait for responses from the oracle.
timeout = "2s"  # Replace "2s" with your desired timeout duration.

# RemoteAddress is the address of the remote oracle server (if it is running out-of-process).
remote_address = ""

# UpdateInterval is the interval at which the oracle will fetch prices from providers.
update_interval = "2s"  # Replace "2s" with your desired update interval duration.

# Providers is the list of providers that the oracle will fetch prices from.

[[providers]]
name = "coinbase"
path = "config/local/providers/coinbase.toml"

[[providers]]
name = "coingecko"
path = "config/local/providers/coingecko.toml"

[[providers]]
name = "coinmarketcap"
path = "config/local/providers/coinmarketcap.toml"

# Currency Pairs

[[currency_pairs]]
base = "BITCOIN"
quote = "USD"
```

# Instrumentation Configuration

> **This section only applies if you are running a validator and the oracle side-car.** If you are running a non-validator node, you can skip this section.

There two types of oracle instrumentations that can be configured: **app-side oracle metrics** and **oracle side-car metrics**. The metrics configuration can be found in [metrics.go](metrics.go). The metrics configuration object is defined as follows:

```go
type MetricsConfig struct {
	// PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
	PrometheusServerAddress string `mapstructure:"prometheus_server_address" toml:"prometheus_server_address"`
	// OracleMetrics is the config for the oracle metrics
	OracleMetrics oracle_metrics.Config `mapstructure:"oracle_metrics" toml:"oracle_metrics"`
	// AppMetrics is the config for the app metrics
	AppMetrics service_metrics.Config `mapstructure:"app_metrics" toml:"app_metrics"`
}
```

## PrometheusServerAddress

This field is utilized to set the address of the prometheus server that will expose these metrics.

## App Metrics

The oracle's `VoteExtensionHandler` and `PreBlockHandler` both utilize the `OracleMetrics` object to instrument how the oracle is interacting with the base app. This allow's us to internally track how long it takes to fetch prices from the oracle service, how frequently a given validator's prices are included in blocks, and more. The `OracleMetrics` object is defined as follows:

```go
type Config struct {
	// Enabled indicates whether metrics should be enabled
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}
```

### Enabled

This flag is utilized to note whether the app-side oracle metrics should be enabled.

## Oracle Metrics

The oracle side-car metrics track's various information about the health of the side-car, including how frequently providers are successfully updating prices, how frequently the oracle is successfully updating prices, and more. The `OracleMetrics` object is defined as follows:

```go
type Config struct {
	// Enabled indicates whether metrics should be enabled
	Enabled bool `mapstructure:"enabled" toml:"enabled"`
}
```

### Enabled

This flag is utilized to note whether the oracle-side car metrics should be enabled.


## Example Configuration

```toml
###############################################################################
###                              Metrics                                    ###
###############################################################################
# MetricsConfig TOML Configuration

# PrometheusServerAddress is the address of the Prometheus server that the oracle will expose metrics to.
prometheus_server_address = "localhost:8000"

# OracleMetrics is the config for the oracle metrics
[oracle_metrics]

# Enabled indicates whether metrics should be enabled
enabled = true

# AppMetrics is the config for the app metrics
[app_metrics]
# Enabled indicates whether metrics should be enabled
enabled = true
```
