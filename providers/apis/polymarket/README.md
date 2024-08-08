# PolyMarket Provider

Docs: https://docs.polymarket.com/

Polymarket is a web3 based prediction market. This provider uses the Polymarket CLOB REST API to access the off-chain order book. 

## How it Works

Polymarket uses [conditional outcome tokens](https://docs.gnosis.io/conditionaltokens/), a token that represents an outcome of a specific event. All tokens in Polymarket are denominated in terms of USD.

Tickers take the form of the `<market_slug>?<outcome>/USD`:

Example: `WILL_BERNIE_SANDERS_WIN_THE_2024_US_PRESIDENTIAL_ELECTION?YES/USD`

The offchain ticker is expected to be _just_ the token_id.

example: `95128817762909535143571435260705470642391662537976312011260538371392879420759`

The Provider can handle both the midpoint and the price endpoints. However, passing in multiple endpoints to the same provider will not yield additional data, as only the first endpoint is considered for the provider.

Example:

Midpoint:

`https://clob.polymarket.com/midpoint?token_id=95128817762909535143571435260705470642391662537976312011260538371392879420759`

Price:

`https://clob.polymarket.com/price?token_id=95128817762909535143571435260705470642391662537976312011260538371392879420759&side=BUY`

## Market Config

Below is an example of a market config for a single Polymarket token.

```json
 {
  "markets": {
    "WILL_BERNIE_SANDERS_WIN_THE_2024_US_PRESIDENTIAL_ELECTION?YES/USD": {
      "ticker": {
        "currency_pair": {
          "Base": "WILL_BERNIE_SANDERS_WIN_THE_2024_US_PRESIDENTIAL_ELECTION?YES",
          "Quote": "USD"
        },
        "decimals": 3,
        "min_provider_count": 1,
        "enabled": true
      },
      "provider_configs": [
        {
          "name": "polymarket_api",
          "off_chain_ticker": "95128817762909535143571435260705470642391662537976312011260538371392879420759"
        }
      ]
    }
  }
 }
```

## Oracle Config

Below is an example of an oracle config with a Polymarket provider.

```json
{
  "providers": {
    "polymarket_api": {
      "name": "polymarket_api",
      "type": "price_provider",
      "api": {
        "name": "polymarket_api",
        "enabled": true,
        "timeout": 3000000000,
        "interval": 500000000,
        "reconnectTimeout": 2000000000,
        "maxQueries": 1,
        "atomic": true,
        "endpoints": [
          {
            "url": "https://clob.polymarket.com/midpoint?token_id=%s",
            "authentication": {
              "apiKey": "",
              "apiKeyHeader": ""
            }
          }
        ],
        "batchSize": 0
      }
    }
  }
}
```

## Rate Limits

While not publicly available on the docs, the following information was given from a Polymarket developer regarding rate limits on the CLOB API:

* In general you can't do more than 20req/s
* You can't request orders and trades combined more than 20req/s
* You can't do more than 10 reqs to /book per second
