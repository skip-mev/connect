# Servers

## Overview

This directory contains all of the servers that are supported for the general purpose oracle as well as an metrics instrumentation server that can be used to expose metrics to Prometheus. Each non-metrics server is responsible for running a GRPC server that internally exposes data from a general purpose oracle.


## Servers

* **[Price Oracle Server](./oracle/)** - This server is responsible for running a GRPC server that exposes price data from a price oracle. This server is meant to be run alongside a Cosmos SDK application that is utilizing the general purpose oracle module. However, it can also be run as a standalone server that exposes price data to any application that is able to connect to a GRPC server.
* **[Metrics Server](./metrics/)** - This server is responsible for running a GRPC server that exposes metrics that can be scraped by Prometheus. Specifically, this server can expose metrics data from the price oracle server as well as the price oracle client.

## Usage

To enable any of these servers, please read over the [oracle configurations](../../oracle/config/README.md) documentation.
