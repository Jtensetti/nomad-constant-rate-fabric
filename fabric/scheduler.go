package fabric

import (
	"context"
	"errors"
	"time"
)

// Source provides the next payload cell. The scheduler never receives user
// activity state, which is an intentional Selection Firewall boundary.
type Source interface {
	NextCell(context.Context) (Cell, error)
}

// RandomSource produces utility/cover cells when no other protocol work is available.
type RandomSource struct{}

func (RandomSource) NextCell(context.Context) (Cell, error) {
	var c Cell
	return c, FillRandom(&c)
}

// Config fixes externally observable traffic shape.
type Config struct {
	Tick         time.Duration
	CellsPerTick int
}

func (c Config) Validate() error {
	if c.Tick <= 0 {
		return errors.New("tick must be positive")
	}
	if c.CellsPerTick <= 0 {
		return errors.New("cells per tick must be positive")
	}
	return nil
}

// Sink receives cells emitted at protocol-determined times.
type Sink interface {
	Send(context.Context, Cell) error
}

// Scheduler emits exactly CellsPerTick cells every Tick. It has deliberately
// no API for application demand or local selection state.
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

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.Tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for i := 0; i < s.cfg.CellsPerTick; i++ {
				c, err := s.source.NextCell(ctx)
				if err != nil {
					return err
				}
				if err := s.sink.Send(ctx, c); err != nil {
					return err
				}
			}
		}
	}
}

// Trace deterministically describes observable traffic without running a wall clock.
func Trace(cfg Config, ticks int) ([]int, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if ticks < 0 {
		return nil, errors.New("ticks must be non-negative")
	}
	out := make([]int, ticks)
	for i := range out {
		out[i] = cfg.CellsPerTick * CellSize
	}
	return out, nil
}
