# dYdX Market Params API

> To learn more about the price aggregation that pairs with this API as well as the structure of the side-car's market config (Market Map), please refer to the [Median Index Price Aggregator readme](../../../pkg/math/oracle/README.md).

The dYdX Market Params API provides a list of all market parameters for the dYdX protocol. Specifically, this includes all markets that the dYdX protocol supports, the exchanges that need to provide prices along with the relevant translations, and a set of operations that convert all markets into a common set of prices. The side-car utilizes this API to update the price providers with all relevant market parameters, such that it can support the dYdX protocol out of the box. 

# Market Map Adapter

In order to utilize the market params API, the dYdX market map provider includes a custom adapter that can translate the market parameters into a format the side-car can understand.

## Ticker Conversions

The ticker conversions are used to convert the specific parameters associated with a given market. Example:

```json
{
    "id": 1000000,
    "pair": "USDT-USD",
    "exponent": -9,
    "min_exchanges": 3,
    "min_price_change_ppm": 1000,
    "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"}]}"
},
```

* `id` - This is currently not applicable to the side car and is not utilized.
* `pair` - This is the pair that the market is associated with. There is a direct conversion here assuming the delimeter is `-` across all markets.
* `exponent` - In the side-car we represent all prices in `Decimals`, which is the absolute value of the `exponent`.
* `min_exchanges` - This is the minimum number of exchanges that need to provide prices for the market to be considered valid. This is a direct mapping to `MinProviderCount` in the side-car.
* `min_price_change_ppm` - This is currently not applicable to the side car and is not utilized.
* `exchange_config_json` - This is a list of exchanges that need to provide prices. Read more about the exchange conversions below.

## Exchange Conversions

There are four cases we consider when converting the exchange configurations:

### 1. Direct Conversion

This is the simplest case where the exchange ticker is the same as the market ticker. Example:

```json
{
    "id": 0,
    "pair": "BTC-USD",
    "exponent": -5,
    "min_exchanges": 3,
    "min_price_change_ppm": 1000,
    "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"}]}"
}
```

```json
{
  "exchanges": [
    {
      "exchangeName": "CoinbasePro",
      "ticker": "BTC-USD"
    }
  ]
}
```

In this case, we update the `BTC-USD`'s set of paths to include a single path with a single operation of CoinbasePro `BTC-USD` -> `BTC-USD`. Additionally we update the `BTC-USD`'s providers to include (CoinBasePro, BTC-USD).

Translated to the side-car, this would look like:

```json
{
  "tickers": {
    "BTC/USD": {
      "currency_pair": {
        "Base": "BTC",
        "Quote": "USD"
      },
      "decimals": 5,
      "min_provider_count": 3
    }
  },
  "paths": {
    "BTC/USD": {
      "paths": [
        {
          "operations": [
            {
              "currency_pair": {
                "Base": "BTC",
                "Quote": "USD"
              },
              "provider": "CoinbasePro"
            }
          ]
        }
      ]
    }
  },
  "providers": {
    "BTC/USD": {
      "providers": [
        {
          "name": "CoinbasePro",
          "off_chain_ticker": "BTC-USD"
        }
      ]
    }
  },
}

```

#### 2. Inverted Conversion

This is the case where the exchange ticker is inverted from the market ticker. Example:

```json
{
  "id": 1000000,
  "pair": "USDT-USD",
  "exponent": -9,
  "min_exchanges": 3,
  "min_price_change_ppm": 1000,
  "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"USDCUSDT\",\"invert\":true}]}"
}
```

```json
{
  "exchanges": [
    {
      "exchangeName": "Binance",
      "ticker": "USDCUSDT",
      "invert": true
    }
  ]
}
```

In this case, we update the `USDT-USD`'s set of paths to include a single path with a single operation of Binance `USDCUSDT` ^-1 -> `USDT-USD`. Additionally we update the `USDT-USD`'s providers to include (Binance, USDCUSDT).

Translated to the side-car, this would look like:

```json
{
  "tickers": {
    "USDT/USD": {
      "currency_pair": {
        "Base": "USDT",
        "Quote": "USD"
      },
      "decimals": 9,
      "min_provider_count": 3
    }
  },
  "paths": {
    "USDT/USD": {
      "paths": [
        {
          "operations": [
            {
              "currency_pair": {
                "Base": "USDT",
                "Quote": "USD"
              },
              "provider": "Binance",
              "invert": true
            }
          ]
        }
      ]
    }
  },
  "providers": {
    "USDT/USD": {
      "providers": [
        {
          "name": "Binance",
          "off_chain_ticker": "USDCUSDT"
        }
      ]
    }
  },
}
```

