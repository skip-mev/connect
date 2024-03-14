# Oracle PreBlock Handler

## Overview

Writing vote extensions to state is possible using the SDK's `PreBlocker` - which allows you to modify the state before the block is executed
and committed. Since the vote extensions not directly accessible from the `PreBlocker`, we inject the vote extensions in `PrepareProposal` and verify them in `ProcessProposal` before a block is accepted by the network. 

The `PreBlockHandler` assumes that the vote extensions are already verified by validators in the network and are ready to be aggregated. A bad vote extension included in a proposal implies that the 
network has accepted a bad proposal.

## Usage

To use the preblock handler, you need to initialize the preblock handler in your `app.go` file. By default, we encourage users to use the aggregation function defined in `abci/preblock/math` to aggregate the votes. This will aggregate all prices and calculate a stake-weighted median for each supported asset. 

The `PreBlockHandler` currently only supports assets that are initialized in the oracle keeper. However, allowing any type of asset can be supported with a small modification to `WritePrices` (TBD whether we will support this).
