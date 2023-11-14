# Coinbase Provider

## Overview

> **NOTE:** This specific provider should not be used in any production setting as it may be rate limited by Coinbase.

The Coinbase provider is used to fetch the spot price for cryptocurrencies from the [Coinbase API](https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-currencies). 

## Configuration

The configuration structure for this provider looks like the following:

```golang
// CoinbaseConfig is the configuration for the Coinbase provider.
type CoinbaseConfig struct {
	// NameToSymbol is a map of currency names to their symbols.
	NameToSymbol map[string]string `mapstructure:"name_to_symbol" toml:"name_to_symbol"`
}
```


Sample `coinbase.toml`:
    
```toml
###############################################################################
###                                Coinbase                                 ###
###############################################################################
# This section contains the configuration for the Coinbase API. This meant to be
# used in a testing environment as this API is not production ready.

# NameToSymbol is a map of currency names to their symbols.
[name_to_symbol]
  "BITCOIN" = "BTC"
  "USD" = "USD"
  # Add more currency mappings as needed

```

This allows the currency pairs that are passed in to the provider to be correctly mapped to the Coinbase API. The provider will then return the spot price for the currency pair.
