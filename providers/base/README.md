# Base Provider

## Overview

The Base Provider is a special provider that is used to provide the base functionality for all other providers. It is not intended to be used directly, but instead serves as a base for other providers to extend.

The base provider is responsible for the following:

* Implementing the oracle `Provider` interface directly just by inheriting the base provider.
* Having the provider data available in constant time when requested.

Each base provider implementation will be run in a separate goroutine by the main oracle process. This allows the provider to fetch data from the underlying data source asynchronously. The base provider will then store the data in a thread safe map. The main oracle service utilizing this provider can determine if the data is stale or not based on result timestamp associated with each data point.

The base provider constructs a response channel that it is always listening to and making updates as needed. Every interval, the base provider will fetch the data from the underlying data source and send the response to the response channel, respecting the number of concurrent requests to the rate limit parameters of the underlying source (if it has any).

![Architecture Overview](./architecture.png)


## API (HTTP) Based Providers

In order to implement API based providers, you must implement the [`APIDataHandler`](./api/handlers/api_data_handler.go) interface and the [`RequestHandler`](./api/handlers/request_handler.go) interfaces. The `APIDataHandler` is responsible for creating the URL to be sent to the HTTP client and parsing the response from the HTTP response. The `RequestHandler` is responsible for making the HTTP request and returning the response.

Once these two interfaces are implemented, you can then instantiate an [`APIQueryHandler`](./api/handlers/api_query_handler.go) and pass it to the base provider. The `APIQueryHandler` is abstracts away the logic for making the HTTP request and parsing the response. The base provider will then take care of the rest. The responses from the `APIQueryHandler` are sent to the base provider via a buffered channel. The base provider will then store the data in a thread safe map. To read more about the various API provider configurations available, please visit the [API provider configuration](../../oracle/config/api.go) documentation.

Alternatively, you can directly implement the [`APIFetcher`](./api/handlers/api_query_handler.go) interface. This is appropriate if you want to abstract over the various processes of interacting with GRPC, JSON-RPC, REST, etc. APIs.

### APIDataHandler

The `APIDataHandler` interface is primarily responsible for constructing the URL that will fetch the desired data and parsing the response. The interface is purposefully built with generics in mind. This allows the provider to fetch data of any type from the underlying data source.

```golang
// APIDataHandler defines an interface that must be implemented by all providers that
// want to fetch data from an API using HTTP requests. This interface is meant to be
// paired with the APIQueryHandler. The APIQueryHandler will use the APIDataHandler
// to create the URL to be sent to the HTTP client and to parse the response from the
// API.
type APIDataHandler[K providertypes.ResponseKey, V providertypes.providertypes.ResponseValue] interface {
	CreateURL(ids []K) (string, error)
	ParseResponse(ids []K, response *http.Response) GetResponse[K, V]
	Atomic() bool
	Name() string
}
```

#### Determining K and V

First developers must determine the type of data that they want to fetch from the underlying data source. This can be any type that is supported by the oracle. For example, the simplest example is price data for a given currency pair (base / quote). The `K` type would be the currency pair and the `V` type would be the price data.

```golang
APIDataHandler[types.ProviderTicker, *big.Float]
```

#### CreateURL

The `CreateURL` function is responsible for creating the URL that will be sent to the HTTP client. The function should utilize the IDs passed in as references to the data that needs to be fetched. For example, if the data source requires a currency pair to be passed in, the `CreateURL` function should use the currency pair to construct the URL.

#### ParseResponse

The `ParseResponse` function is responsible for parsing the response from the API. The response should be parsed into a map of IDs to results. If any IDs are not resolved, they should be returned in the unresolved map. The timestamp associated with the result should reflect either the time the data was fetched or the time the API last updated the data.

#### Atomic

The `Atomic` function is used to determine whether the handler can make a single request for all IDs or multiple requests for each ID. If true, the handler will make a single request for all IDs. If false, the handler will make a request for each ID.

### RequestHandler

The request handler is responsible for making the HTTP request and returning the response.

```golang
// RequestHandler is an interface that encapsulates sending a request to a data provider.
type RequestHandler interface {
	Do(ctx context.Context, url string) (*http.Response, error)
}
```

#### Do

The `Do` function is responsible for making the HTTP request and returning the response.

This interface is particularly useful if a custom HTTP client is needed. For example, if the data provider requires a custom header to be sent with the request, the `RequestHandler` can be used to implement this logic.

### APIFetcher

The `APIFetcher` interface is used to fetch data from the underlying data source. This interface is used by the `APIQueryHandler` to encapsulate the logic for fetching data - with metrics collection and more.

