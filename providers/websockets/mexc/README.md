# MEXC Provider

## Overview

The MEXC provider is a websocket provider that fetches data from the MEXC exchange API. All documentation for the websocket can be found [here](https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-market-streams).


## Considerations

* A single connection to the MEXC API is made and remains valid for 24 hours before disconnecting and reconnecting. 
* All ticker symbols must be in uppercase in the market configuration eg: `spot@public.deals.v3.api@<symbol>` -> `spot@public.deals.v3.api@BTCUSDT`.
* If there is no valid websocket subscription, the server will disconnect in 30 seconds. If the subscription is successful but there is no streams, the server will disconnect in 1 minute. The client can send PING to maintain the connection.
* Every websocket connection can support a maximum of 30 subscriptions. If the client needs to subscribe to more than 30 streams, it needs to connect multiple websocket connections.

To determine all supported markets, you can run the following command:

```bash
curl https://api.mexc.com/api/v3/defaultSymbols 
```
