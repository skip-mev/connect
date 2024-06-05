# Oracle Client

## Overview

The oracle client can be run concurrently with the oracle server. It is meant to be a useful tool for interacting with the oracle server and ensuring the price feed is working as expected.

## Usage

The oracle client can be run with the following command:

```bash
make run-oracle-client
```

However, before running the oracle client, you must first run the oracle server. To start the oracle server run the following command:

```bash
make run-oracle-server
```

One additionally useful tool is to run the prometheus server and check the metrics that are being collected. To start the prometheus server run the following command:

```bash
make run-prom-client
```

After starting the prometheus server, you can view the metrics by navigating to `localhost:9090` in your browser.
