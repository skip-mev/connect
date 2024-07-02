# Crypto.com Provider

## Overview

The Crypto.com provider is used to fetch the ticker price from the [Crypto.com websocket API](https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#ticker-instrument_name). The websocket is [rate limited](https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#rate-limits) with a maximum of 100 requests per second. This provider does not require any API keys. To determine the acceptable set of base and quote currencies, you can reference the [get instruments API](https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#reference-and-market-data-api).

To better distribute system load, a single market data websocket connection is limited to a maximum of 400 subscriptions. Once this limit is reached, further subscription requests will be rejected with the EXCEED_MAX_SUBSCRIPTIONS error code. A user should establish multiple connections if additional market data subscriptions are required. The names of the markets (BTCUSD-PERP vs. BTC_USD) represent the perpetual vs. spot markets.
