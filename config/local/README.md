# Local Oracle Configuration

## Overview

This directory contains the [configuration file](./oracle.json) for the local oracle instance as well as [market map file](./market.json) that contains all of the markets that are supported by the oracle. To update the set of provider's utilized, update the `LocalOracleConfig` in `generate.go`. To update the set of markets supported, update the `ProvidersToMarkets` in `generate.go`. To add a custom conversion path, update the `TickerPaths` in `generate.go`.

```bash
make update-local-config
```

## Considerations

Note that not every provider supports every currency pair. If you configure pairs that are not supported, some providers may stop returning responses. As such, please read over the documentation pertaining to each provider before adding new price feeds.


