# BitFinex Provider

## Overview

The BitFinex provider is used to fetch the ticker price from the [BitFinex websocket API](https://docs.bitfinex.com/docs/ws-general). The total amount of subscriptions per connection is [30](https://docs.bitfinex.com/docs/ws-general#how-to-connect).

BitFinex provides [public](https://docs.bitfinex.com/docs/ws-public) and [private (authenticated)](https://docs.bitfinex.com/docs/ws-auth) channels.

* Public channels -- No authentication is required, include tickers channel, K-Line channel, limit price channel, order book channel, and mark price channel etc.
* Private channels -- including account channel, order channel, and position channel, etc -- require log in.

The exact channel that is used to subscribe to the ticker price is the [`Tickers`](https://docs.bitfinex.com/reference/ws-public-ticker). This pushes data regularly regarding a ticker status.

To retrieve all supported [tickers](https://docs.bitfinex.com/reference/rest-public-tickers), please run the following command:

```bash
curl https://api-pub.bitfinex.com/v2/conf/pub:list:currency
```
