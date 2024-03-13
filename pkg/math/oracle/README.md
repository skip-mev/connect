# Oracle Price Aggregator

## Overview

The aggregator defined in this package is used to convert a set of price feeds into a common set of prices. The prices used in aggregation contain 

* the most recent prices fetched by each provider.
* `index` prices which are the previous aggregated prices.

Each of the prices provided are safe to use as they are already filtered by the oracle - specifically each price is within the `MaxPriceAge` window.

## Configuration

### MarketMap

The MarketMap contains:

* `Tickers` - the list of tickers that we care about.
* `Providers` - each ticker maps to a list of providers that fetch prices for that ticker.
* `Paths` - each ticker maps to a set of paths that are used to calculate the price of that ticker. Each path is a set of operations that define how to convert a given provider price to the desired ticker price. These are either direct conversions or indirect conversions that use the `index` price.

> Say that the oracle is configured to fetch prices for the following feeds:
>
> * BITCOIN/USD (8 decimal precision)
> * BITCOIN/USDT (8 decimal precision)
> * ETHEREUM/USD (18 decimal precision)
> * ETHEREUM/USDT (18 decimal precision)
> * USDT/USD (6 decimal precision)
> * USDC/USD (6 decimal precision)
> 
> We care about calculating the price of 
> 
> * BITCOIN/USD (8 decimal precision)
> * ETHEREUM/USD (8 decimal precision)
> * USDT/USD (6 decimal precision)

Using the example above, a sample market map might look like the following (assume we have correctly configured the providers):

```golang
marketmap = mmtypes.MarketMap{
    Tickers: map[string]mmtypes.Ticker{
        BTC_USD.String():   BTC_USD,
        BTC_USDT.String():  BTC_USDT,
        USDT_USD.String():  USDT_USD,
        USDC_USDT.String(): USDC_USDT,
        ETH_USD.String():   ETH_USD,
        ETH_USDT.String():  ETH_USDT,
    },
    Paths: map[string]mmtypes.Paths{
        BTC_USD.String(): {
            Paths: []mmtypes.Path{
                {
                    // COINBASE BTC/USD = BTC/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: BTC_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     coinbase.Name,
                        },
                    },
                },
                {
                    // COINBASE BTC/USDT * INDEX USDT/USD = BTC/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: BTC_USDT.CurrencyPair,
                            Invert:       false,
                            Provider:     coinbase.Name,
                        },
                        {
                            CurrencyPair: USDT_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     oracle.IndexPrice,
                        },
                    },
                },
                {
                    // BINANCE BTC/USDT * INDEX USDT/USD = BTC/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: BTC_USDT.CurrencyPair,
                            Invert:       false,
                            Provider:     binance.Name,
                        },
                        {
                            CurrencyPair: USDT_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     oracle.IndexPrice,
                        },
                    },
                },
            },
        },
        ETH_USD.String(): {
            Paths: []mmtypes.Path{
                {
                    // COINBASE ETH/USD = ETH/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: ETH_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     coinbase.Name,
                        },
                    },
                },
                {
                    // COINBASE ETH/USDT * INDEX USDT/USD = ETH/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: ETH_USDT.CurrencyPair,
                            Invert:       false,
                            Provider:     coinbase.Name,
                        },
                        {
                            CurrencyPair: USDT_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     oracle.IndexPrice,
                        },
                    },
                },
                {
                    // BINANCE ETH/USDT * INDEX USDT/USD = ETH/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: ETH_USDT.CurrencyPair,
                            Invert:       false,
                            Provider:     binance.Name,
                        },
                        {
                            CurrencyPair: USDT_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     oracle.IndexPrice,
                        },
                    },
                },
            },
        },
        USDT_USD.String(): {
            Paths: []mmtypes.Path{
                {
                    // COINBASE USDT/USD = USDT/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: USDT_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     coinbase.Name,
                        },
                    },
                },
                {
                    // COINBASE USDC/USDT ^ -1 = USDT/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: USDC_USDT.CurrencyPair,
                            Invert:       true,
                            Provider:     coinbase.Name,
                        },
                    },
                },
                {
                    // BINANCE USDT/USD = USDT/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: USDT_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     binance.Name,
                        },
                    },
                },

                {
                    // KUCOIN BTC/USDT ^-1 * INDEX BTC/USD = USDT/USD
                    Operations: []mmtypes.Operation{
                        {
                            CurrencyPair: BTC_USDT.CurrencyPair,
                            Invert:       true,
                            Provider:     kucoin.Name,
                        },
                        {
                            CurrencyPair: BTC_USD.CurrencyPair,
                            Invert:       false,
                            Provider:     oracle.IndexPrice,
                        },
                    },
                },
            },
        },
    },
}
```

From the above market map, we can see that the following paths are used to calculate the price of each ticker:

* BTC/USD
    * COINBASE BTC/USD
    * COINBASE BTC/USDT * INDEX USDT/USD
    * BINANCE BTC/USDT * INDEX USDT/USD
* ETH/USD
    * COINBASE ETH/USD
    * COINBASE ETH/USDT * INDEX USDT/USD
    * BINANCE ETH/USDT * INDEX USDT/USD
* USDT/USD
    * COINBASE USDT/USD
    * COINBASE USDC/USDT ^ -1
    * BINANCE USDT/USD
    * KUCOIN BTC/USDT ^-1 * INDEX BTC/USD

A few important considerations:

1. Each ticker (BTC/USD, ETH/USD, USDT/USD) can have a configured `MinimumProviderCount` which is the minimum number of providers that are required to calculate the price of the ticker.
2. Each path that is not a direct conversion (e.g. BTC/USD) must configure the second operation to utilize the `index` price i.e. [`IndexPrice`](./median.go).

## Aggregation

### Precision

As a given path is being traversed, we convert each price to a common decimal precision - [`ScaledDecimals`](./math.go) - in order to correctly multiply & invert. We choose 36 decimal precision as the common precision - primarily to retain the maximum precision possible. The conversion is done as follows:

* BITCOIN/USDT: 8 -> 36
* BITCOIN/USDC: 8 -> 36
* ETHEREUM/USDT: 18 -> 36
* ETHEREUM/USDC: 18 -> 36
* USDT/USD: 6 -> 36
* USDC/USD: 6 -> 36

### Example Aggregation

Given the market map above, let's assume that we have the following prices fetched by the providers:

* COINBASE BTC/USD: 71_000
* COINBASE BTC/USDT: 70_000
* BINANCE BTC/USDT: 70_500
* INDEX USDT/USD: 1.05

The aggregator will calculate the price of BTC/USD using the following paths:

1. COINBASE BTC/USD
2. COINBASE BTC/USDT * INDEX USDT/USD
3. BINANCE BTC/USDT * INDEX USDT/USD

This gives a set of prices:

* COINBASE BTC/USD: 71_000
* COINBASE BTC/USDT * INDEX USDT/USD: 73_500
* BINANCE BTC/USDT * INDEX USDT/USD: 73_575

The final price of BTC/USD is the median of the above prices, which is 73_500. In the case of an even number of prices, the median is the average of the two middle numbers.

## Other Considerations

### Cycle Detection

It is possible to have cycles in the market map. If the price of a ticker is dependent on a different ticker, which in turn is dependent on the first ticker, then we have a cycle. This can affect price liveness and can cause the oracle to be stuck in a loop. To prevent this, we recommend that markets that are dependent on each other have a sufficient amount of providers, have considerable `MinProviderCount`, and have sufficent amounts of direct conversions (i.e. not dependent on other tickers).

If a cycle does exist, it will ~likely~ be resolved after a few iterations of the oracle.
