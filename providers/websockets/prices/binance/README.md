# Binance Provider

## Overview

The Binance provider is used to fetch the ticker price from the [Binance websocket API](https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams). A single connection is only valid for 24 hours; after that, a new connection must be established. The Websocket server will send a ping frame every 3 minutes. If the client does not receive a pong frame within 10 minutes, the connection will be closed. Note that all symbols are in lowercase.

The WebSocket connection has a limit of 5 incoming messages per second. A message is considered:

* A Ping frame
* A Pong frame
* A JSON controlled message (e.g. a subscription)
* A connection that goes beyond the rate limit will be disconnected. IPs that are repeatedly disconnected for going beyond the rate limit may be banned for a period of time.

A single connection can listen to a maximum of 1024 streams. If a user attempts to listen to more streams, the connection will be disconnected. There is a limit of 300 connections per attempt every 5 minutes per IP.

The specific channels / streams that are subscribed to is the [Aggregate Trade Stream](https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams) and the [Ticker Stream](https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams). The Aggregate Trade Streams push trade information that is aggregated for a single taker order in real time. The ticker stream pushes the ticker spot price every second.
