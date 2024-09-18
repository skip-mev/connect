# Oracle Vote Extensions

> **This readme provides an overview of how the oracle, base app, and vote extension interact.**

## Overview

Each validator in the network running the connect module either runs an oracle in process or out of process. In either 
case, the oracle is responsible for fetching data offchain that will include the validator's vote extension, broadcast to the network, included in a block, and subsequently committed on-chain.

The process of fetching data offchain and broadcasting it to the network is handled by the vote extension handlers. Note, when a validator broadcasts their vote extensions, they will only be available to other validators in the network at the next height. This means that oracle prices
are always one height behind.

Additionally, each validator has a local view of the network's set of vote extensions, meaning there
is not a canonical set of vote extensions. As such, we must utilize the next proposer's local view of
the network's vote extensions as canonical in order to maintain determinism across the network.

## Extend Vote Extension

The extend vote extension handler has access to the oracle service via a remote or local gRPC client. Extend vote is responsible for the following:

1. Fetching the data from the oracle service within the specified timeout period.
2. Creating a vote extension with the data fetched from the oracle service.
3. Broadcasting the vote extension to the network.

> Note: In the case where the oracle service is unavailable, returns a bad response, or times out, a nil vote extension will be broadcast to the network. We do not want to halt the chain because of an oracle failure.

## Verify Vote Extension

The verify vote extension handler acknowledges and verifies the vote extensions currently in transit across the network. The verify vote extension handler is responsible for the following:

1. Verifying the vote extension is valid. If the vote extension is empty, the vote extension is considered valid.
2. Verifying the vote extension is not expired. If the vote extension is expired, the vote extension is considered invalid.
3. Verifying that the prices provided in the vote extension are valid. If the prices are invalid, the vote extension is considered invalid.
