# Clients

## Overview

This directory contains all of the clients that are supported for the general purpose oracle. Each client is responsible for fetching data from a specific type of provider (price, random number, etc.) and returning a standardized response. The client is utilized by the a validator's application (Cosmos SDK block chain) to fetch data from the out of process oracle service before processing the data and including it in their vote extensions.

## Clients

* **[Prices](./oracle/)** - This client supports fetching prices from a oracle that is aggregating price data.
