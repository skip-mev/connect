# OKX Provider

## Overview

The OKX provider is used to fetch the ticker price from the [OKX web socket API](https://www.okx.com/docs-v5/en/#overview-websocket-overview). The web socket request size for data transmission between the client and server is only 2 bytes. The total number of requests to subscribe to new markets is limited to 3 requests per second. The total number of 'subscribe'/'unsubscribe'/'login' requests per connection is limited to 480 times per hour. WebSocket login and subscription rate limits are based on connection.

Connections will break automatically if the subscription is not established or data has not been pushed for more than 30 seconds. [Per OKX documentation](https://www.okx.com/docs-v5/en/#overview-websocket-overview),

```text
If there’s a network problem, the system will automatically disable the connection.

The connection will break automatically if the subscription is not established or data has not been pushed for more than 30 seconds.

To keep the connection stable:

1. Set a timer of N seconds whenever a response message is received, where N is less than 30.
2. If the timer is triggered, which means that no new message is received within N seconds, send the String 'ping'.
3. Expect a 'pong' as a response. If the response message is not received within N seconds, please raise an error or reconnect.
```

OKX provides [public and private channels](https://www.okx.com/docs-v5/en/?shell#overview-websocket-subscribe). 

* Public channels -- No authentication is required, include tickers channel, K-Line channel, limit price channel, order book channel, and mark price channel etc.
* Private channels -- including account channel, order channel, and position channel, etc -- require log in.

Users can choose to subscribe to one or more channels, and the total length of multiple channels cannot exceed 64 KB. This provider is implemented assuming that the user is only subscribing to public channels.

The exact channel that is used to subscribe to the ticker price is the [`Index Tickers Channel`](https://www.okx.com/docs-v5/en/?shell#public-data-websocket-index-tickers-channel). This pushes data every 100ms if there are any price updates, otherwise it will push updates once a minute.

To retrieve all of the supported [spot markets](https://www.okx.com/docs-v5/en/?shell#public-data-rest-api-get-instruments), please run the following command:

```bash
curl https://www.okx.com/api/v5/public/instruments?instType=SPOT
```

## Configuration

The configuration structure for this provider looks like the following:

```golang
type Config struct {
	// Markets is the list of markets to subscribe to. The key is the currency pair and the value
	// is the instrument ID. The instrument ID must correspond to the spot market. For example,
	// the instrument ID for the BITCOIN/USDT market is BTC-USDT.
	Markets map[string]string `json:"markets"`

	// Production is true if the config is for production.
	Production bool `json:"production"`
}
```

Note that if production is set to false, all prices returned by any subscribed markets will be static. A sample configuration is shown below:

```json
{
    "markets": {
        "BITCOIN/USD": "BTC-USD", // Spot market
        "ETHEREUM/USD": "ETH-USD", // Spot market
        "SOLANA/USD": "SOL-USD", // Spot market
        "ATOM/USD": "ATOM-USD", // Spot market
        "POLKADOT/USD": "DOT-USD", // Spot market
        "DYDX/USD": "DYDX-USD", // Spot market
        "ETHEREUM/BITCOIN": "ETH-BTC" // Spot market
    },
    "production": true
}
```
