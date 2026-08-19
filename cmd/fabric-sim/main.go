package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Jtensetti/nomad-constant-rate-fabric/fabric"
)

func main() {
	ticks := flag.Int("ticks", 1000, "number of protocol ticks")
	cells := flag.Int("cells", 16, "cells per tick")
	flag.Parse()

	cfg := fabric.Config{Tick: 100 * time.Millisecond, CellsPerTick: *cells}
	trace, err := fabric.Trace(cfg, *ticks)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"tick", "bytes", "cells"})
	for i, n := range trace {
		_ = w.Write([]string{strconv.Itoa(i), strconv.Itoa(n), strconv.Itoa(*cells)})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "generated", len(trace), "constant-rate ticks")
}
