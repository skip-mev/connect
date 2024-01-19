# Web Socket Providers

## Overview

Web socket providers utilize web socket APIs / clients to retrieve data from external sources. The data is then transformed into a common format and aggregated across multiple providers. To implement a new provider, please read over the base provider documentation in [`providers/base/README.md`](../base/README.md).

Web sockets are preferred over HTTP APIs for real-time data as they only require a single connection to the server, whereas HTTP APIs require a new connection for each request. This makes web sockets more efficient for real-time data. Additionally, web sockets typically have lower latency than HTTP APIs, which is important for real-time data.

## Supported Providers

The current set of supported providers are:

* [Crypto.com](./cryptodotcom/README.md) - Crypto.com is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Crypto.com is a **primary data source** for the oracle.   
* [OKX](./okx/README.md) - OKX is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. OKX is a **primary data source** for the oracle.
* [ByBit](./bybit/README.md) - ByBit is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. ByBit is a **primary data source** for the oracle.
