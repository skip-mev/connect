# Bitstamp Provider

## Overview

Bitstamp is a cryptocurrency exchange that provides a free API for fetching cryptocurrency data. Bitstamp is a **primary data source** for the oracle. It supports connecting to a websocket without authentication. Once you open a connection via websocket handshake (using HTTP upgrade header), you can subscribe to desired channels. After this is accomplished, you will start to receive a stream of live events for every channel you are subscribed to. Maximum connection age is 90 days from the time the connection is established. When that period of time elapses, you will be automatically disconnected and will need to re-connect.

To see the supported set of markets, you can reference the websocket documentation [here](https://www.bitstamp.net/websocket/v2/).
