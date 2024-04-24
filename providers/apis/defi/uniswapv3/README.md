# Uniswap v3 API Provider

> Please read over the [Uniswap v3 documentation](https://blog.uniswap.org/uniswap-v3-math-primer) to understand the basics of Uniswap v3.

## Overview

The Uniswap v3 API Provider allows you to interact with the Uniswap v3 pools - otherwise known as concentrated liquidity pools (CLPs) - on the Ethereum blockchain. The provider utilizes JSON-RPC to interact with an ethereum node - batching multiple requests into a single HTTP request to reduce latency and improve performance.

Uniswap v3 shows the current price of the pool in `slot0` of the pool contract. `slot0` is where most of the commonly accessed values are stored, making it a good starting point for data collection. You can get the price from two places; either from the `sqrtPriceX96` or calculating the price from the pool `tick` value. Using `sqrtPriceX96` should be preferred over calculating the price from the current tick, because the current tick may lose precision due to the integer constraints. As such, this provider uses the `sqrtPriceX96` value to calculate the price of the pool.

Based on the [analysis](https://docs.chainstack.com/docs/http-batch-request-vs-multicall-contract#performance-comparison) of various approaches for querying EVM state, this implementation utilizes `BatchCallContext` available on any client that implements the go-ethereum's `ethclient` interface. This allows for multiple requests to be batched into a single HTTP request, reducing latency and improving performance. This is preferable to using the `multicall` contract, which is a contract that aggregates multiple calls into a single call.

To generate the ABI for the Uniswap v3 pool contract, you can use the `abigen` tool provided by the go-ethereum library. The ABI is used to interact with the Uniswap v3 pool contract.

```bash
abigen --sol ./contracts/UniswapV3Pool.sol --pkg uniswap --out ./uniswap_v3_pool.go
```
