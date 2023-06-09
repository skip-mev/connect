# slinky

<!-- markdownlint-disable MD013 -->
<!-- markdownlint-disable MD041 -->
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#wip)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://godoc.org/github.com/skip-mev/slinky)
[![Go Report Card](https://goreportcard.com/badge/github.com/skip-mev/slinky?style=flat-square)](https://goreportcard.com/report/github.com/skip-mev/slinky)
[![Version](https://img.shields.io/github/tag/skip-mev/slinky.svg?style=flat-square)](https://github.com/skip-mev/slinky/releases/latest)
[![License: Apache-2.0](https://img.shields.io/github/license/skip-mev/slinky.svg?style=flat-square)](https://github.com/skip-mev/slinky/blob/main/LICENSE)
[![Lines Of Code](https://img.shields.io/tokei/lines/github/skip-mev/slinky?style=flat-square)](https://github.com/skip-mev/slinky)

A general purpose price oracle leveraging ABCI++

## Install

```shell
$ go install github.com/skip-mev/slinky
```

## Status

Current, slinkly contains the following structure:

```text
.
├── config        <- Configuration that instructs the Oracle
├── oracle        <- The oracle implementation along with providers and types
│   ├── provider
│   └── types
└── pkg           <- Package types and utilities
    └── sync
```

The following components have been implemented:

* `oracle.go`: The main oracle implementation that is responsible for fetching prices
  from providers given a set of assets and a set of providers.
* `provider.go`: The provider interface that is used to fetch prices from a given
  provider along with a mock implementation.
* `client/local.go`: A local client that is used to communicate with an ABCI++
  application.
* `client/grpc.go`: A gRPC client that is used to communicate with an ABCI++
  application.
* `abci/vote_ext.go`: A vote extension handler implementation.

The following components are still in progress or need to be implemented completely:

* gRPC server & corresponding CLI command to launch it
* x/oracle module
* various providers in `oracle/provider/` such Kraken, Binance, etc...
* config type definition in `config/`
* tests!!!
