# KuCoin Provider

## Overview

The KuCoin provider is utilized to fetch pricing data from the KuCoin websocket API. You need to apply for one of the two tokens below to create a websocket connection. It should be noted that: if you subscribe to spot/margin data, you need to obtain tokens through the spot base URL; if you subscribe to futures data, you need to obtain tokens through the futures base URL, which cannot be mixed. **Data is pushed every 100ms.** Note that the KuCoin provider requires a custom websocket connection handler to be used, as the WSS is dynamically generated at start up. 

This implementation subscribes to the spot markets by default, but support for future and orderbook data is also available.

To determine all supported markets, you can use the [get all tickers](https://docs.kucoin.com/#get-all-tickers) endpoint.

```bash
curl https://api.kucoin.com/api/v1/market/allTickers
```
