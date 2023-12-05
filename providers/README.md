# Creating a Provider

## Overview

This readme will walk you through the process of creating a new provider that can be utilized by your application. The process is broken down into the following steps:

1. Determining the provider you want i.e. Binance, Coinbase, etc.
2. Implementing the interface for the provider.
3. Creating a configuration file that will be consumed by the provider.
4. Instantiating the provider with the application.

## 1. Determining the Provider

A few considerations that should be made when determining the provider you want to use include:

* Does the provider require any API keys?
* What are the different types of responses that the provider returns? How do you want to handle these responses?
* Are the prices returned in a format that is compatible with the application? What type of parsing is required to get the price in the correct format?
* What is the frequency of the price updates? How does this compare to the frequency of the price updates in the application?
* Does the provider do any aggregation of prices? If so, is it a TWAP, VWAP, or some other type of aggregation?
* Does the provider have any rate limits? If so, how can you handle these rate limits?
* How many decimals does the provider return? How does this compare to the number of decimals that the application supports or requires?

## 2. Implementing the Interface

The interface required to implement a provider is defined in the `oracle/provider.go` file. The interface is as follows:

```go
// Provider defines an interface an exchange price provider must implement.
type Provider interface {
	// Name returns the name of the provider.
	Name() string

	// GetPrices returns the aggregated prices based on the provided currency pairs.
	GetPrices(context.Context) (map[types.CurrencyPair]aggregator.QuotePrice, error)

	// SetPairs sets the pairs that the provider should fetch prices for.
	SetPairs(...types.CurrencyPair)

	// GetPairs returns the pairs that the provider is fetching prices for.
	GetPairs() []types.CurrencyPair
}
```

The `Name` function returns the name of the provider. This name is used to identify the provider in the configuration file. The `GetPrices` function returns a map of currency pairs to price. The `SetPairs` function sets the currency pairs that the provider should fetch prices for. The `GetPairs` function returns the currency pairs that the provider is fetching prices for.

> **CAUTION: It is critical that the amount of decimal points returned for each currency pair matches that of what is stored on-chain.**
> 
> For example, if the `x/oracle` module is storing the price of `ETH/USDT` with 8 decimal points, then the provider must return the price of `ETH/USDT` with 8 decimal points. If the provider returns the price of `ETH/USDT` with 18 decimal points then **this may result in in-correct prices being stored on-chain and ultimately may result in slashing of a validator**. 
>
> Please consult the `x/oracle` module documentation for more information on how prices are stored on-chain as well as chain maintainers for the number of decimal points that are required for each currency pair.

The oracle will make call `GetPrices` to the provider every so often - configured as `UpdateInterval` on the oracle. The provider is responsible for returning a response for every supported currency pair in the correct format. 

## 3. Creating the Configuration File

Based on the considerations mentioned above, your provider may require additional configurations - such as an API key - in order to provide reliable price data. To that, you will need to create a configuration file that will be consumed by the provider. The configuration file can be of any format - including `toml`, `json`, `yaml`, etc. - as long as it can be parsed by the provider. We recommend that all code related to the same provider be placed in the same directory. For example, the `coinmarketcap` provider has the following directory structure:

```bash
├── providers/
│   ├── coinmarketcap/
│   │   ├── provider.go 
│   │   ├── config.go
│   │   ├── utils.go
│   │   ├── ...

```

The `config.go` file contains the configuration for the provider. The configuration is defined as follows:

```go
// Config is the config struct for the coinmarketcap provider.
type Config struct {
	// APIKey is the API key used to make requests to the coinmarketcap API.
	APIKey string `mapstructure:"api_key" toml:"api_key"`
	// TokenNameToMetadata is a map of token names to their metadata.
	TokenNameToSymbol map[string]string `mapstructure:"token_name_to_symbol" toml:"token_name_to_symbol"`
}
```

For this specific provider, we need an API key that allows us to have consistent access to their APIs without worries about rate limiting. Additionally, the provider needs a way to map currency pairs - ex. (ETHEREUM/USDT) - to the symbols that are required for their APIs. 


## 4. Instantiating the Provider in BaseApp

Now that we have a provider that implements the interface and a configuration file that can be consumed by the provider, we can instantiate the provider with the application.

