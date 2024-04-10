# Factories

## Overview

Factories are used to create an underlying set of data providers that will be utilized by the oracle sidecar. Currently, the factory is primarily built to support price feeds, but later will be extended to support other data types.

## Supported Provider Factories

* **Price Feed Factory**: This factory is used to construct a set of API and Websocket oracle price feed providers that fetch price data from various sources.
* **Market Map Factory**: This factory is used to construct a set of API oracle market providers that fetch market data from market map providers - providers that are responsible for determining the markets the oracle should be fetching prices for.
