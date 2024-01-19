# Binance Provider

## Overview

The Binance provider is used to fetch the spot price for cryptocurrencies from the [Binance API](https://binance-docs.github.io/apidocs/spot/en/#general-info). Note that only supported markets can be added to this provider's market configurations, otherwise all requests will fail. To see all of the supported markets please visit ,.....

## Configuration

The configuration structure for this provider looks like the following:

```golang
// Config is the config struct for the Binance provider.
type Config struct {
	// SupportedBases maps the base currencies to the Binance API's supported
	// base currencies.
	SupportedBases map[string]string `json:"supportedBases" validate:"required"`

	// SupportedQuotes is the list of supported quotes for the Binance API.
	SupportedQuotes map[string]string `json:"supportedQuotes" validate:"required"`
}
```

Where a properly formatted `BinanceConfig` json object looks like the following:

Sample `binance.json`:

```json
{
  "supportedBases": {
    "BITCOIN": "BTC",
    "ETHEREUM": "ETH",
    "ATOM": "ATOM",
    "SOLANA": "SOL",
    "POLKADOT": "DOT",
    "DYDX": "DYDX"
  },
  "supportedQuotes": {
    "USD": "USDT",
    "ETHEREUM": "ETH"
  }
}
```

The Binance API fetches an aggregated TWAP price any given currency pair.

## Supported Pairs

To determine the pairs (in the form `BASEQUOTE`) currencies that the Binance provider supports, you can run the following command:

```bash
$ curl -X GET https://api.binance.vision/api/v3/ticker/price         
```
