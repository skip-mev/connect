# OKX Provider

## Overview

The OKX provider is used to fetch the ticker price from the [OKX websocket API](https://www.okx.com/docs-v5/en/#overview-websocket-overview). The websocket request size for data transmission between the client and server is only 2 bytes. The total number of requests to subscribe to new markets is limited to 3 requests per second. The total number of 'subscribe'/'unsubscribe'/'login' requests per connection is limited to 480 times per hour. WebSocket login and subscription rate limits are based on connection.

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

To retrieve all supported [spot markets](https://www.okx.com/docs-v5/en/?shell#public-data-rest-api-get-instruments), please run the following command:

```bash
curl https://www.okx.com/api/v5/public/instruments?instType=SPOT
```
