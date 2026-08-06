package app

import (
	"bytes"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/runmanifest"
)

// THE GATE (FINAL.md G1). Persisting the snapshot is not enough — a continuation must
// CONSUME it. This is the acceptance test the gate demanded: freeze a run, then change
// the configuration underneath it, and prove the continuation still launches what the run
// froze. Without it, `roster sync` (which exists to change deck values) could move a
// running deliberation onto a different model mid-flight.
func TestContinuationUsesTheFrozenSnapshotNotTheChangedConfig(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{
		{Agent: "claude-1", Adapter: "claude", Model: "claude-opus-5[1m]", Effort: "max", Speed: "deep"},
	}
	// The world moved on after the run was created: config now says a different model.
	discovered := []agents.Discovery{{
		Spec: agents.Spec{ID: "claude", Model: "claude-sonnet-9", Reasoning: "low", Speed: "fast"},
	}}

	got := applyRosterSnapshot(discovered, frozen, nil)
	if got[0].Spec.Model != "claude-opus-5[1m]" {
		t.Fatalf("continuation used %q; the run froze claude-opus-5[1m]", got[0].Spec.Model)
	}
	if got[0].Spec.Reasoning != "max" {
		t.Fatalf("continuation used effort %q; the run froze max", got[0].Spec.Reasoning)
	}
	if got[0].Spec.Speed != "deep" {
		t.Fatalf("continuation used speed %q; the run froze deep", got[0].Spec.Speed)
	}
}

// A run created before snapshots exist must keep working exactly as it did.
func TestApplyRosterSnapshotIsANoOpWithoutASnapshot(t *testing.T) {
	discovered := []agents.Discovery{{Spec: agents.Spec{ID: "claude", Model: "m"}}}
	got := applyRosterSnapshot(discovered, nil, nil)
	if got[0].Spec.Model != "m" {
		t.Fatalf("pre-snapshot run was altered: %q", got[0].Spec.Model)
	}
}

// An agent absent from the snapshot launches as configured, but never silently: a frozen
// run gaining an unfrozen participant is the same class of surprise the gate prevents.
func TestApplyRosterSnapshotWarnsAboutAgentsItDidNotFreeze(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{{Agent: "claude-1", Adapter: "claude", Model: "opus"}}
	discovered := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "changed"}},
		{Spec: agents.Spec{ID: "kimi", Model: "k3"}},
	}
	var warn bytes.Buffer
	got := applyRosterSnapshot(discovered, frozen, &warn)
	if got[1].Spec.Model != "k3" {
		t.Fatalf("unfrozen agent should launch as configured, got %q", got[1].Spec.Model)
	}
	if warn.Len() == 0 {
		t.Fatal("an unfrozen agent inside a frozen run must be reported")
	}
}

// `unknown` in the snapshot means "the launch passed no such flag". Pinning that literal
// back onto a spec would send the string "unknown" to the vendor CLI.
func TestApplyRosterSnapshotIgnoresUnknownCells(t *testing.T) {
	frozen := []runmanifest.RosterSnapshotEntry{{Agent: "codex-1", Adapter: "codex", Model: agents.Unknown, Effort: agents.Unknown}}
	discovered := []agents.Discovery{{Spec: agents.Spec{ID: "codex", Model: "gpt-5.6-sol", Reasoning: "xhigh"}}}
	got := applyRosterSnapshot(discovered, frozen, nil)
	if got[0].Spec.Model != "gpt-5.6-sol" || got[0].Spec.Reasoning != "xhigh" {
		t.Fatalf("unknown cells must not overwrite live config: %+v", got[0].Spec)
	}
}
