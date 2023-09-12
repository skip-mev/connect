# ERC4626 Provider

Provides prices from contracts implementing ERC4626 by calling `previewRedeem()`
with the corresponding pair's correct decimals as defined in the metadata map.

> Beware: The `Base` denom is ignored by this provider as the actual quote/base
> pairing is determined by the provided contract address for each pair. 
> Operators must be careful to make sure the contract address actually matches 
> the denom pair it's configured for.

## Generate ABI bindings

Prereq: Install the abigen tool 

The ABI can be found 
[here](https://gist.github.com/cbrit/9d657ac2b08a7237df551f0fce3bfbfe)

The following command reads the abi from the current directory, generates 
`erc4626.go`, and names or prefixes the generated types with "ERC4626".

```bash
abigen --abi=erc4626_abi.json --pkg=erc4626 --out=erc4626.go --type ERC4626
```

## Config

The `Symbol` field of each `TokenMetadata` is expected to be the ERC4626 
contract address of the corresponding token. It is important that the `Decimals`
value is accurate, as the price reported will be scaled based on this value. 

For example, Real Yield ETH uses 18 decimal places. It's metadata entry would 
look like:

```toml
[oracle.providers.token_name_to_metadata]
ryETH = {
    symbol = "0xb5b29320d2Dde5BA5BAFA1EbcD270052070483ec",
    decimals = 18
}
```

The resulting contract function call for price retrieval would be 
`previewRedeem(1000000000000000000)`. 

