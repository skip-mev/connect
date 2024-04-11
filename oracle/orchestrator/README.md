# Orchestrator

The Orchestrator is the central component of the oracle. It is responsible for:

* Spinning up the necessary services (e.g. the sidecar, the price fetchers, etc.).
* Managing the lifecycle of the providers.
* Determining the set of markets that need to be fetched and updating the providers accordingly.

## Configuration

At a high level the orchestrator is configured with a `oracle.json` file that contains all of the providers that need to be instantiated. To read more about the configuration of `oracle.json`, please refer to the [oracle configuration documentation](./../config/README.md).

Each provider is instantiate using the `PriceAPIQueryHandlerFactory`, `PriceWebSocketQueryHandlerFactory`, and `MarketMapFactory` factory functions. Think of these as the constructors for the providers. 

* `PriceAPIQueryHandlerFactory` - This is used to create the API query handler for the provider - which is then passed into a base provider.
* `PriceWebSocketQueryHandlerFactory` - This is used to create the WebSocket query handler for the provider - which is then passed into a base provider.
* `MarketMapFactory` - This is used to create the market map provider.

## Lifecycle

The orchestrator can be initialized with an option of `WithMarketMap` which allows each provider to be instantiated with a predetermined set of markets. If this option is not provided, the orchestrator will fetch the markets from the market map provider. **Both options can be set.**

The orchestrator will then start each provider in a separate goroutine. Additionally, if the orchestrator has a market map provider, it will start a goroutine that will periodically fetch the markets from the market map provider and update the providers accordingly.

All of the providers are running concurrently and will do so until the main context is canceled (what is passed into `Start`). If the orchestrator is canceled, it will cancel all of the providers and wait for them to finish before returning.

