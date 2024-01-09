# CoinGecko Provider

## Overview

> **NOTE:** This specific provider should not be used in any production setting as it is not sufficiently tested.

The CoinGecko provider is used to fetch the spot price for cryptocurrencies from the [CoinGecko API](https://www.coingecko.com/en/api). This provider can be configured to fetch with or without an API key. Note that without an API key, it is very likely that the CoinGecko API will rate limit your requests. 

## Configuration

The configuration structure for this provider looks like the following:

```golang
// Config is the config struct for the CoinGecko provider.
type Config struct {
	// APIKey is the API key used to make requests to the CoinGecko API.
	APIKey string `json:"apiKey" validate:"required"`

	// SupportedBases maps the base currencies to the CoinGecko API's supported
	// base currencies.
	SupportedBases map[string]string `json:"supportedBases" validate:"required"`

	// SupportedQuotes is the list of supported quotes for the CoinGecko API.
	SupportedQuotes map[string]string `json:"supportedQuotes" validate:"required"`
}
```

Where a properly formatted `CoinGeckoConfig` json object looks like the following:

Sample `coingecko.json`:
    
```json
{
  "apiKey": "",
  "supportedBases": {
    "BITCOIN": "bitcoin",
    "ETHEREUM": "ethereum",
    "ATOM": "cosmos",
    "SOLANA": "solana",
    "POLKADOT": "polkadot",
    "DYDX": "dydx-chain"
  },
  "supportedQuotes": {
    "USD": "usd",
    "ETHEREUM": "eth"
  }
}
```

The CoinGecko API fetches an aggregated TWAP price any given currency pair.

## Supported Bases

To determine the base currencies that the CoinGecko provider supports, you can run the following command:

```bash
$ curl -X GET https://api.coingecko.com/api/v3/coins/list
```

## Supported Quotes

To determine the quote currencies that the CoinGecko provider supports, you can run the following command:

```bash
$ curl -X GET https://api.coingecko.com/api/v3/simple/supported_vs_currencies
```
