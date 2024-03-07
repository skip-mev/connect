# CoinGecko Provider

## Overview

The CoinGecko provider is used to fetch the spot price for cryptocurrencies from the [CoinGecko API](https://www.coingecko.com/en/api). This provider can be configured to fetch with or without an API key. Note that without an API key, it is very likely that the CoinGecko API will rate limit your requests. The CoinGecko API fetches an aggregated TWAP price any given currency pair.

## Supported Bases

To determine the base currencies that the CoinGecko provider supports, you can run the following command:

```bash
$ curl -X GET https://api.coingecko.com/api/v3/coins/list
```

## Supported Quotes

To determine the quote currencies that the CoinGecko provider supports, you can run the following command:

```bash
$ curl -X GET https://api.coingecko.com/api/v3/simple/supported_vs_currencies
```
