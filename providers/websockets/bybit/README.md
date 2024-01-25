# OKX Provider

## Overview

The OKX provider is used to fetch the ticker price from the [ByBit web socket API](https://linear.app/skip/issue/BLO-733/bybit-provider). The total number of connections that can be opened is 500 per 5 minutes.  Spot subscription requests can only be made with up to 10 arguments per request.

Connections may be disconnected if a heartbeat ping is not sent to the server every 20 seconds to maintain the connection.

```text
If there’s a network problem, the system will automatically disable the connection.

The connection will break automatically if the subscription is not established or data has not been pushed for more than 30 seconds.

To keep the connection stable:

1. Set a timer of N seconds, where N is less than 30.
2. If the timer is triggered, send the String 'ping'.
3. Expect a 'pong' as a response. If the response message is not received within N seconds, please raise an error or reconnect.
```

ByBit provides [public and private channels](https://www.okx.com/docs-v5/en/?shell#overview-websocket-subscribe).

* Public channels -- No authentication is required, include tickers topic, K-Line topic, limit price topic, order book topic, and mark price topic etc.
* Private channels -- including account topic, order topic, and position topic, etc -- require log in.

Users can choose to subscribe to one or more topic, and the total length of multiple topics cannot exceed 21,000 characters. This provider is implemented assuming that the user is only subscribing to public topics.

The exact topic that is used to subscribe to the ticker price is the [`Tickers`](https://bybit-exchange.github.io/docs/v5/websocket/public/ticker). This pushes data in real time if there are any price updates.

To retrieve all of the supported [spot markets](https://bybit-exchange.github.io/docs/v5/market/instrument), please run the following command:

```bash
curl "https://api.bybit.com/v5/market/instruments-info" 
```

## Configuration

The configuration structure for this provider looks like the following:

```golang
type Config struct {
	// Markets is the list of markets to subscribe to. The key is the currency pair and the value
	// is the pair ID. The pair ID must correspond to the spot market. For example,
	// the pair ID for the BITCOIN/USDT market is BTCUSDT.
	Markets map[string]string `json:"markets"`

	// Production is true if the config is for production.
	Production bool `json:"production"`
}
