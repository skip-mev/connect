# Oracle Price Conversion Aggregation

## Overview

The aggregation function defined in this package is used to convert a set of price feeds into a common set of prices. This is done by converting all prices to a common decimal precision and then aggregating the prices based on the price feed conversions defined for each price feed.

## Conversion

Say that the oracle is configured to fetch prices for the following feeds:

* BITCOIN/USDT (8 decimal precision)
* BITCOIN/USDC (8 decimal precision)
* ETHEREUM/USDT (18 decimal precision)
* ETHEREUM/USDC (18 decimal precision)
* USDT/USD (6 decimal precision)
* USDC/USD (6 decimal precision)

We care about calculating the price of BITCOIN and ETHEREUM in USD. To do this, we need to convert all prices to a common decimal precision - [`ScaledDecimals`](./utils.go). We choose 36 decimal precision as the common precision - primarily to retain the maximum precision possible. The conversion is done as follows:

* BITCOIN/USDT: 8 -> 36
* BITCOIN/USDC: 8 -> 36
* ETHEREUM/USDT: 18 -> 36
* ETHEREUM/USDC: 18 -> 36
* USDT/USD: 6 -> 36
* USDC/USD: 6 -> 36

## Aggregation

The main oracle contains a [list of valid price conversions per desired price feed](./../../../oracle/config/README.md#aggregate-market-configurations). For example, to calculate the price of BITCOIN in USD, we need to convert the price of BITCOIN/USDT to USD, and the price of BITCOIN/USDC to USD. If the list contains multiple valid conversions, the aggregation function will return the median of the prices - where an average is taken if the number of prices is even.

Following the example above, the aggregation function will return the median of the following prices:

* BITCOIN/USDT * USDT/USD -> BITCOIN/USD
* BITCOIN/USDC * USDC/USD -> BITCOIN/USD
* ETHEREUM/USDT * USDT/USD -> ETHEREUM/USD
* ETHEREUM/USDC * USDC/USD -> ETHEREUM/USD

It will then return the median of the two prices to get the final price of BITCOIN in USD.
