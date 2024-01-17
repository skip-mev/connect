# API (HTTP) Providers

## Overview

API providers utilize rest APIs to retrieve data from external sources. The data is then transformed into a common format and aggregated across multiple providers. To implement a new provider, please read over the base provider documentation in [`providers/base/README.md`](../base/README.md).

## Supported Providers

The current set of supported providers are:

* [CoinGecko](./coingecko/README.md) - CoinGecko is a cryptocurrency data aggregator that provides a free API for fetching cryptocurrency data. CoinGecko is a **secondary data source** for the oracle.
* [Coinbase](./coinbase/README.md) - Coinbase is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Coinbase is a **primary data source** for the oracle.
* 