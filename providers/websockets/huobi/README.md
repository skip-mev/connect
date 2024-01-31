# OKX Provider

## Overview

The Huobi provider is used to fetch the ticker price from the [Huobi websocket API](https://huobiapi.github.io/docs/spot/v1/en/#introduction-10). All data of websocket Market APIs are compressed with GZIP and need to be unzipped.

The server will send a ping message and expect a pong sent back promptly.  If a pong is not sent back after 2 pings, the connection will be disconnected.

Huobi provides [public channels](https://huobiapi.github.io/docs/spot/v1/en/#introduction-10).

* Public channels -- No authentication is required, include tickers channel, K-Line channel, limit price channel, order book channel, and mark price channel etc.

The exact channel that is used to subscribe to the ticker price is the [`Market Tickers Topic`](https://huobiapi.github.io/docs/spot/v1/en/#market-ticker). This pushes data every 100ms.

To retrieve all supported [pais](https://huobiapi.github.io/docs/spot/v1/en/#get-latest-tickers-for-all-pairs), please run the following command

```bash
 curl "https://api.huobi.pro/market/tickers"   
```
