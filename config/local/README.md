# Local Oracle Configuration

## Overview

This directory contains the [configuration file](./oracle.json) for the local oracle instance. To update the configuration, edit the `LocalConfig` in [generate_toml.go](./generate_toml.go) and run the following command:

```bash
make update-local-config
```

## Considerations

Note that not every provider supports every currency pair. If you configure pairs that are not supported, some providers may stop returning responses. As such, please read over the documentation pertaining to each provider before adding new price feeds.


