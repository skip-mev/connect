# Index Price Aggregator

## Overview

The price aggregator defined in this package is used to convert a set of price feeds into a common set of prices. The prices used in aggregation contain:

* the most recent prices fetched by each provider.
* `index` prices which are the previous aggregated prices.
* `scaled` prices which are effectively the index price but with a common decimal precision (as determined by the market map).

Each of the prices provided are safe to use as they are already filtered by the oracle - specifically each price is within the `MaxPriceAge` window.

## Configuration

### MarketMap

The MarketMap contains:

* `Tickers` - the list of tickers that we care about.
* `Providers` - each ticker maps to a list of providers that fetch prices for that ticker.

> Say that the oracle is configured to fetch prices for the following feeds:
>
> * BITCOIN/USD (8 decimal precision)
> * ETHEREUM/USD (18 decimal precision)
> * USDT/USD (6 decimal precision)

Using the example above, a sample market map might look like the following (assume we have correctly configured the providers):

```golang
	marketmap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			BTC_USD.String(): {
                Ticker: mmtypes.Ticker{
                    CurrencyPair:     constants.BITCOIN_USD,
                    Decimals:         8,
                    MinProviderCount: 3,
                },
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:           coinbase.Name,
						OffChainTicker: "BTC-USDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
					{
						Name:           binance.Name,
						OffChainTicker: "BTCUSDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
			ETH_USD.String(): {
                Ticker: mmtypes.Ticker{
                    CurrencyPair:     constants.ETHEREUM_USD,
                    Decimals:         11,
                    MinProviderCount: 3,
                },
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "ETH-USD",
					},
					{
						Name:           coinbase.Name,
						OffChainTicker: "ETH-USDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
					{
						Name:           binance.Name,
						OffChainTicker: "ETHUSDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
			USDT_USD.String(): {
                Ticker: mmtypes.Ticker{
                    CurrencyPair:     constants.USDT_USD,
                    Decimals:         6,
                    MinProviderCount: 2,
                },
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "USDT-USD",
					},
					{
						Name:           coinbase.Name,
						OffChainTicker: "USDC-USDT",
						Invert:         true,
					},
					{
						Name:           binance.Name,
						OffChainTicker: "USDTUSD",
					},
					{
						Name:           kucoin.Name,
						OffChainTicker: "BTC-USDT",
						Invert:         true,
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "BTC",
							Quote: "USD",
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
2. Each path that is not a direct conversion (e.g. BTC/USD) must configure the second operation to utilize the `index` price i.e. of a primary ticker i.e. market.

## Aggregation

### Precision

Precision is retained as much as possible in the aggregator. Each price included by each provider is converted to the maximum amount of precision that is possible for the price (and what big.Float is capable of handling). The index prices are always big.Floats with minimal precision lost between conversions, scaling, and aggregation.

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

It is possible to have cycles in the market map. If the price of a ticker is dependent on a different ticker, which in turn is dependent on the first ticker, then we have a cycle. This can affect price liveness and can cause the oracle to be stuck in a loop. To prevent this, we recommend that markets that are dependent on each other have a sufficient amount of providers, have considerable `MinProviderCount`, and have sufficient amounts of direct conversions (i.e. not dependent on other tickers).

If a cycle does exist, it will likely be resolved after a few iterations of the oracle.
