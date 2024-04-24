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
* `pair` - This is the pair that the market is associated with. There is a direct conversion here assuming the delimiter is `-` across all markets.
* `exponent` - In the side-car we represent all prices in `Decimals`, which is the absolute value of the `exponent`.
* `min_exchanges` - This is the minimum number of exchanges that need to provide prices for the market to be considered valid. This is a direct mapping to `MinProviderCount` in the side-car.
* `min_price_change_ppm` - This is currently not applicable to the side car and is not utilized.

* `exchange_config_json` - This is a list of exchanges that need to provide prices. We have a very similar structure in the side-car, but we use a different format. The side-car uses a `MarketMap` to store the exchange and ticker pairs.


### Additional Considerations

#### MEXC Ticker Conversions

The dYdX Mexc API uses the [Spot V2 endpoint](https://mexcdevelop.github.io/apidocs/spot_v2_en/#ticker-information) which has a different representation of tickers relative to the [Spot V3 websocket API](https://mexcdevelop.github.io/apidocs/spot_v3_en/#miniticker). V2 includes underscores that we omit in the V3 API.

#### Bitstamp Ticker Conversions

The dYdX Bitstamp API uses the [Ticker endpoint](https://www.bitstamp.net/api/v2/ticker/) which has a different representation of tickers relative to the [Websocket API](https://www.bitstamp.net/websocket/v2/). The ticker endpoint uses a `/` delimiter, while the websocket connection does not have a delimiter and lowercases the ticker.


