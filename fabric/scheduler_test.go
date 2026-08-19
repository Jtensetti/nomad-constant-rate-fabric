package fabric

import (
	"context"
	"errors"
	"testing"
	"time"
)

type byteSource byte

func (s byteSource) NextCell(context.Context) (Cell, error) {
	var c Cell
	for i := range c {
		c[i] = byte(s)
	}
	return c, nil
}

type recordingSink struct{ cells []Cell }

func (s *recordingSink) Send(_ context.Context, c Cell) error {
	s.cells = append(s.cells, c)
	return nil
}

func TestCellInterval(t *testing.T) {
	cfg := Config{Epoch: 100 * time.Millisecond, CellsPerEpoch: 16}
	if got, want := cfg.CellInterval(), 6250*time.Microsecond; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestDifferentPayloadSourcesHaveSameEpochShape(t *testing.T) {
	cfg := Config{Epoch: time.Second, CellsPerEpoch: 16}
	left := &recordingSink{}
	right := &recordingSink{}
	a, err := NewScheduler(cfg, byteSource(0x11), left)
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewScheduler(cfg, byteSource(0xee), right)
	if err != nil {
		t.Fatal(err)
	}
	for epoch := 0; epoch < 32; epoch++ {
		if err := a.EmitEpoch(context.Background()); err != nil {
			t.Fatal(err)
		}
		if err := b.EmitEpoch(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if len(left.cells) != 32*cfg.CellsPerEpoch || len(right.cells) != len(left.cells) {
		t.Fatalf("unexpected cell counts: %d and %d", len(left.cells), len(right.cells))
	}
	if left.cells[0] == right.cells[0] {
		t.Fatal("test sources unexpectedly produced identical payloads")
	}
}

type failingSource struct{}

func (failingSource) NextCell(context.Context) (Cell, error) {
	return Cell{}, errors.New("source failure")
}

func TestSourceFailureStopsEmission(t *testing.T) {
	sink := &recordingSink{}
	s, err := NewScheduler(Config{Epoch: time.Second, CellsPerEpoch: 2}, failingSource{}, sink)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.EmitOne(context.Background()); err == nil {
		t.Fatal("expected source error")
	}
	if len(sink.cells) != 0 {
		t.Fatalf("emitted %d cells after source failure", len(sink.cells))
	}
}

func TestEpochTraceMatchesConfiguredShape(t *testing.T) {
	cfg := Config{Epoch: time.Second, CellsPerEpoch: 16}
	trace, err := EpochTrace(cfg, 10)
	if err != nil {
		t.Fatal(err)
	}
	for i, bytes := range trace {
		if bytes != 16*CellSize {
			t.Fatalf("epoch %d: got %d bytes", i, bytes)
		}
	}
}

func TestInvalidConfig(t *testing.T) {
	cases := []Config{
		{},
		{Epoch: time.Second},
		{CellsPerEpoch: 1},
		{Epoch: time.Nanosecond, CellsPerEpoch: 2},
		{Epoch: 10 * time.Nanosecond, CellsPerEpoch: 3},
	}
	for _, cfg := range cases {
		if err := cfg.Validate(); err == nil {
			t.Fatalf("expected validation error for %#v", cfg)
		}
	}
}
