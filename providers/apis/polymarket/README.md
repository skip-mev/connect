# PolyMarket Provider

Docs: https://docs.polymarket.com/

Polymarket is a web3 based prediction market. This provider uses the Polymarket CLOB REST API to access the off-chain order book. 

## How it Works

Polymarket uses [conditional outcome tokens](https://docs.gnosis.io/conditionaltokens/), a token that represents an outcome of a specific event. All tokens in Polymarket are denominated in terms of USD.

We suggest tickers take the form of the `<market_slug>?<outcome>/USD`. However, tickers are ignored by the polymarket provider, and they can be whatever arbitrary data that suits your use case. The ONLY required text is your ticker must end in `/USD`. 

Example: `WILL_BERNIE_SANDERS_WIN_THE_2024_US_PRESIDENTIAL_ELECTION?YES/USD`
Example2: `BernieBecomesPresident/USD`

The offchain ticker **must** be <market_id>/<token_id>

example: `0x08f5fe8d0d29c08a96f0bc3dfb52f50e0caf470d94d133d95d38fa6c847e0925/95128817762909535143571435260705470642391662537976312011260538371392879420759`

The Provider queries the `/markets` endpoint, and looks for the token_id in the response. The provider will throw an error if the token_id in the offchain ticker is not present in the response data.

Example:

`https://clob.polymarket.com/markets/0xc6485bb7ea46d7bb89beb9c91e7572ecfc72a6273789496f78bc5e989e4d1638`


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
          "off_chain_ticker": "0x08f5fe8d0d29c08a96f0bc3dfb52f50e0caf470d94d133d95d38fa6c847e0925/95128817762909535143571435260705470642391662537976312011260538371392879420759"
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
        "atomic": false,
        "endpoints": [
          {
            "url": "https://clob.polymarket.com/markets/%s",
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
