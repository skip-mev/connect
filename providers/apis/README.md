# API (HTTP) Providers

## Overview

API providers utilize rest APIs to retrieve data from external sources. The data is then transformed into a common format and aggregated across multiple providers. To implement a new provider, please read over the base provider documentation in [`providers/base/README.md`](../base/README.md).

## Supported Providers

The current set of supported providers are:

> Note: The URLs provided are endpoints that can be used to determine the set of available currency pairs and their respective symbols. The `jq` command is used to format the JSON response for readability. Note that some of these may require a VPN to access. Depending on the provider, the markets supported as well as the URL may differ.

* [Binance](./binance/README.md) - Binance is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Binance is a **primary data source** for the oracle.
    * Check all supported markets: 
        * `curl https://api.binance.us/api/v3/ticker/price | jq`
        * `curl https://api.binance.com/api/v3/exchangeInfo | jq`
    * Check if a given market is supported:
        * `curl https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT | jq`
* [Bitstamp](./bitstamp/README.md) - Bitstamp is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Bitstamp is a **primary data source** for the oracle.
    * Check all supported markets: 
        * `curl https://www.bitstamp.net/api/v2/trading-pairs-info/ | jq`
    * Check if a given market is supported: 
        * `curl https://www.bitstamp.net/api/v2/ticker/{btcusd} | jq`
* [Coinbase](./coinbase/README.md) - Coinbase is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Coinbase is a **primary data source** for the oracle.
    * Check all supported markets: 
        * `curl https://api.exchange.coinbase.com/currencies | jq`
        * `curl https://api.exchange.coinbase.com/products | jq`
    * Check if a given market is supported: 
        * `curl https://api.coinbase.com/v2/prices/{DYDX-USDC}/spot | jq`
* [CoinGecko](./coingecko/README.md) - CoinGecko is a cryptocurrency data aggregator that provides a free API for fetching cryptocurrency data. CoinGecko is a **secondary data source** for the oracle. This is not recommended for use in production.
    * Check all supported markets: 
        * `curl https://api.coingecko.com/api/v3/coins/list | jq`
    * Check if a given market is supported: 
        * `curl https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd | jq`
* [dYdX](./dydx/README.md) - dYdX is a decentralized exchange built using the Cosmos SDK. dYdX is a market map provider - we use it to fetch the list of markets the side-car should fetch prices for.
* [GeckoTerminal](./geckoterminal/README.md) - GeckoTerminal is price provider that aggregates prices of tokens on a variety of blockchains, pools,  and decentralized exchanges. To fetch the price of a token, you need to provide the token's address. 
* [Kraken](./kraken/README.md) - Kraken is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Kraken is a **primary data source** for the oracle.
    * Check all supported markets: 
        * `curl https://api.kraken.com/0/public/AssetPairs | jq`
    * Check if a given market is supported: 
        * `curl https://api.kraken.com/0/public/Ticker?pair=ETHUSD | jq`
* [Raydium](./defi/raydium/price_fetcher.go) - Raydium is a decentralized exchange on the Solana blockchain. Raydium is a **primary data source** for the oracle.
* [Uniswap V3](./defi/uniswapv3/README.md) - Uniswap V3 is a decentralized exchange on the Ethereum blockchain. Uniswap V3 is a **primary data source** for the oracle.
