package fabric

import "testing"

func TestTraceIndependentOfActivity(t *testing.T) {
	cfg := Config{Tick: 1, CellsPerTick: 16}
	a, err := Trace(cfg, 10_000)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Trace(cfg, 10_000)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != len(b) {
		t.Fatal("trace length differs")
	}
	for i := range a {
		if a[i] != b[i] || a[i] != 16*CellSize {
			t.Fatalf("unexpected observable trace at %d: %d %d", i, a[i], b[i])
		}
	}
}

func TestInvalidConfig(t *testing.T) {
	if _, err := Trace(Config{}, 1); err == nil {
		t.Fatal("expected validation error")
	}
}
