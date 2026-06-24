package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/store"
)

// LE-5: loopBudgetBreach reports a breach only when a non-zero ceiling is exceeded.
func TestLoopBudgetBreach(t *testing.T) {
	mk := func(steps int, wall time.Duration, cost float64) *Driver {
		return New(Config{IdeaSlug: "demo", RunDir: t.TempDir(), Events: store.New(t.TempDir()),
			MaxDriverSteps: steps, MaxWallClock: wall, MaxCostUSD: cost}, &fakeRunner{})
	}
	d := mk(3, 0, 0)
	if d.loopBudgetBreach(2, time.Now()) != "" {
		t.Fatal("2/3 steps must not breach")
	}
	if d.loopBudgetBreach(3, time.Now()) == "" {
		t.Fatal("3/3 steps must breach")
	}

	dw := mk(0, 50*time.Millisecond, 0)
	if dw.loopBudgetBreach(0, time.Now().Add(-time.Second)) == "" {
		t.Fatal("a start 1s ago against a 50ms budget must breach")
	}
	if dw.loopBudgetBreach(0, time.Now()) != "" {
		t.Fatal("a fresh start must not breach")
	}

	// 0 ceilings everywhere = unlimited (backward-compatible default).
	du := mk(0, 0, 0)
	if du.loopBudgetBreach(1_000_000, time.Now().Add(-24*time.Hour)) != "" {
		t.Fatal("0 ceilings must never breach (unlimited)")
	}
}

func TestIsProgressAction(t *testing.T) {
	for _, a := range []Action{ActionPromoted, ActionConsensusDrafted, ActionFinalized, ActionImplemented, ActionReviewOpened, ActionReviewDrafted, ActionFixup} {
		if !isProgressAction(a) {
			t.Fatalf("%s should be a progress action", a)
		}
	}
	for _, a := range []Action{ActionAwait, ActionComplete, ActionConsensus, ActionSurfaceOnly, ActionEscalated} {
		if isProgressAction(a) {
			t.Fatalf("%s should not be a progress action", a)
		}
	}
}

func TestEmitLoopBudgetEvent(t *testing.T) {
	runDir := t.TempDir()
	d := New(Config{IdeaSlug: "demo", RunDir: runDir, Events: store.New(runDir), MaxDriverSteps: 10}, &fakeRunner{})
	d.emitLoopBudget(4, time.Now())
	evs, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == "loop.budget" {
			return
		}
	}
	t.Fatal("expected a loop.budget event")
}

// LE-6: loopCostUSD sums cost_usd across agent.usage events (best-effort).
func TestLoopCostUSDSumsAgentUsage(t *testing.T) {
	runDir := t.TempDir()
	st := store.New(runDir)
	_ = st.Append(store.Event{Time: time.Now(), Type: "agent.usage", Data: map[string]any{"cost_usd": 1.5}})
	_ = st.Append(store.Event{Time: time.Now(), Type: "agent.usage", Data: map[string]any{"cost_usd": 2.25}})
	_ = st.Append(store.Event{Time: time.Now(), Type: "other", Data: map[string]any{"cost_usd": 99.0}})
	d := New(Config{IdeaSlug: "demo", RunDir: runDir, Events: st}, &fakeRunner{})
	if got := d.loopCostUSD(); got < 3.74 || got > 3.76 {
		t.Fatalf("loopCostUSD = %v, want ~3.75 (only agent.usage events count)", got)
	}
}

// LE-5: a budget breach writes a durable blocking inbox note and halts (not complete).
func TestEscalateLoopBudgetWritesInbox(t *testing.T) {
	deck := t.TempDir()
	runDir := filepath.Join(deck, "runs", "r1")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d := New(Config{IdeaSlug: "demo", RunDir: runDir, Events: store.New(runDir)}, &fakeRunner{})
	if err := d.escalateLoopBudget(Cursor{Phase: PhaseReview}, "driver-step budget exhausted (3/3 steps)"); err != nil {
		t.Fatal(err)
	}
	note := filepath.Join(deck, "inbox", "claude-to-user_demo_loop-budget.md")
	if _, err := os.Stat(note); err != nil {
		t.Fatalf("expected a loop-budget inbox note at %s: %v", note, err)
	}
}

// F-T2-2: the loop.budget event reports observed cost even when cost is unlimited
// (MaxCostUSD == 0) — unlimited disables enforcement, not observability.
func TestEmitLoopBudgetReportsCostWhenUnlimited(t *testing.T) {
	runDir := t.TempDir()
	st := store.New(runDir)
	_ = st.Append(store.Event{Time: time.Now(), Type: "agent.usage", Data: map[string]any{"cost_usd": 4.0}})
	d := New(Config{IdeaSlug: "demo", RunDir: runDir, Events: st, MaxCostUSD: 0}, &fakeRunner{})
	d.emitLoopBudget(1, time.Now())
	evs, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type != "loop.budget" {
			continue
		}
		c, ok := e.Data["cost_usd"].(float64)
		if !ok || c < 3.9 || c > 4.1 {
			t.Fatalf("loop.budget cost_usd = %v, want ~4.0 even when MaxCostUSD=0", e.Data["cost_usd"])
		}
		return
	}
	t.Fatal("expected a loop.budget event")
}

// F-T2-3: a Run with an already-elapsed wall-clock budget escalates via the inbox note
// on the first pre-Advance check and never reaches Complete.
func TestRunEscalatesOnLoopBudget(t *testing.T) {
	deck := t.TempDir()
	runDir := filepath.Join(deck, "runs", "r1")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d := New(Config{IdeaSlug: "demo", RunDir: runDir, Events: store.New(runDir), MaxWallClock: time.Nanosecond}, &fakeRunner{})
	if err := d.Run(context.Background()); err != nil {
		t.Fatalf("Run should halt cleanly on a budget breach, got err=%v", err)
	}
	note := filepath.Join(deck, "inbox", "claude-to-user_demo_loop-budget.md")
	if _, err := os.Stat(note); err != nil {
		t.Fatalf("Run must write the loop-budget inbox note on breach: %v", err)
	}
}
