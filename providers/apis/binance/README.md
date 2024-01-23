# Binance Provider

## Overview

The Binance provider is used to fetch the spot price for cryptocurrencies from the [Binance API](https://binance-docs.github.io/apidocs/spot/en/#general-info).

## Supported Pairs

To determine the pairs (in the form `BASEQUOTE`) currencies that the Binance provider supports, you can run the following command:

```bash
$ curl -X GET https://api.binance.vision/api/v3/ticker/price         
```
