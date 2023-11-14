# CoinGecko Provider

## Overview

> **NOTE:** This specific provider should not be used in any production setting as it is not sufficiently tested.

The CoinGecko provider is used to fetch the spot price for cryptocurrencies from the [CoinGecko API](https://www.coingecko.com/en/api). 

## Configuration

The configuration structure for this provider looks like the following:

```golang
// CoinGeckoConfig is the config struct for the CoinGecko provider.
type CoinGeckoConfig struct {
	// APIKey is the API key used to make requests to the CoinGecko API.
	APIKey string `mapstructure:"api_key" toml:"api_key"`
}
```

Where a properly formatted `CoinGeckoConfig` json object looks like the following:

Sample `coingecko.toml`:
    
```toml
###############################################################################
###                               CoinGecko                                 ###
###############################################################################
# This section contains the configuration for the CoinGecko API. This meant to be
# used in a testing environment as this API is not production ready.

# APIKey is the API key used to make requests to the CoinGecko API.
api_key = "my-api-key"
```

The CoinGecko API fetches an aggregated TWAP price any given currency pair.
