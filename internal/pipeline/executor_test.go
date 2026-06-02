package pipeline

import (
	"os"
	"testing"
	"time"
)

func testDriver(t *testing.T, autonomy Autonomy) (Driver, *PipelineRun) {
	t.Helper()
	m, err := Parse([]byte(validLinearYAML))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	m.Autonomy = autonomy
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	run := NewPipelineRun(m, now)
	d := Driver{DeckDir: t.TempDir(), Manifest: m, Now: now}
	return d, &run
}

// allDone reports every block complete; noneDone reports none complete.
func allDone(Block) (bool, error)  { return true, nil }
func noneDone(Block) (bool, error) { return false, nil }

func TestAdvanceRunBlockWhenIncomplete(t *testing.T) {
	d, run := testDriver(t, AutonomySupervised)
	res, err := d.Advance(run, noneDone)
	if err != nil {
		t.Fatalf("Advance: %v", err)
	}
	if res.Action != ActionRunBlock {
		t.Fatalf("action = %q, want run-block", res.Action)
	}
	if run.Status != StatusRunning {
		t.Fatalf("status = %q", run.Status)
	}
}

func TestAdvanceSupervisedOpensGate(t *testing.T) {
	d, run := testDriver(t, AutonomySupervised)
	res, err := d.Advance(run, allDone)
	if err != nil {
		t.Fatalf("Advance: %v", err)
	}
	if res.Action != ActionAwaitGate {
		t.Fatalf("action = %q, want await-gate", res.Action)
	}
	if run.Status != StatusBlockedOnGate || run.PendingGate != EdgeID("business-spec", "technical-spec") {
		t.Fatalf("run = %+v", run)
	}
	// Gate file should exist and be open.
	g, ok, _ := LoadGate(d.DeckDir, run.PipelineSlug, run.PendingGate)
	if !ok || g.Status != GateOpen {
		t.Fatalf("gate ok=%v status=%q", ok, g.Status)
	}
}

func TestAdvanceSeedsNextAfterGateApproved(t *testing.T) {
	d, run := testDriver(t, AutonomySupervised)
	// First advance opens the gate.
	if _, err := d.Advance(run, allDone); err != nil {
		t.Fatalf("advance 1: %v", err)
	}
	// Human approves it.
	g, _, _ := LoadGate(d.DeckDir, run.PipelineSlug, run.PendingGate)
	g.Resolve(true, "tomas", d.Now)
	if err := SaveGate(d.DeckDir, g); err != nil {
		t.Fatalf("save gate: %v", err)
	}
	// Second advance seeds the next block.
	res, err := d.Advance(run, allDone)
	if err != nil {
		t.Fatalf("advance 2: %v", err)
	}
	if res.Action != ActionSeededNext {
		t.Fatalf("action = %q, want seeded-next", res.Action)
	}
	if run.CurrentBlock != "technical-spec" {
		t.Fatalf("current block = %q", run.CurrentBlock)
	}
	if len(run.CompletedBlocks) != 1 || run.CompletedBlocks[0] != "business-spec" {
		t.Fatalf("completed = %v", run.CompletedBlocks)
	}
	seeded := BlockWorkspace(d.DeckDir, run.PipelineSlug, "technical-spec") + "/00-prompt.md"
	if _, err := os.Stat(seeded); err != nil {
		t.Fatalf("seeded prompt missing: %v", err)
	}
}

func TestAdvanceAutoLeftAutoApprovesLowRisk(t *testing.T) {
	d, run := testDriver(t, AutonomyAutoLeft)
	// blocks have no risk set -> treated as auto-approvable under auto-left.
	res, err := d.Advance(run, allDone)
	if err != nil {
		t.Fatalf("Advance: %v", err)
	}
	if res.Action != ActionSeededNext {
		t.Fatalf("action = %q, want seeded-next (auto-approved)", res.Action)
	}
	if run.CurrentBlock != "technical-spec" {
		t.Fatalf("current = %q", run.CurrentBlock)
	}
}

func TestAdvanceCompletesAtLastBlock(t *testing.T) {
	d, run := testDriver(t, AutonomyAutoLeft)
	// Walk all boundaries (auto-approved) until done.
	for i := 0; i < 10; i++ {
		res, err := d.Advance(run, allDone)
		if err != nil {
			t.Fatalf("advance: %v", err)
		}
		if res.Action == ActionDone {
			if run.Status != StatusCompleted {
				t.Fatalf("status = %q, want completed", run.Status)
			}
			if len(run.CompletedBlocks) != 3 {
				t.Fatalf("completed = %v, want all 3", run.CompletedBlocks)
			}
			return
		}
	}
	t.Fatal("pipeline never reached done")
}

func TestAdvanceStopsOnRejectedGate(t *testing.T) {
	d, run := testDriver(t, AutonomySupervised)
	if _, err := d.Advance(run, allDone); err != nil {
		t.Fatalf("advance 1: %v", err)
	}
	g, _, _ := LoadGate(d.DeckDir, run.PipelineSlug, run.PendingGate)
	g.Resolve(false, "tomas", d.Now)
	_ = SaveGate(d.DeckDir, g)
	res, err := d.Advance(run, allDone)
	if err != nil {
		t.Fatalf("advance 2: %v", err)
	}
	if res.Action != ActionRejected || run.Status != StatusFailed {
		t.Fatalf("res=%+v status=%q", res, run.Status)
	}
}
