# ERC4626 Share Price Oracle Provider

Provides prices from contracts implementing [ERC4626SharePriceOracle.sol][1] by 
calling `getLatest()`.

> Beware: The `Base` denom is ignored by this provider as the actual quote/base
> pairing is determined by the provided contract address for each pair. 
> Operators must be careful to make sure the contract address actually matches 
> the denom pair it's configured for.

## Generate ABI bindings

Prereq: Install the abigen tool 

The ABI can be found [here][2]

The following command reads the abi from the current directory, generates 
`erc4626_share_price_oracle.go`, and names or prefixes the generated types with
"ERC4626SharePriceOracle".

```bash
abigen --abi=erc4626_share_price_oracle_abi.json \ 
    --pkg=erc4626_share_price_oracle \
    --out=erc4626_share_price_oracle.go \
    --type ERC4626SharePriceOracle
```

## Config

The `Symbol` field of each `TokenMetadata` is expected to be the
ERC4626SharePriceOracle contract address of the corresponding token. The price 
function takes no arguments, thus the decimals metadata is ignored by this
provider.

Because `isLatest()` returns both an instantanious price and a TWAP price, 
two config entries (and consequently two pair entries) are needed to capture
both prices. The `is_twap` field should be used to indicate which entry should
be used to record the time weighted price.

For example, Real Yield ETH's metadata entry would look like:

```toml
ryETH = {
    symbol = "0xb5b29320d2Dde5BA5BAFA1EbcD270052070483ec",
    is_twap = false
}
ryETH_TWAP = {
    symbol = "0xb5b29320d2Dde5BA5BAFA1EbcD270052070483ec",
    is_twap = true
}
```

The resulting contract function call for price retrieval would be `getLatest()` 

## Price Return Value

The `getLatest()` function returns the instantanious price, and time weighted
average price, and a boolean indicating whether the time weighted average answer
is safe to use. A `true` value means it is *not* safe to use, and both answers 
returned will be 0). "Safe" means the time weighted average price answer is not
too stale, or too new. If the price is not safe, an empty quote will be 
returned.

For more information, refer to [the contract][1].

[1]: https://github.com/PeggyJV/cellar-contracts/blob/97cc9070e6bdb3e55d8a3f074df5ac80d64f6be3/src/base/ERC4626SharePriceOracle.sol
[2]: https://gist.github.com/cbrit/75eb6aeffcaddc17c4f226d8a40db591

## RPC URL

By default, the `AlchemyURL` "https://eth-mainnet.g.alchemy.com/v2/" is recommended.
