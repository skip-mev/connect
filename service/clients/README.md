# Clients

## Overview

This directory contains all clients that are supported for the general purpose oracle. Each client is responsible for fetching data from a specific type of provider (price, random number, etc.) and returning a standardized response. The client is utilized by the validator's application (Cosmos SDK blockchain) to fetch data from the out of process oracle service before processing the data and including it in their vote extensions.

## Clients

* **[Prices](./oracle/)** - This client supports fetching prices from the oracle that is aggregating price data.
