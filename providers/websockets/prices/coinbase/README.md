# Coinbase Provider

## Overview

The Coinbase WebSocket feed is publicly available and provides real-time market data updates for orders and trades. Two endpoints are supported in both production and sandbox. Coinbase Market Data is the traditional feed which is available without authentication. New message types can be added at any time. Clients are expected to ignore messages they do not support. To begin receiving feed messages, you must send a subscribe message to the server indicating which channels and products to receive. This message is mandatory — you are disconnected if no subscribe has been received within 5 seconds.

The Coinbase Websocket feed enables websocket compression. Websocket compression, defined in RFC7692, compresses the payload of WebSocket messages which can increase total throughput and potentially reduce message delivery latency. The permessage-deflate extension can be enabled by adding the extension header. Currently, it is not possible to specify the compression level.

### Sequence Numbers

Most feed messages contain a sequence number. Sequence numbers are increasing integer values for each product, with each new message being exactly one sequence number greater than the one before it.Sequence numbers that are greater than one integer value from the previous number indicate that a message has been dropped. Sequence numbers that are less than the previous number can be ignored or represent a message that has arrived out of order.

### Rate Limits

Real-time market data updates provide the fastest insight into order flow and trades. This means that you are responsible for reading the message stream and using the message relevant for your needs—this can include building real-time order books or tracking real-time trades.

* Requests per second per IP: 8
* Requests per second per IP in bursts: Up to 20
* Messages sent by the client every second per IP: 100

### Other Considerations

* Connected clients should increase their websocket receive buffer to the largest configurable amount possible (given any client library or infrastructure limitations), due to the potential volume of data for any given product.
* Space out websocket requests to adhere to the above rate limits.

To determine all markets available, you can use the [Get Products](https://docs.pro.coinbase.com/#get-products) API call.
