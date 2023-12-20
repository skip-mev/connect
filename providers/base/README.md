# Base Provider

## Overview

The Base Provider is a special provider that is used to provide the base functionality for all other providers. It is not intended to be used directly, but instead serves as a base for other providers to extend.

The base provider is responsible for the following:

* Implementing the oracle `Provider` interface directly just by inheriting the base provider.
* Having the provider data avaiable in constant time when the oracle requests it.

Each base provider implementation will be run in a separate goroutine by the main oracle process. This allows the provider to fetch data from the underlying data source asynchronusly. The base provider will then store the data in a thread safe map. The main oracle service utilizing this provider can determine if the data is stale or not based on the last time the data was fetched.

## Implementation

In order to implement a provider, you must inherit the base provider and implement the `APIDataHandler` interface. The `APIDataHandler` interface is responsible for fetching data from the underlying data source. The base provider will then take care of the rest.

```golang
// APIDataHandler interface defines the methods that need to be implemented by the extender.
type APIDataHandler[K comparable, V any] interface {
	// Get is used to fetch data from the API.
	Get(ctx context.Context) (map[K]V, error)
}
```

The `APIDataHandler` interface is purposefully built with generics in mind. This allows the provider to fetch data of any type from the underlying data source.

## APIDataHandler Usage

### Determining K and V

> **Currently the oracle only supports `*big.Int` as the `V` type and `oracletypes.CurrencyPair` as the `K` type.** This will change in the future once generics are supported on the chain side.

First developers must determine the type of data that they want to fetch from the underlying data source. This can be any type that is supported by the oracle. For example, the simplest example is price data for a given currency pair (base / quote). The `K` type would be the currency pair and the `V` type would be the price data.

```golang
APIDataHandler[oracletypes.CurrencyPair, *big.Int]
```

### Get

This method is used to fetch data from the underlying data source. Following the price example above, the `Get` method would be used to fetch the price data for a given currency pair.
