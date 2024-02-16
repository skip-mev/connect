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

TODO need to finalize data structure

### Params

The `x/marketmap` module stores its params in the keeper state.  The params can be updated with governance or the
keeper authority address.

The `x/marketmap` module contains the following parameters:

| Key             | Type   | Example                                        |
|-----------------|--------|------------------------------------------------|
| MarketAuthority | string | "cosmos1vq93x443c0fznuf6...q4jd28ke6r46p999s0" |
| Version         | uint64 | 20                                             |

#### MarketAuthority

The MarketAuthority is the bech32 address that is permitted to submit market updates to the chain.

#### Version

Version is the version of the MarketMap schema. This version is returned in the `GetMarketMap` query and can be used
by oracle service providers to verify the schema they are consuming.  When being modified via governance, the new value
must always be greater than the current value.

## Events

TODO BLO-921

## Hooks

Other modules can register routines to execute after a certain event has occurred in `x/marketmap`.
The following hooks can be registered:

### AfterMarketCreated

* `AfterMarketCreated(ctx sdk.Context, TODO) error`
  * Called after a new market is created in `CreateMarket` message server.

### AfterMarketUpdated

* `AfterMarketUpdated(ctx sdk.Context, TODO) error`
  * Called after a new market is updated in `UpdateMarket` message server.

## Client

### CLI

TODO BLO-920

### gRPC

TODO BLO-919

### Rest

TODO BLO-919