#### 3. Indirect Conversion

This is the case where the ticker is not directly associated with the market ticker, but can be converted to the market ticker using another market (the index price). Example:

```json
{
  "id": 0,
  "pair": "BTC-USD",
  "exponent": -5,
  "min_exchanges": 3,
  "min_price_change_ppm": 1000,
  "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
}
```

```json
{
  "exchanges": [
    {
      "exchangeName": "Okx",
      "ticker": "BTC-USDT",
      "adjustByMarket": "USDT-USD"
    }
  ]
}
```

In this case, we update the `BTC-USD`'s set of paths to include a single path with two operations of Okx `BTC-USDT` * Index `USDT-USD` -> `BTC-USD`. Additionally we update the `BTC-USD`'s providers to include (Okx, BTC-USDT).

Translated to the side-car, this would look like:

```json
{
  "tickers": {
    "BTC/USD": {
      "currency_pair": {
        "Base": "BTC",
        "Quote": "USD"
      },
      "decimals": 5,
      "min_provider_count": 3
    }
  },
  "paths": {
    "BTC/USD": {
      "paths": [
        {
          "operations": [
            {
              "currency_pair": {
                "Base": "BTC",
                "Quote": "USD"
              },
              "provider": "Okx"
            },
            {
              "currency_pair": {
                "Base": "USDT",
                "Quote": "USD"
              },
              "provider": "Index" // This is the index price (previously computed median price)
            }
          ]
        }
      ]
    }
  },
  "providers": {
    "BTC/USD": {
      "providers": [
        {
          "name": "Okx",
          "off_chain_ticker": "BTC-USDT"
        }
      ]
    }
  }
}
```

### 4. Indirect Inverted Conversion

This is the case where the ticker is not directly associated with the market ticker, but can be converted to the market ticker using another market, and the ticker is inverted. Example:

```json
{
  "id": 1000000,
  "pair": "USDT-USD",
  "exponent": -9,
  "min_exchanges": 3,
  "min_price_change_ppm": 1000,
  "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"BTC-USD\",\"invert\":true}]}"
}
```

```json
{
  "exchanges": [
    {
      "exchangeName": "Kucoin",
      "ticker": "BTC-USDT",
      "adjustByMarket": "BTC-USD",
      "invert": true
    }
  ]
}
```

In this case, we update the `USDT-USD`'s set of paths to include a single path with two operations of Kucoin (`BTC-USDT` * ^-1)  * Index `BTC-USD` -> `USDT-USD`. In this case, we DO NOT update the `USDT-USD`'s providers to include (Kucoin, BTC-USDT) as we assume the index price is provided by the side-car in combination with other providers.

We assume that the BTC/USD market will define the relevant providers including Kucoin separately.

Translated to the side-car, this would look like:

```json
{
  "tickers": {
    "USDT/USD": {
      "currency_pair": {
        "Base": "USDT",
        "Quote": "USD"
      },
      "decimals": 9,
      "min_provider_count": 3
    }
  },
  "paths": {
    "USDT/USD": {
      "paths": [
        {
          "operations": [
            {
              "currency_pair": {
                "Base": "BTC",
                "Quote": "USD"
              },
              "provider": "Kucoin",
              "invert": true
            },
            {
              "currency_pair": {
                "Base": "BTC",
                "Quote": "USD"
              },
              "provider": "Index" // This is the index price (previously computed median price)
            }
          ]
        }
      ]
    }
  },
  "providers": {
    "BTC/USD": { // Populated by the BTC/USD market
      "providers": [
        {
          "name": "Kucoin",
          "off_chain_ticker": "BTC-USDT"
        },
      ]
    },
    "USDT/USD": {
      "providers": []
    }
  }
}
```


### Additional Considerations

#### MEXC Ticker Conversions

The dYdX Mexc API uses the [Spot V2 endpoint](https://mexcdevelop.github.io/apidocs/spot_v2_en/#ticker-information) which has a different representation of tickers relative to the [Spot V3 websocket API](https://mexcdevelop.github.io/apidocs/spot_v3_en/#miniticker). V2 includes underscores that we omit in the V3 API.

#### Bitstamp Ticker Conversions

The dYdX Bitstamp API uses the [Ticker endpoint](https://www.bitstamp.net/api/v2/ticker/) which has a different representation of tickers relative to the [Websocket API](https://www.bitstamp.net/websocket/v2/). The ticker endpoint uses a `/` delimiter, while the websocket connection does not have a delimeter and lowercases the ticker.


