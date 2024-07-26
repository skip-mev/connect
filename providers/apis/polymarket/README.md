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