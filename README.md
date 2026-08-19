# nomad-constant-rate-fabric

A small Go traffic-shaping experiment for Nomad's fixed-size, fixed-cadence cell stream.

`Run` emits **one 1200-byte cell per configured cell interval**. Earlier versions of this repository emitted an epoch's cells as a burst and therefore did not implement constant-rate timing; that behavior has been removed from the wall-clock path.

## Implemented

- fixed 1200-byte cells,
- cryptographically random filler source,
- one-cell-at-a-time wall-clock scheduler,
- explicit traffic-class accounting epochs,
- deterministic source/sink execution helpers for tests.

`EmitEpoch` deliberately skips delays and exists only for tests/batch processing. It must not be used as evidence of wire-level constant-rate behavior.

## Not implemented

There is no UDP transport, congestion response, NAT traversal, peer discovery, kernel scheduling control or packet-capture validation here. OS scheduling jitter can still change real packet timing. Wire-level claims belong in `nomad-testnet` and must be based on packet captures, not this package's planned trace.

```bash
go test -race ./...
go vet ./...
go run ./cmd/fabric-sim -epochs 1000 -cells 16 -epoch 100ms > trace.csv
```