```golang
// APIFetcher is an interface that encapsulates fetching data from a provider. This interface
// is meant to abstract over the various processes of interacting w/ GRPC, JSON-RPC, REST, etc. APIs.
type APIFetcher[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	Fetch(ctx context.Context,ids []K) providertypes.GetResponse[K, V]
}
```

## Websocket-Based Providers

In order to implement websocket-based providers, you must implement the [`WebSocketDataHandler`](./websocket/handlers/ws_data_handler.go) interface and the [`WebSocketConnHandler`](./websocket/handlers/ws_conn_handler.go) interfaces. The `WebSocketDataHandler` is responsible for parsing messages from the websocket connection, constructing heartbeats, and constructing the initial subscription message(s). This handler must manage all state associated with the websocket connection i.e. connection identifiers. The `WebSocketConnHandler` is responsible for making the websocket connection and maintaining it - including reads, writes, dialing, and closing.

Once these two interfaces are implemented, you can then instantiate an [`WebSocketQueryHandler`](./websocket/handlers/ws_query_handler.go) and pass it to the base provider. The `WebSocketQueryHandler` abstracts away the logic for connecting, reading, sending updates, and parsing responses all using the two interfaces above. The base provider will then take care of the rest - including storing the data in a thread safe manner. To read more about the various configurations available for websocket providers, please visit the [websocket provider configuration](../../oracle/config/websocket.go) documentation.

### WebSocketDataHandler

The `WebSocketDataHandler` interface is primarily responsible for constructing the initial set of subscription messages, parsing messages received from the websocket connection, and constructing heartbeat updates. The interface is purposefully built with generics in mind. This allows the provider to fetch data of any type from the underlying data source.

```golang
// WebSocketDataHandler defines an interface that must be implemented by all providers that
// want to fetch data from a websocket. This interface is meant to be paired with the
// WebSocketQueryHandler. The WebSocketQueryHandler will use the WebSocketDataHandler to
// create establish a connection to the correct host, create subscription messages to be sent
// to the data provider, and handle incoming events accordingly.
type WebSocketDataHandler[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	HandleMessage(message []byte) (response providertypes.GetResponse[K, V], updateMessages []WebsocketEncodedMessage, err error)
	CreateMessages(ids []K) ([]WebsocketEncodedMessage, error)
	HeartBeatMessages() ([]WebsocketEncodedMessage, error)
	Copy() WebSocketDataHandler[K, V]
}
```

#### Determining K and V

First developers must determine the type of data that they want to fetch from the underlying data source. This can be any type that is supported by the oracle. For example, the simplest example is price data for a given currency pair (base / quote). The `K` type would be the currency pair and the `V` type would be the price data.

```golang
WebSocketDataHandler[types.ProviderTicker, *big.Float]
```

#### HandleMessage

HandleMessage is used to handle a message received from the data provider. Message parsing and response creation should be handled by this data handler. Given a message from the websocket the handler should either return a response or a set of update messages.

#### CreateMessages

CreateMessages is used to update the connection to the data provider. This can be used to subscribe to new events or unsubscribe from events.

#### HeartBeatMessages

HeartBeatMessages is used to construct a heartbeat messages to be sent to the data provider. As the provider is receiving messages from the data provider, it should store any relevant identification data that is required to construct the heartbeat messages.

#### Copy

Copy is used to create a copy of the data handler. This is useful if the data handler needs to be shared across multiple providers.

### WebSocketConnHandler

WebSocketConnHandler is an interface the encapsulates the functionality of a websocket connection to a data provider.

```golang
// WebSocketConnHandler is an interface the encapsulates the functionality of a websocket
// connection to a data provider. It provides the simple CRUD operations for a websocket
// connection. The connection handler is responsible for managing the connection to the
// data provider. This includes creating the connection, reading messages, writing messages,
// and closing the connection.
type WebSocketConnHandler interface {
	Read() ([]byte, error)
	Write(message []byte) error
	Close() error
	Dial() error
	Copy() WebSocketConnHandler
}
```

#### Read

`Read()` is used to read data from the data provider. This should block until data is received from the data provider.

#### Write

`Write()` is used to write data to the data provider. This should block until the data is sent to the data provider.

#### Close

`Close()` is used to close the connection to the data provider. Any resources associated with the connection should be cleaned up.

#### Dial

`Dial()` is used to establish a connection to the data provider. This should block until the connection is established.

#### Copy

`Copy()` is used to create a copy of the connection handler. This is useful if the connection handler needs to be shared across multiple providers.

