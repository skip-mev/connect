# PolyMarket Provider

Docs: https://docs.polymarket.com/

Polymarket is a web3 based prediction market. This provider uses the Polymarket CLOB REST API to access the off-chain order book. 

## How it Works

Polymarket uses [conditional outcome tokens](https://docs.gnosis.io/conditionaltokens/), a token that represents an outcome of a specific event. All tokens in Polymarket are denominated in terms of USD.

Tickers take the form of:

`<POLYMARKET_TOKEN_ID>/USD`

example: `21742633143463906290569050155826241533067272736897614950488156847949938836455/USD`

The offchain ticker is expected to be _just_ the token_id.

The Provider simply calls the `/price` endpoint of the CLOB API. There are two query parameters:
* token_id
* side

Side can be either `buy` or `sell`. For this provider, we hardcode the side to `buy`.

## Market Config

Below is an example of a market config for a single Polymarket token.

```json
 {
  "markets": {
   "21742633143463906290569050155826241533067272736897614950488156847949938836455/USD": {
    "ticker": {
     "currency_pair": {
      "Base": "21742633143463906290569050155826241533067272736897614950488156847949938836455",
      "Quote": "USD"
     },
     "decimals": 3,
     "min_provider_count": 1,
     "enabled": true
    },
    "provider_configs": [
     {
      "name": "polymarket_api",
      "off_chain_ticker": "21742633143463906290569050155826241533067272736897614950488156847949938836455"
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
            "url": "https://clob.polymarket.com/price?token_id=%s&side=BUY",
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