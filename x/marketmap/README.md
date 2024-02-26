# `x/marketmap`

## Contents

* [Concepts](#concepts)
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
    * [REST](#rest)

## Concepts

The `x/marketmap` module encapsulates a system for creating and updating a unified configuration that is stored on-chain
and consumed by a set of oracle service providers (Slinky oracle, etc.).

The core goal of the system is to collect off-chain market updates and to post them on chain, informing oracle service
providers to fetch prices for new markets.

The data is stored in a `MarketMap` data structure which can be queried and consumed by oracle services.

## State

### MarketMap

The market map data is as follows:

```protobuf
// Ticker represents a price feed for a given asset pair i.e. BTC/USD. The price
// feed is scaled to a number of decimal places and has a minimum number of
// providers required to consider the ticker valid.
message Ticker {
  option (gogoproto.goproto_stringer) = false;
  option (gogoproto.stringer) = false;

  // CurrencyPair is the currency pair for this ticker.
  slinky.types.v1.CurrencyPair currency_pair = 1
      [ (gogoproto.nullable) = false ];

  // Decimals is the number of decimal places for the ticker. The number of
  // decimal places is used to convert the price to a human-readable format.
  uint64 decimals = 3;
  // MinProviderCount is the minimum number of providers required to consider
  // the ticker valid.
  uint64 min_provider_count = 4;

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
}

// Path is the list of convertable markets that will be used to convert the
// prices of a set of tickers to a common ticker.
message Path {
  // Operations is an ordered list of operations that will be taken. These must
  // be topologically sorted to ensure that the conversion is possible i.e. DAG.
  repeated Operation operations = 1 [ (gogoproto.nullable) = false ];
}

// Operation represents the operation configuration for a given ticker.
message Operation {
  // CurrencyPair is the on-chain currency pair for this ticker.
  slinky.types.v1.CurrencyPair currency_pair = 1
      [ (gogoproto.nullable) = false ];

  // Invert is a boolean that indicates whether the price of the ticker should
  // be inverted.
  bool invert = 2;
}

message Paths {
  // Paths is the list of convertable markets that will be used to convert the
  // prices of a set of tickers to a common ticker.
  repeated Path paths = 1 [ (gogoproto.nullable) = false ];
}

message Providers {
  // Providers is the list of provider configurations for the given ticker.
  repeated ProviderConfig providers = 1 [ (gogoproto.nullable) = false ];
}

message MarketMap {
  option (gogoproto.goproto_stringer) = false;
  option (gogoproto.stringer) = false;

  // Tickers is the full list of tickers and their associated configurations
  // to be stored on-chain.
  map<string, Ticker> tickers = 1 [ (gogoproto.nullable) = false ];

  // Paths is a map from CurrencyPair to all paths that resolve to that pair
  map<string, Paths> paths = 2 [ (gogoproto.nullable) = false ];

  // Providers is a map from CurrencyPair to each of to provider-specific
  // configs associated with it.
  map<string, Providers> providers = 3 [ (gogoproto.nullable) = false ];
}
```

The `MarketMap` message itself is not stored in state.  Rather, ticker strings are used as key prefixes
so that the data can be stored in a map-like structure, while retaining determinism.

### Params

The `x/marketmap` module stores its params in the keeper state.  The params can be updated with governance or the
keeper authority address.

The `x/marketmap` module contains the following parameters:

| Key               | Type     | Example                                          |
| MarketAuthority | string | "cosmos1vq93x443c0fznuf6...q4jd28ke6r46p999s0" |
| Version         | uint64 | 20                                             |

#### MarketAuthority

The MarketAuthority is the bech32 address that is permitted to submit market updates to the chain.

#### Version

Version is the version of the MarketMap schema. This version is returned in the `GetMarketMap` query and can be used
by oracle service providers to verify the schema they are consuming.  When being modified via governance, the new value
must always be greater than the current value.

## Events

The marketmap module emits the following events:

### CreateMarket

| Attribute Key      | Attribute Value |
|--------------------|-----------------|
| currency_pair      | {CurrencyPair}  |
| decimals           | {uint64}        |
| min_provider_count | {uint64}        |
| metadata           | {json string}   |
| providers          | {[]Provider}    |
| paths              | {[]Path]}       |

## Hooks

Other modules can register routines to execute after a certain event has occurred in `x/marketmap`.
The following hooks can be registered:

### AfterMarketCreated

* `AfterMarketCreated(ctx sdk.Context, ticker marketmaptypes.Ticker) error`
    * Called after a new market is created in `CreateMarket` message server.

### AfterMarketUpdated

* `AfterMarketUpdated(ctx sdk.Context, ticker marketmaptypes.Ticker) error`
    * Called after a new market is updated in `UpdateMarket` message server.
TODO BLO-866

## Client

### CLI

TODO BLO-920

### gRPC

A user can query the `marketmap` module using gRPC endpoints.

#### MarketMap

The `MarketMap` endpoint queries the full state of the market map as well as associated information such as
`LastUpdated` and `Version`.

Example:

```shell
grpcurl -plaintext localhost:9090 slinky.marketmap.v1.Query/MarketMap
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
grpcurl -plaintext localhost:9090 slinky.marketmap.v1.Query/LastUpdated
```

Example response:

```json
{
  "lastUpdated": "1"
}
```

#### Params

The params command allows users to query values set as marketmap parameters.

Example:

```shell
grpcurl -plaintext localhost:9090 slinky.marketmap.v1.Query/Params
```

Example response:

```json
{
  "params": {
    "marketAuthority": "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
  }
}
```
