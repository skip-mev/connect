# slinky

<!-- markdownlint-disable MD013 -->
<!-- markdownlint-disable MD041 -->
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#wip)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://godoc.org/github.com/skip-mev/slinky)
[![Go Report Card](https://goreportcard.com/badge/github.com/skip-mev/slinky?style=flat-square)](https://goreportcard.com/report/github.com/skip-mev/slinky)
[![Version](https://img.shields.io/github/tag/skip-mev/slinky.svg?style=flat-square)](https://github.com/skip-mev/slinky/releases/latest)
[![License: Apache-2.0](https://img.shields.io/github/license/skip-mev/slinky.svg?style=flat-square)](https://github.com/skip-mev/slinky/blob/main/LICENSE)
[![Lines Of Code](https://img.shields.io/tokei/lines/github/skip-mev/slinky?style=flat-square)](https://github.com/skip-mev/slinky)

A general purpose price oracle leveraging ABCI++.

## Install

```shell
$ go install github.com/skip-mev/slinky
```

## Overview

The slinky repository is composed of the following core packages:

* **abci** - This package contains the [vote extension](./abci/ve/README.md), [proposal](./abci/proposals/README.md), and [preblock handlers](./abci/preblock/oracle/README.md) that are used to broadcast oracle data to the network and to store it in the blockchain.
* **oracle** - This [package](./oracle/) contains the main oracle that aggregates external data sources before broadcasting it to the network. You can reference the provider documentation [here](./providers/base/README.md) to get a high level overview of how the oracle works.
* **providers** - This package contains a collection of [websocket](./providers/websockets/README.md) and [API](./providers/apis/README.md) based data providers that are used by the oracle to collect external data. 
* **x/oracle** - This package contains a Cosmos SDK module that allows you to store oracle data on a blockchain.
* **x/alerts** - This package contains a Cosmos SDK module that allows network participants to create alerts when oracle data that is in violation of some condition is broadcast to the network and stored on the blockchain.
* **x/sla** - This package contains a Cosmos SDK module that allows you to create service level agreements (SLAs) that can be used to incentivize network participants to consistently, reliably provide data with high uptime.
* **x/marketmap** - This [package](./x/marketmap/README.md) contains  a Cosmos SDK module that allows for market configuration to be stored and updated on a blockchain.

## Usage

To run the oracle, run the following command.

```bash
$ make start-oracle
```

This will:

1. Start a blockchain with a single validator node. It may take a few minutes to build and reach a point where vote extensions can be submitted.
2. Start the oracle side-car that will aggregate prices from external data providers and broadcast them to the network. To check the current aggregated prices on the side-car, you can run `curl localhost:8080/slinky/oracle/v1/prices`.
3. Host a prometheus instance that will scrape metrics from the oracle side-car. Navigate to http://localhost:9090 to see all network traffic and metrics pertaining to the oracle sidecar. Navigate to http://localhost:8001 to see all application-side oracle metrics.
4. Host a profiler that will allow you to profile the oracle side-car. Navigate to http://localhost:6060 to see the profiler.

After a few minutes, run the following commands to see the prices written to the blockchain:

```bash
# access the blockchain container
$ docker exec -it slinky-blockchain-1 bash

# query the price of bitcoin in USD on the node
$ (slinky-blockchain-1) ./build/slinkyd q oracle price BITCOIN USD
```

Result: 

```bash
decimals: "8"
id: "0"
nonce: "44"
price:
  block_height: "46"
  block_timestamp: "2024-01-29T01:43:48.735542Z"
  price: "4221100000000"
```

To stop the oracle, run the following command:

```bash
$ make stop-oracle
```

## Metrics

### Oracle Service Metrics

* metrics relevant to the oracle service's health + operation are [here](./oracle/metrics/README.md)
* metrics relevant to the operation / health of the oracle's providers are [here](./providers/base/metrics/README.md)

### Oracle Application / Network Metrics

* metrics relevant to the network's (that is running the instance slinky) performance are [here](./service/metrics/README.md)

## Basic Perfomance Analysis

> **Note: These are numbers based on 14 providers and 9 currency pairs over a 24 hour period.**

* ~**5 milliseconds** between price updates across all providers and price feeds.
* ~**14 million** total price updates.
* ~**60 go routines** are running at any given time.
* ~**7x** improvement in performance of websocket providers over API providers.

To test these numbers yourself, spin up the the oracle server following the instructions above and then navigate to http://localhost:9090. From here, you can run the prometheus queries defined in the packages below to get insight into the oracle's performance.

* [Oracle Graphs & Queries](./oracle/metrics/README.md#usage): Provides insight into the oracle's performance by provider, price feed, and currency pair. All includes nice visualizations of the oracle's aggregated prices and the individual prices that are aggregated to produce the oracle's aggregated prices.
* [Data Provider Queries](./providers/base/metrics/README.md#usage): Provides general insight into how often price feeds are updated by status (success/failure), provider (binance, coinbase, etc.), price feed (BTC/USD, ETH/USD), and provider type (api/websocket).
* [Websocket Performance Queries](./providers/base/websocket/metrics/README.md#usage): Provides insight into how often websocket providers are successfully updating their data. This is a combination of metrics related to the underlying connection as well as the data handler which is responsible for processing the data received from the Websocket connection.
* [API Performance Queries](./providers/base/api/metrics/README.md#usage): Provides insight into how often API providers are successfully updating their data.

## DEFI

- Defi providers are currently an experimental feature, to generate local configs to run the defi providers locally, run
```bash
$ DEFI_PROVIDERS_ENABLED=true SOLANA_NODE_ENDPOINT=<SOLANA_NODE>  make update-local-configs
```
Then, to run
```bash
$ make run-oracle-server
```
## Future Work

The oracle side car is a combination of the oracle and provider packages. This is being moved to a [separate repository](https://github.com/skip-mev/slinky-sidecar).
