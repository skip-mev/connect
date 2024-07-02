# Gate.io Provider

## Overview

The Gate.io provider is used to fetch the ticker price from the [Gate.io websocket API](https://www.gate.io/docs/developers/apiv4/ws/en/#api-overview).

Gate.io provides [public](https://www.gate.io/docs/developers/apiv4/ws/en/#public-trades-channel) and [authenticated](https://www.gate.io/docs/developers/apiv4/ws/en/#funding-balance-channel) channels.

The Gate.io provider uses _protocol-level_ ping-pong, so no handlers need to be specifically implemented.

[Application level ping messages](https://www.gate.io/docs/developers/apiv4/ws/en/#application-ping-pong) can be sent which should be responded to with pong messages.

* Public channels -- No authentication is required, include tickers topic, K-Line topic, limit price topic, order book topic, and mark price topic etc.
* Private channels -- including account topic, order topic, and position topic, etc. -- require log in.

Users can choose to subscribe to one or more topic. This provider is implemented assuming that the user is only subscribing to public topics.

The exact topic that is used to subscribe to the ticker price is the [`Tickers`](https://www.gate.io/docs/developers/apiv4/ws/en/#tickers-channel). This pushes data every 1000ms.

To retrieve all supported [spot markets](https://www.gate.io/docs/developers/apiv4/en/#get-details-of-a-specific-currency), please run the following command:

```bash
curl -X GET https://api.gateio.ws/api/v4/spot/currency_pairs \
  -H 'Accept: application/json'
```
