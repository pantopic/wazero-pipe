# Wazero Pipe

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK providing pipes between WASI modules.

## Host Module

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-pipe/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-pipe/host)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-pipe/host)](https://goreportcard.com/report/github.com/pantopic/wazero-pipe/host)
[![Go Coverage](https://github.com/pantopic/wazero-pipe/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-pipe/host/coverage.html)

First register the host module with the runtime

```go
import (
    "github.com/tetratelabs/wazero"
    "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

    "github.com/pantopic/wazero-pipe/host"
)

func main() {
    ctx := context.Background()
    r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
    wasi_snapshot_preview1.MustInstantiate(ctx, r)

    module := wazero_pipe.New()
    module.Register(ctx, r)

    // ...
}
```

## Guest SDK (Go)

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-pipe/sdk-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-pipe/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-pipe/sdk-go)](https://goreportcard.com/report/github.com/pantopic/wazero-pipe/sdk-go)

Then you can import the guest SDK into your WASI module to send messages from one WASI module to another.

```go
package main

import (
    "github.com/pantopic/wazero-pipe/sdk-go"
)

var p *pipe.Pipe[uint64]

func main() {
    p = pipe.New[uint64]()
}

//export test
func test() {
    p.Send(12345)
}
```

This is useful for emitting [watch events](https://etcd.io/docs/v3.3/learning/api/#watch-api) from WASM state machines.

## Roadmap

This project is in alpha. Breaking API changes should be expected until Beta.

- `v0.0.x` - Alpha
  - [ ] Stabilize API
- `v0.x.x` - Beta
  - [ ] Finalize API
  - [ ] Test in production
- `v1.x.x` - General Availability
  - [ ] Proven long term stability in production
