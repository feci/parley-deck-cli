package driver

import (
	"context"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestDriverActionIdentitySurvivesNewDriverWithoutDuplicateCharge(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	runner := &fakeRunner{writeOnRun: func(n int) { writeAll(t, ideaDir, n, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, runner)
	d.cfg.MaxDriverSteps = 3
	if action, _, e := d.Advance(context.Background()); e != nil || action != ActionPromoted {
		t.Fatalf("first advance %s %v", action, e)
	}
	step, e := budget.LoadStepBinding(context.Background(), d.cfg.Root, d.cfg.IdeaSlug)
	if e != nil || step == nil {
		t.Fatalf("missing step binding: %v", e)
	}
	cycle, e := budget.LoadCycleBinding(context.Background(), d.cfg.Root, d.cfg.IdeaSlug, budget.CrossReview)
	if e != nil || cycle == nil {
		t.Fatalf("missing cycle binding: %v", e)
	}
	originals := map[string]budget.ActionReceipt{}
	for _, s := range []budget.Store{step.Store, cycle.Store} {
		rows, e := s.InspectActions(context.Background())
		if e != nil || len(rows) != 1 {
			t.Fatalf("expected one group charge: %v", e)
		}
		r := rows[0]
		if r.Identity == nil || r.Identity.Input.Operation != "driver-advance" || r.Identity.Input.RunID != filepath.Base(runDir) || r.Identity.Input.Basis != "runtime-input" {
			t.Fatalf("driver input not retained: %+v", r)
		}
		originals[s.Scope] = r
	}
	// The canonical completed round lets the rebuilt driver advance its gate
	// without rerunning the finished group. The durable receipt remains original.
	resumed := newTestDriver(ideaDir, runDir, parts, 1, true, runner)
	resumed.cfg.MaxDriverSteps = 3
	if _, _, e := resumed.Advance(context.Background()); e != nil {
		t.Fatal(e)
	}
	if len(runner.calls) != 1 {
		t.Fatal("resume reran a completed group")
	}
	for _, s := range []budget.Store{step.Store, cycle.Store} {
		rows, e := s.InspectActions(context.Background())
		if e != nil || len(rows) != 1 {
			t.Fatal("resume duplicated accounting")
		}
		old := originals[s.Scope]
		r, e := s.ReplayAction(context.Background(), old.EntryKey, old.IdentitySHA256)
		if e != nil || r.Permission != "none" {
			t.Fatal("replay changed accounting authority")
		}
	}
}
