package runaction

import (
	"strings"
	"testing"
)

func TestCommandUsesRoundAndFallbackTargets(t *testing.T) {
	draft := Command(NextAction{
		Kind:     KindDraftConsensus,
		IdeaSlug: "sample",
		Round:    "round-02",
	}, "run-1", "")
	if draft != "parley consensus draft --round 2 sample" {
		t.Fatalf("draft command=%q", draft)
	}
	if strings.Contains(draft, "codex") {
		t.Fatalf("draft command hardcodes agent: %q", draft)
	}

	finalize := Command(NextAction{Kind: KindFinalize}, "run-1", "sample")
	if finalize != "parley consensus finalize sample" {
		t.Fatalf("finalize command=%q", finalize)
	}
	if strings.Contains(finalize, "codex") {
		t.Fatalf("finalize command hardcodes agent: %q", finalize)
	}
}

func TestCommandLeavesRetryUnmapped(t *testing.T) {
	if command := Command(NextAction{Kind: KindRetryAgent, RunID: "run-1", IdeaSlug: "sample"}, "", ""); command != "" {
		t.Fatalf("retry command=%q, want empty", command)
	}
}

func TestRoundNumberRejectsZeroRound(t *testing.T) {
	if round := RoundNumber("round-00"); round != "" {
		t.Fatalf("round=%q, want empty", round)
	}
}