### 4.1. Registering the Provider

To register the provider with the base application, you must either create your own `ProviderFactory` or add your provider to the `DefaultProviderFactory` provided by your application. The `ProviderFactory` is defined as follows:


```go
// ProviderFactory inputs the oracle configuration and returns a set of providers. Developers
// can implement their own provider factory to create their own providers.
type ProviderFactory func(log.Logger, config.OracleConfig) ([]Provider, error)
```

Here is a sample implementation of a `ProviderFactory`:

```go
// DefaultProviderFactory returns a sample implementation of the provider factory.
func DefaultProviderFactory() oracle.ProviderFactory {
	return func(logger log.Logger, oracleCfg config.OracleConfig) ([]oracle.Provider, error) {
		providers := make([]oracle.Provider, len(oracleCfg.Providers))

		var err error
		for i, p := range oracleCfg.Providers {
			if providers[i], err = providerFromProviderConfig(logger, oracleCfg.CurrencyPairs, p); err != nil {
				return nil, err
			}
		}

		return providers, nil
	}
}

// providerFromProviderConfig returns a provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func providerFromProviderConfig(logger log.Logger, cps []types.CurrencyPair, cfg config.ProviderConfig) (oracle.Provider, error) {
	switch cfg.Name {
	case "coingecko":
		return coingecko.NewProvider(logger, cps, cfg)
	case "coinbase":
		return coinbase.NewProvider(logger, cps, cfg)
	case "coinmarketcap":
		return coinmarketcap.NewProvider(logger, cps, cfg)
	case "erc4626":
		return erc4626.NewProvider(logger, cps, cfg)
	case "erc4626-share-price-oracle":
		return erc4626sharepriceoracle.NewProvider(logger, cps, cfg)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
}
```

Register this with your `app.go` file as follows:

```go
...
// Create the oracle service.
app.oracleService, err = serviceclient.NewOracleService(
    app.Logger(),
    oracleCfg,
    metricsCfg,
    DefaultProviderFactory(), // Register the provider factory here.
    aggregator.ComputeMedian(),
)
if err != nil {
    panic(err)
}
...
```

Now your provider can be resolved by the application!

### 4.2. Registering the Provider in your oracle configuration

Now that the provider is registered with the application, you must configure your provider within your oracle configuration file that is consumed by the application. The oracle configuration file has two parts:

1. The basic oracle configuration paths itself that is consumed by the application - `app.toml`.
2. the oracle configuration file that is consumed by the oracle service - `oracle.toml` \ `metrics.toml`.

#### `app.toml`

The `app.toml` file is the basic configuration file that is consumed by the application. This file contains over-arching configurations for your entire Cosmos SDK application, as well as a few new configurations for the oracle. The oracle configurations are as follows:

```toml
###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Oracle path is the path for the config file for the oracle.
oracle_path = "{{ .Oracle.OraclePath }}"

# Metrics path is the path for the config file for the metrics.
metrics_path = "{{ .Oracle.MetricsPath }}"
```

In your `app.toml`, you will see something that looks like this.

```toml
...

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
# Oracle path is the path for the config file for the oracle.
oracle_path = "config/path/to/oracle.toml"

# Metrics path is the path for the config file for the metrics.
metrics_path = "config/path/to/metrics.toml"

...
```

#### `oracle.toml`

The `oracle.toml` file is the configuration file that is consumed by the oracle service. This file contains configurations for the oracle service, as well as configurations for the providers that are registered with the application. The oracle configurations are as follows:

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

To add your new provider to the oracle configuration, you must add a new entry to the `providers` array. The `name` field must match the name of the provider that you registered with the application. The `path` field must point to the configuration file for the provider.

### 4.3. Managing the Provider Configuration

We recommend that operators maintain a separate directory for all oracle configurations. The directory structure may look like the following:

```bash
├── config/
│   ├── local/
│   │   ├── oracle.toml
│   │   ├── metrics.toml
│   │   ├── providers/
│   │   │   ├── coinbase.toml
│   │   │   ├── coingecko.toml
│   │   │   ├── coinmarketcap.toml
│   │   │   ├── ...
```

This directory should then be referenced by the `oracle_path` and `metrics_path` fields in the `app.toml` file.
