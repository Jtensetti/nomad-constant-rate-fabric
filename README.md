# nomad-constant-rate-fabric

Research implementation of Nomad's **activity-independent traffic shaper**.

The package emits a protocol-defined number of fixed-size cells per time slice. Local user activity is intentionally absent from the scheduler API. This makes the principal invariant testable:

> The externally observable traffic shape must not depend on whether the local user is idle, reading, searching, or reconstructing content.

## What is implemented

- Fixed 1200-byte cells.
- Deterministic traffic-trace model.
- Wall-clock scheduler for local experiments.
- Cryptographically random utility cells.
- Tests asserting invariant traffic shape.
- CSV trace generator.

## What this is not

This repository is a research component, not an audited anonymity transport and not a public-network deployment tool. It deliberately omits peer discovery, NAT traversal, Internet routing and deployment automation.

## Build and test

```bash
go test ./...
go test -race ./...
go vet ./...
go run ./cmd/fabric-sim -ticks 1000 -cells 16 > trace.csv
```

## Security invariant

`Scheduler` has no reference to local selection state. If future code adds such a dependency, treat that as a security regression.

See the architecture repository for the complete threat model and cross-component invariants.
