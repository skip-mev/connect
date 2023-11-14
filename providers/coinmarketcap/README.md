# CoinMarketCap Provider

## Overview

> **NOTE:** This specific provider should not be used in any production setting as it is not sufficiently tested.

The CoinMarketCap provider is used to fetch the spot price for cryptocurrencies from the [CoinMarketCap API](https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyQuotesLatest). 

## Configuration

The configuration structure for this provider looks like the following:

```golang
// CoinMarketCapConfig is the config struct for the coinmarketcap provider.
type CoinMarketCapConfig struct {
	// APIKey is the API key used to make requests to the coinmarketcap API.
	APIKey string `mapstructure:"api_key" toml:"api_key"`
	// TokenNameToMetadata is a map of token names to their metadata.
	TokenNameToSymbol map[string]string `mapstructure:"token_name_to_symbol" toml:"token_name_to_symbol"`
}
```

Where a properly formatted `CoinMarketCapConfig` json object looks like the following:
    

Sample `coinmarketcap.toml`:
    
```toml
###############################################################################
###                           CoinMarketCap                                 ###
###############################################################################
# This section contains the configuration for the CoinMarketCap API. This meant to be
# used in a testing environment as this API is not production ready.

# APIKey is the API key used to make requests to the CoinMarketCap API.
api_key = "my-api-key"

# TokenNameToMetadata is a map of token names to their metadata.
[token_name_to_symbol]
  "BITCOIN" = "BTC"
  "USD" = "USD"
  # Add more token name-to-symbol mappings as needed
```
