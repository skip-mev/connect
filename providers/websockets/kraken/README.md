# Kraken Provider

## Overview

The Kraken provider is used to fetch the ticker price from the [Kraken websocket API](https://docs.kraken.com/websockets/).


## General Considerations

* TLS with SNI (Server Name Indication) is required in order to establish a Kraken WebSockets API connection.
* All messages sent and received via WebSockets are encoded in JSON format
* All decimal fields (including timestamps) are quoted to preserve precision.
* Timestamps should not be considered unique and not be considered as aliases for transaction IDs. Also, the granularity of timestamps is not representative of transaction rates.
* Please use REST API endpoint [AssetPairs](https://docs.kraken.com/rest/#tag/Market-Data/operation/getTradableAssetPairs) to fetch the list of pairs which can be subscribed via WebSockets API. For example, field 'wsname' gives the supported pairs name which can be used to subscribe.
* **Recommended reconnection behaviour** is to (1) attempt reconnection instantly up to a handful of times if the websocket is dropped randomly during normal operation but (2) after maintenance or extended downtime, attempt to reconnect no more quickly than once every 5 seconds. There is no advantage to reconnecting more rapidly after maintenance during cancel_only mode.

To check all available pairs, you can use the following REST API call:

```bash
curl "https://api.kraken.com/0/public/Assets"
```
