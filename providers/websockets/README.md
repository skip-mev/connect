# Websocket Providers

## Overview

Websocket providers utilize websocket APIs / clients to retrieve data from external sources. The data is then transformed into a common format and aggregated across multiple providers. To implement a new provider, please read over the base provider documentation in [`providers/base/README.md`](../base/README.md).

Websockets are preferred over REST APIs for real-time data as they only require a single connection to the server, whereas HTTP APIs require a new connection for each request. This makes websockets more efficient for real-time data. Additionally, web sockets typically have lower latency than HTTP APIs, which is important for real-time data.

## Supported Providers

The current set of supported providers are:

> Note: The URLs provided are endpoints that can be used to determine the set of available currency pairs and their respective symbols. The `jq` command is used to format the JSON response for readability. Note that some of these may require a VPN to access.

* [Binance](./binance/README.md) - Binance is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Binance is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://api.binance.com/api/v3/exchangeInfo | jq`
    * Check if a given market is supported:
        * `curl https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT | jq`
* [BitFinex](./bitfinex/README.md) - BitFinex is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. BitFinex is a **primary data source** for the oracle.
    * Check all supported markets: 
        * `curl https://api-pub.bitfinex.com/v2/conf/pub:list:currency | jq`
    * Check if a given market is supported: 
        * `curl https://api-pub.bitfinex.com/v2/ticker/t{BTCUSD} | jq`
* [Bitstamp](./bitstamp/README.md) - Bitstamp is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Bitstamp is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://www.bitstamp.net/api/v2/currencies/ | jq`
    * Check if a given market is supported:
        * `curl https://www.bitstamp.net/api/v2/ticker/{btcusd}/ | jq`
* [ByBit](./bybit/README.md) - ByBit is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. ByBit is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://api.bybit.com/v5/market/tickers?category=spot | jq`
        * `curl https://api.bybit.com/v5/market/instruments-info | jq`
    * Check if a given market is supported:
        *  `curl https://api.bybit.com/v5/market/tickers?category=spot&symbol={BTCUSDT} | jq`
* [Coinbase](./coinbase/README.md) - Coinbase is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Coinbase is a **primary data source** for the oracle.
    * Check all supported markets: 
        * `curl https://api.exchange.coinbase.com/currencies | jq`
        * `curl https://api.exchange.coinbase.com/products | jq`
    * Check if a given market is supported: 
        * `curl https://api.coinbase.com/v2/prices/{DYDX-USDC}/spot | jq`
* [Crypto.com](./cryptodotcom/README.md) - Crypto.com is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Crypto.com is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://api.crypto.com/v2/public/get-instruments | jq`
        * `curl https://api.crypto.com/v2/public/get-ticker | jq`
    * Check if a given market is supported:
        * `curl https://api.crypto.com/v2/public/get-ticker?instrument_name={BTCUSD-PERP} | jq`
* [Gate](./gate/README.md) - Gate.io is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Gate.io is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://api.gateio.ws/api/v4/spot/currency_pairs | jq`
    * Check if a given market is supported:
        * `curl https://api.gateio.ws/api/v4/spot/currency_pairs/{ETH_USDT} | jq`
* [Huobi](./huobi/README.md) - Huobi is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Huobi is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://api.huobi.pro/market/tickers | jq`
    * Check if a given market is supported:
        * `curl https://api.huobi.pro/market/trade?symbol=ethusdt |jq`
* [Kraken](./kraken/README.md) - Kraken is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Kraken is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl "https://api.kraken.com/0/public/Assets"` 
        * `curl https://api.kraken.com/0/public/AssetPairs | jq`
        * `curl "https://api.kraken.com/0/public/Ticker" | jq`
    * Check if a given market is supported:
        * `curl https://api.kraken.com/0/public/Ticker?pair={XBTUSD} | jq`
* [KuCoin](./kucoin/README.md) - KuCoin is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. KuCoin is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://api.kucoin.com/api/v1/symbols | jq`
        * `curl https://api.kucoin.com/api/v3/currencies | jq`
        * `curl https://api.kucoin.com/api/v1/market/allTickers | jq`
    * Check if a given market is supported:
        * `curl https://api.kucoin.com/api/v1/market/orderbook/level1?symbol={BTC-USDT} | jq`
* [MEXC](./mexc/README.md) - MEXC is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. MEXC is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://www.mexc.com/open/api/v2/market/ticker | jq`
        * `curl https://www.mexc.com/open/api/v2/market/symbols | jq`
        * `curl https://api.mexc.com/api/v3/exchangeInfo | jq`
    * Check if a given market is supported:
        * `curl https://www.mexc.com/open/api/v2/market/ticker?symbol={BTC_USDT} | jq`
* [OKX](./okx/README.md) - OKX is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. OKX is a **primary data source** for the oracle.
    * Check all supported markets:
        * `curl https://www.okx.com/api/v5/market/index-tickers?quoteCcy={USD} | jq`
        
    * Check if a given market is supported:
        * `curl https://www.okx.com/api/v5/market/index-tickers?instId={BTC-USDT} | jq`
