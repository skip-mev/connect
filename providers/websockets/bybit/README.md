# ByBit Provider

## Overview

The ByBit provider is used to fetch the ticker price from the [ByBit websocket API](https://bybit-exchange.github.io/docs/v5/ws/connect).

Connections may be disconnected if a heartbeat ping is not sent to the server every 20 seconds to maintain the connection.


ByBit provides [public and private channels](https://bybit-exchange.github.io/docs/v5/ws/connect#how-to-subscribe-to-topics).

* Public channels -- No authentication is required, include tickers topic, K-Line topic, limit price topic, order book topic, and mark price topic etc.
* Private channels -- including account topic, order topic, and position topic, etc. -- require log in.

Users can choose to subscribe to one or more topic, and the total length of multiple topics cannot exceed 21,000 characters. This provider is implemented assuming that the user is only subscribing to public topics.

The exact topic that is used to subscribe to the ticker price is the [`Tickers`](https://bybit-exchange.github.io/docs/v5/websocket/public/ticker). This pushes data in real time if there are any price updates.

To retrieve all supported [spot markets](https://bybit-exchange.github.io/docs/v5/market/instrument), please run the following command:

```bash
curl "https://api.bybit.com/v5/market/instruments-info" 
```
