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

## Validator Usage

To read how to run the oracle as a validator based on the chain, please reference the [validator documentation](https://docs.skip.money/slinky/integrations).

## Developer Usage

To run the oracle, run the following command.

```bash
$ make start-all-core-dev
```

This will:

1. Start a blockchain with a single validator node. It may take a few minutes to build and reach a point where vote extensions can be submitted.
2. Start the oracle side-car that will aggregate prices from external data providers and broadcast them to the network. To check the current aggregated prices on the side-car, you can run `curl localhost:8080/slinky/oracle/v1/prices`.
3. Host a prometheus instance that will scrape metrics from the oracle side-car. Navigate to http://localhost:9091 to see all network traffic and metrics pertaining to the oracle sidecar. Navigate to http://localhost:8002 to see all application-side oracle metrics.
4. Host a profiler that will allow you to profile the oracle side-car. Navigate to http://localhost:6060 to see the profiler.
5. Host a grafana instance that will allow you to visualize the metrics scraped by prometheus. Navigate to http://localhost:3000 to see the grafana dashboard. The default username and password are `admin` and `admin`, respectively.

After a few minutes, run the following commands to see the prices written to the blockchain:

```bash
# access the blockchain container
$ docker exec -it slinky-blockchain-1 bash

# query the price of bitcoin in USD on the node
$ (slinky-blockchain-1) ./build/slinkyd q oracle price BTC USD
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
$ make stop-all-dev
```

## Metrics

### Oracle Service Metrics

* metrics relevant to the oracle service's health + operation are [here](./metrics.md)

### Oracle Application / Network Metrics

* metrics relevant to the network's (that is running the instance slinky) performance are [here](./service/metrics/README.md)

