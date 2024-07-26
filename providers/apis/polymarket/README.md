# PolyMarket Provider

Docs: https://docs.polymarket.com/

Polymarket is a web3 based prediction market. This provider uses their CLOB API which access their off-chain orderbook. 

## How it Works

Polymarket uses [conditional outcome tokens](https://docs.gnosis.io/conditionaltokens/), a token that represents an outcome of a specific event. All tokens in Polymarket are denominated in terms of USD.

Naturally, tickers should take the form of:

`<TOKEN_ID>/USD`

example: `21742633143463906290569050155826241533067272736897614950488156847949938836455/USD`

The offchain ticker is expected to be _just_ the token_id.

The Provider simply calls the `/price` endpoint. There are two query parameters:
* token_id
* side

Side can be either `buy` or `sell`. For this provider, we hardcode the side to `buy`.