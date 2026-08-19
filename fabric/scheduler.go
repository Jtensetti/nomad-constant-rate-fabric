package fabric

import (
	"context"
	"errors"
	"time"
)

// Source supplies one complete fixed-size cell. Choosing what protocol work
// fills a cell is intentionally separate from the wall-clock scheduler.
type Source interface {
	NextCell(context.Context) (Cell, error)
}

type RandomSource struct{}

func (RandomSource) NextCell(context.Context) (Cell, error) { return RandomCell() }

// Config describes a traffic class. Epoch is the accounting window; cells are
// evenly spaced across that window by Run rather than emitted as an epoch burst.
type Config struct {
	Epoch         time.Duration
	CellsPerEpoch int
}

func (c Config) Validate() error {
	if c.Epoch <= 0 {
		return errors.New("epoch must be positive")
	}
	if c.CellsPerEpoch <= 0 {
		return errors.New("cells per epoch must be positive")
	}
	if c.CellInterval() <= 0 {
		return errors.New("epoch is too short for the configured cell count")
	}
	if c.Epoch%time.Duration(c.CellsPerEpoch) != 0 {
		return errors.New("epoch must divide exactly into equal cell intervals")
	}
	return nil
}

// CellInterval is the target spacing between externally visible cells.
func (c Config) CellInterval() time.Duration {
	if c.CellsPerEpoch <= 0 {
		return 0
	}
	return c.Epoch / time.Duration(c.CellsPerEpoch)
}

type Sink interface {
	Send(context.Context, Cell) error
}

type Scheduler struct {
	cfg    Config
	source Source
	sink   Sink
}

func NewScheduler(cfg Config, source Source, sink Sink) (*Scheduler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if source == nil || sink == nil {
		return nil, errors.New("source and sink are required")
	}
	return &Scheduler{cfg: cfg, source: source, sink: sink}, nil
}

// EmitOne executes one scheduled emission. It exists so the actual source/sink
// path can be exercised without wall-clock sleeps in unit tests.
func (s *Scheduler) EmitOne(ctx context.Context) error {
	cell, err := s.source.NextCell(ctx)
	if err != nil {
		return err
	}
	return s.sink.Send(ctx, cell)
}

// EmitEpoch emits one epoch's worth of cells without delays. It is a test and
// batch-processing helper; Run is the method that enforces wall-clock spacing.
func (s *Scheduler) EmitEpoch(ctx context.Context) error {
	for i := 0; i < s.cfg.CellsPerEpoch; i++ {
		if err := s.EmitOne(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Run emits one cell per CellInterval. No application/read state is consulted.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.CellInterval())
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.EmitOne(ctx); err != nil {
				return err
			}
		}
	}
}

// EpochTrace returns the planned byte count per accounting epoch. It is a
// planning helper, not a packet capture or timing measurement.
func EpochTrace(cfg Config, epochs int) ([]int, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if epochs < 0 {
		return nil, errors.New("epochs must be non-negative")
	}
	out := make([]int, epochs)
	for i := range out {
		out[i] = cfg.CellsPerEpoch * CellSize
	}
	return out, nil
}
