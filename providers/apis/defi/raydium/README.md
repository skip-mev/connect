# Raydium Provider

The Raydium provider fetches prices from the Raydium dex via JSON-RPC requests to Solana nodes. 

## How It Works

For each ticker (i.e. RAY/SOL), we query 4 accounts:

* BaseTokenVault
* QuoteTokenVault
* AMMInfo
* OpenOrders

To calculate the price, we need to get the base and quote token balances, subtract PNL feels, and add the value of open orders.

With the above values, we calculate the price by dividing quote / base and multiplying by the scaling factor.

