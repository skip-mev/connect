# Crypto.com Provider

## Overview

The Crypto.com provider is used to fetch the ticker price from the [Crypto.com web socket API](https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#ticker-instrument_name). The websocket is [rate limited](https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#rate-limits) with a maximum of 100 requests per second. This provider does not require any API keys. To determine the acceptable set of base and quote currencies, you can reference the [get instruments API](https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#reference-and-market-data-api).

To better distribute system load, a single market data websocket connection is limited to a maximum of 400 subscriptions. Once this limit is reached, further subscription requests will be rejected with the EXCEED_MAX_SUBSCRIPTIONS error code. A user should establish multiple connections if additional market data subscriptions are required.

## Configuration

The configuration structure for this provider looks like the following:

```golang
// Config is the configuration for the Crypto.com provider. To access the available
// markets, please check the following link:
// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#reference-and-market-data-api
type Config struct {
	// Markets is a map of currency pair to perpetual market ID. For an example of the
	// how to configure the markets, please check the readme.
	Markets map[string]string `json:"markets" validate:"required"`

	// Production is true if the provider is running in production mode. Note that if the
	// production setting is set to false, all prices returned by any subscribed markets
	// will be static.
	Production bool `json:"production" validate:"required"`
}
```

Note that if production is set to false, all prices returned by any subscribed markets will be static. A sample configuration is shown below:

```json
{
    "markets": {
        "BITCOIN/USD": "BTCUSD-PERP", // Perpetual market
        "ETHEREUM/USD": "ETHUSD-PERP", // Perpetual market
        "SOLANA/USD": "SOLUSD-PERP", // Perpetual market
        "ATOM/USD": "ATOMUSD-PERP", // Perpetual market
        "POLKADOT/USD": "DOTUSD-PERP", // Perpetual market
        "DYDX/USD": "DYDXUSD-PERP", // Perpetual market
        "ETHEREUM/BITCOIN": "ETH_BTC" // Spot market
    },
    "production": true
}
```

The names of the markets (BTCUSD-PERP vs. BTC_USD) represent the perpetual vs. spot markets.
