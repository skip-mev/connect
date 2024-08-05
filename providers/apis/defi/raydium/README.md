# Raydium Provider

The Raydium provider fetches prices from the Raydium dex via JSON-RPC requests to Solana nodes. 

## How It Works

For each ticker (i.e. RAY/SOL), we query 4 accounts:

* BaseTokenVault
* QuoteTokenVault
* AMMInfo
* OpenOrders

We get the token balances for both the base and quote tokens. This involves getting the balance, subtracting PNL feels, and finally adding the value of open orders.

With the values for the base and quote tokens, we calculate the price by dividing quote / base and multiplying by the scaling factor.

