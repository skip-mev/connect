# Coinbase Provider

## Overview

> **NOTE:** This specific provider should not be used in any production setting as it may be rate limited by Coinbase.

The Coinbase provider is used to fetch the spot price for cryptocurrencies from the [Coinbase API](https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-prices#get-spot-price). 

## Configuration

The configuration structure for this provider looks like the following:

```golang
// Config is the configuration for the Coinbase provider.
type Config struct {
	// SymbolMap maps the oracle's equivalent of an asset to the expected coinbase
	// representation of the asset.
	SymbolMap map[string]string `json:"symbolMap" validate:"required"`
}
```


Sample `coinbase.json`:
    
```json
{
  "symbolMap": {
      "BITCOIN": "BTC",
      "USD": "USD",
      "ETHEREUM": "ETH",
      "ATOM": "ATOM",
      "SOLANA": "SOL",
      "POLKADOT": "DOT",
      "DYDX": "DYDX"
  }
}


```

This allows the currency pairs that are passed in to the provider to be correctly mapped to the Coinbase API. The provider will then return the spot price for the currency pair.
