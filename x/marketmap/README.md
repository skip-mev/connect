# `x/marketmap`

## Contents

* [Concepts](#concepts)
* [Integration](#integtration)
* [State](#state)
    * [MarketMap](#marketmap)
    * [Params](#params)
        * [MarketAuthority](#marketauthority)
        * [Version](#version)
* [Events](#events)
* [Hooks](#hooks)
    * [AfterMarketCreated](#aftermarketcreated)
    * [AfterMarketUpdated](#aftermarketupdated)
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)

## Concepts

The `x/marketmap` module encapsulates a system for creating and updating a unified configuration that is stored on-chain
and consumed by a set of oracle service providers (Connect oracle, etc.).

The core goal of the system is to collect off-chain market updates and to post them on chain, informing oracle service
providers to fetch prices for new markets.

The data is stored in a `MarketMap` data structure which can be queried and consumed by oracle services.

## Integration

When integrating `x/marketmap` into your Cosmos SDK application, some considerations must be made:

### Module Hooks

Integrating modules can use the [hooks](#hooks) exposed by `x/marketmap` to update their state whenever
changes are made to the marketmap.

An example of this can be seen in `x/oracle`'s implementation of the `AfterMarketCreated` hook.  This hook
triggers the creation of a `CurrencyPairState` that corresponds to the new `Ticker` that was created in the marketmap.
This allows for a unified flow where updates to the market map prepare the `x/oracle` module for new price feeds.

### Genesis Order

Any modules that integrate with `x/marketmap` must set their `InitGenesis` to occur _before_ the `x/marketmap` module's
`InitGenesis`.  This is so that logic any consuming modules may want to implement in `AfterMarketGenesis` will be
run properly.

## State

### MarketMap

The market map data is as follows:

```protobuf
// Market encapsulates a Ticker and its provider-specific configuration.
message Market {
  option (gogoproto.goproto_stringer) = false;
  option (gogoproto.stringer) = false;

  // Ticker represents a price feed for a given asset pair i.e. BTC/USD. The
  // price feed is scaled to a number of decimal places and has a minimum number
  // of providers required to consider the ticker valid.
  Ticker ticker = 1 [ (gogoproto.nullable) = false ];

  // ProviderConfigs is the list of provider-specific configs for this Market.
  repeated ProviderConfig provider_configs = 2 [ (gogoproto.nullable) = false ];
}

// Ticker represents a price feed for a given asset pair i.e. BTC/USD. The price
// feed is scaled to a number of decimal places and has a minimum number of
// providers required to consider the ticker valid.
message Ticker {
  option (gogoproto.goproto_stringer) = false;
  option (gogoproto.stringer) = false;

  // CurrencyPair is the currency pair for this ticker.
  connect.types.v2.CurrencyPair currency_pair = 1 [ (gogoproto.nullable) = false ];

  // Decimals is the number of decimal places for the ticker. The number of
  // decimal places is used to convert the price to a human-readable format.
  uint64 decimals = 2;

  // MinProviderCount is the minimum number of providers required to consider
  // the ticker valid.
  uint64 min_provider_count = 3;

  // Enabled is the flag that denotes if the Ticker is enabled for price
  // fetching by an oracle.
  bool enabled = 14;

  // MetadataJSON is a string of JSON that encodes any extra configuration
  // for the given ticker.
  string metadata_JSON = 15;
}

message ProviderConfig {
  // Name corresponds to the name of the provider for which the configuration is
  // being set.
  string name = 1;

  // OffChainTicker is the off-chain representation of the ticker i.e. BTC/USD.
  // The off-chain ticker is unique to a given provider and is used to fetch the
  // price of the ticker from the provider.
  string off_chain_ticker = 2;

  // NormalizeByPair is the currency pair for this ticker to be normalized by.
  // For example, if the desired Ticker is BTC/USD, this market could be reached
  // using: OffChainTicker = BTC/USDT NormalizeByPair = USDT/USD This field is
  // optional and nullable.
  connect.types.v2.CurrencyPair normalize_by_pair = 3;

  // Invert is a boolean indicating if the BASE and QUOTE of the market should
  // be inverted. i.e. BASE -> QUOTE, QUOTE -> BASE
  bool invert = 4;

  // MetadataJSON is a string of JSON that encodes any extra configuration
  // for the given provider config.
  string metadata_JSON = 15;
}

// MarketMap maps ticker strings to their Markets.
message MarketMap {
  option (gogoproto.goproto_stringer) = false;
  option (gogoproto.stringer) = false;

  // Markets is the full list of tickers and their associated configurations
  // to be stored on-chain.
  map<string, Market> markets = 1 [ (gogoproto.nullable) = false ];
}

```

The `MarketMap` message itself is not stored in state.  Rather, ticker strings are used as key prefixes
so that the data can be stored in a map-like structure, while retaining determinism.

### Params

The `x/marketmap` module stores its params in the keeper state.  The params can be updated with governance or the
keeper authority address.

The `x/marketmap` module contains the following parameters:

| Key               | Type     | Example                                          |
| MarketAuthorities | []string | "cosmos1vq93x443c0fznuf6...q4jd28ke6r46p999s0" |

#### MarketAuthority

A MarketAuthority is the bech32 address that is permitted to submit market updates to the chain.

## Events

The marketmap module emits the following events:

### CreateMarket

| Attribute Key      | Attribute Value |
|--------------------|-----------------|
| currency_pair      | {CurrencyPair}  |
| decimals           | {uint64}        |
| min_provider_count | {uint64}        |
| metadata           | {json string}   |

## Hooks

Other modules can register routines to execute after a certain event has occurred in `x/marketmap`.
The following hooks can be registered:

### AfterMarketCreated

* `AfterMarketCreated(ctx sdk.Context, ticker marketmaptypes.Market) error`
    * Called after a new market is created in `CreateMarket` message server.

### AfterMarketUpdated

* `AfterMarketUpdated(ctx sdk.Context, ticker marketmaptypes.Market) error`
    * Called after a new market is updated in `UpdateMarket` message server.

### AfterMarketGenesis

* `AfterMarketGenesis(ctx sdk.Context, tickers map[string]marketmaptypes.Market) error`
    * Called at the end of `InitGenesis` for the `x/marketmap` keeper.

## Client

### gRPC

A user can query the `marketmap` module using gRPC endpoints.

#### MarketMap

The `MarketMap` endpoint queries the full state of the market map as well as associated information such as
`LastUpdated`.

Example:

```shell
grpcurl -plaintext localhost:9090 connect.marketmap.v2.Query/MarketMap
```

Example response:

```json
{
  "marketMap": {
    "tickers": {
      "BITCOIN/USD": {
        "currencyPair": {
          "Base": "BITCOIN",
          "Quote": "USD"
        },
        "decimals": "8",
        "minProviderCount": "3"
      }
    },
    "paths": {
      "BITCOIN/USD": {
        "paths": [
          {
            "operations": [
              {
                "currencyPair": {
                  "Base": "BITCOIN",
                  "Quote": "USD"
                }
              }
            ]
          }
        ]
      }
    },
    "providers": {
      "BITCOIN/USD": {
        "providers": [
          {
            "name": "kucoin",
            "offChainTicker": "btc_usd"
          },
          {
            "name": "mexc",
            "offChainTicker": "btc-usd"
          },
          {
            "name": "binance",
            "offChainTicker": "BTCUSD"
          }
        ]
      }
    }
  },
  "lastUpdated": "1"
}
```

#### LastUpdated

The `LastUpdated` endpoint queries the last block height that the market map was updated.
This can be consumed by oracle service providers to recognize when their local configurations
must be updated using the heavier `MarketMap` query.

Example:

```shell
grpcurl -plaintext localhost:9090 connect.marketmap.v2.Query/LastUpdated
```

Example response:

```json
{
  "lastUpdated": "1"
}
```

#### Params

The params query allows users to query values set as marketmap parameters.

Example:

```shell
grpcurl -plaintext localhost:9090 connect.marketmap.v2.Query/Params
```

Example response:

```json
{
  "params": {
    "marketAuthorities": "[cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn]"
  }
}
```

### CLI

A user can query the `marketmap` module using the CLI.

#### MarketMap

The `MarketMap` endpoint queries the full state of the market map as well as associated information such as
`LastUpdated` and `Version`.

Example:

```shell
  connectd q marketmap market-map
```

#### LastUpdated

The `LastUpdated` query queries the last block height that the market map was updated.
This can be consumed by oracle service providers to recognize when their local configurations
must be updated using the heavier `MarketMap` query.

Example:

```shell
  connectd q marketmap last-updated
```

#### Params

The params query allows users to query values set as marketmap parameters.

Example:

```shell
  connectd q marketmap params
```
