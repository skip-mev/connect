# Aggregating Oracle Prices On-Chain

## Overview

> **Definitions**
>
> * **Price**: The price of an asset at a given time.
> * **Delta**: The difference between the current price and the previous price. This can be positive or negative.
> * **Aggregation Strategy**: The method by which the deltas are aggregated across multiple oracle sources.
> 
> **Note**: The price aggregation strategy might vary based on the currency pair strategy utilized by the vote extension handler. For example, if the vote extension is implemented to submit a delta, the aggregation strategy must aggregate deltas. If the vote extension is implemented to submit a price, the aggregation strategy must aggregate prices. By default, the vote extension handler will submit deltas which are converted to real prices before being aggregated. This is the recommended approach.

This module provides the default implementation for a price aggregation strategy across multiple oracle sources on-chain. It is designed to be used as the aggregation function for the `PreBlock` handler. This implementation can be used as a reference for other implementations. For a reference as to how these strategies can be implemented, please reference the [aggregator](../../../aggregator/README.md) module.

## Implementation

The `PreBlock` handler is implemented in `preblock.go`. This handler is registered with the application in `app.go`. The handler is called by the ABCI application before the block is committed to the chain - in `PreBlock` during `FinalizeBlock`. 

This implementation will first associate the relative voting power of each validator given the price update that they have submitted through their vote extensions. After the voting power is calculated, the price information is aggregated across all validators by first sorting the prices in ascending order and then taking the stake weighted median of the prices.

For example, if there are 3 validators with the following prices:

```golang
Validator 1: 100
Validator 2: 200
Validator 3: 300
```

Assume the validators have the following voting power:

```golang
Validator 1: 10
Validator 2: 20
Validator 3: 20
```

The final aggregated price will be `200` which is the median of the sorted prices.

Another example, if there are 3 validators with the following prices:

```golang
Validator 1: 100
Validator 2: 200
Validator 3: 300
```

Assume the validators have the following voting power:

```golang
Validator 1: 10
Validator 2: 20
Validator 3: 30
```

The final aggregated price will be `250` which is the median of the sorted prices. Notice, that the price aggregation strategy selects the price update where the voting power is greater than or equal to 50%. In this case, the price update from `Validator 2` is selected.

As a final example, if there are 3 validators with the following prices:

```golang
Validator 1: 100
Validator 2: 200
Validator 3: 300
```

Assume the validators have the following voting power:

```golang
Validator 1: 10
Validator 2: 10
Validator 3: 100
```

The final aggregated price will be `300` which is the median of the sorted prices.
