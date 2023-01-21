<!-- markdownlint-disable MD041 -->
<!-- markdownlint-disable MD013 -->

![Logo!](assets/nebula-logo.png)

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#wip)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://godoc.org/github.com/tessornetwork/nebula)
[![Go Report Card](https://goreportcard.com/badge/github.com/tessornetwork/nebula?style=flat-square)](https://goreportcard.com/report/github.com/tessornetwork/nebula)
[![Version](https://img.shields.io/github/tag/tessornetwork/nebula.svg?style=flat-square)](https://github.com/tessornetwork/nebula/releases/latest)
[![License: Apache-2.0](https://img.shields.io/github/license/tessornetwork/nebula.svg?style=flat-square)](https://github.com/tessornetwork/nebula/blob/main/LICENSE)
[![Lines Of Code](https://img.shields.io/tokei/lines/github/tessornetwork/nebula?style=flat-square)](https://github.com/tessornetwork/nebula)
[![GitHub Super-Linter](https://img.shields.io/github/workflow/status/tessornetwork/nebula/Lint?style=flat-square&label=Lint)](https://github.com/marketplace/actions/super-linter)

> A Golang implementation of the Nebula network, a decentralized universal capital
> facility in the Cosmos ecosystem.

Nebula is a Universal Capital Facility that can collateralize assets on one blockchain
towards borrowing assets on another blockchain. The platform specializes in
allowing staked assets from PoS blockchains to be used as collateral for borrowing
across blockchains. The platform uses a combination of algorithmically determined
interest rates based on market driven conditions. As a cross chain DeFi protocol,
Nebula will allow a multitude of decentralized debt products.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Releases](#releases)
  - [Release Compatibility Matrix](#release-compatibility-matrix)
- [Active Networks](#active-networks)
  - [Public](#public)
- [Install](#install)
  - [Swagger](#swagger)

## Releases

See [Release procedure](CONTRIBUTING.md#release-procedure) for more information about the release model.

Since `nebud v3.2` there is a new runtime dependency: `libwasmvm.x86_64.so v1.1.1` is required.
Building from source will automatically link the `libwasmvm.x86_64.so` created as a part of the build process (you must build on same host as you run the binary, or copy the `libwasmvm.x86_64.so` your lib directory).

### Release Compatibility Matrix

| Nebula Version | Mainnet | Experimental | Cosmos SDK |  IBC   |  Peggo  | Price Feeder |       Gravity Bridge       |
| :----------: | :-----: | :----------: | :--------: | :----: | :-----: | :----------: | :------------------------: |
|    v0.8.x    |    ✗    |      ✓       |  v0.45.x   | v2.0.x | v0.2.x  |    v0.1.x    |                            |
|    v1.x.x    |    ✓    |      ✗       |  v0.45.x   | v2.0.x | v0.2.x  |     N/A      | nebula/v1 module/v1.4.x-nebula |
|    v2.x.x    |    ✗    |      ✓       |  v0.45.x   | v2.3.x | v0.2.x  |    v0.2.x    |   nebula/v2 module/v1.4.x    |
|   v3.0-1.x   |    ✓    |      ✗       |  v0.46.x   | v5.0.x | v1.3.x+ |    v1.0.x    | nebula/v3 module/v1.5.x-nebula |
|    v3.2.x    |    ✓    |      ✗       |  v0.46.6+  | v5.0.x | v1.3.x+ |    v2.0.x    |   nebula/v3 v1.5.3-nebula-3    |

## Active Networks

### Public

- [nebula-1](networks/nebula-1) (mainnet)
- [canon-2](networks/canon-2) (testnet)

## Install

To install the `nebud` binary:

```shell
$ make install
```

### Swagger

- To enable it, modify the node config at `$NEBULA_HOME/config/app.toml` to `api.swagger` `true`
- Run the node normally `nebud start`
- Enter the swagger docs `http://localhost:1317/swagger/`

### Cosmovisor

> [Docs](https://github.com/cosmos/cosmos-sdk/tree/main/tools/cosmovisor)
> Note: `cosmovisor` only works for upgrades in the `nebud`, for off-chain processes updates like `peggo` or `price-feeder`, manual steps are required.

- `cosmovisor` is a small process manager for Cosmos SDK application binaries that monitors the governance module for incoming chain upgrade proposals. If it sees a proposal that gets approved, `cosmovisor` can automatically download
  the new binary, stop the current binary, switch from the old binary to the new one, and finally restart the node with the new binary.

- Install it with go

```shell
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
```

- For the usual use of `cosmovisor`, we recomend setting theses env variables

```shell
export DAEMON_NAME=nebud
export DAEMON_HOME={NODE_HOME}
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_PREUPGRADE_MAX_RETRIES=3
```
- If you didn't build binary from source in the machine, you have to download the respective `libwasmvm` into your machine.  
```bash
$ wget https://raw.githubusercontent.com/CosmWasm/wasmvm/v1.1.1/internal/api/libwasmvm.$(uname -m).so -O /lib/libwasmvm.$(uname -m).so
```
- To use `cosmovisor` for starting `nebud` process, instead of calling `nebud start`, use `cosmovisor run start [nebud flags]`
