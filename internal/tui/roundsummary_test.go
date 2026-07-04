package tui

import (
	"strings"
	"testing"

	"parley-deck-cli/internal/store"
)

func TestLatestRoundDigestDecodes(t *testing.T) {
	m := liveModel{events: []store.Event{
		{Type: "round.completed", Data: map[string]any{"round": "round-01"}},
		{Type: "round.digest", Data: map[string]any{
			"round":  "round-01",
			"digest": `{"idea":"demo","round":1,"total":2,"completed":2,"lines":[{"agent":"claude-1","position":"Accepts v2.","present":true}],"flag_block":1,"next":"drafting consensus"}`,
		}},
	}}
	dv, ok := m.latestRoundDigest()
	if !ok {
		t.Fatal("expected a digest")
	}
	if dv.Round != 1 || dv.Completed != 2 || dv.Next != "drafting consensus" {
		t.Fatalf("decoded = %+v", dv)
	}
}

func TestRenderRoundDigestBoundedAndLabelsHints(t *testing.T) {
	dv := digestView{Round: 2, Total: 3, Completed: 3, Next: "opening round-03", FlagBlock: 2, FlagCounter: 1}
	for i, a := range []string{"claude-1", "codex-1", "hermes-1"} {
		dv.Lines = append(dv.Lines, struct {
			Agent    string `json:"agent"`
			Position string `json:"position"`
			Fell     bool   `json:"fell_back"`
			Present  bool   `json:"present"`
		}{Agent: a, Position: "position " + string(rune('A'+i)), Present: true})
	}
	out := renderRoundDigest(dv, 80, 10)
	if !strings.Contains(out, "complete (3/3)") {
		t.Fatalf("missing completeness header: %q", out)
	}
	// Flags MUST be framed as mentions, never verdicts.
	if !strings.Contains(out, "mentions:") {
		t.Fatalf("flags not labeled as mentions: %q", out)
	}
	if !strings.Contains(out, "next: opening round-03") {
		t.Fatalf("missing next line: %q", out)
	}

	// A tiny row budget must not overflow: at most maxRows lines.
	small := renderRoundDigest(dv, 80, 4)
	if lines := strings.Count(small, "\n") + 1; lines > 4 {
		t.Fatalf("render exceeded row budget: %d lines", lines)
	}
}

func TestRenderRoundDigestZeroRows(t *testing.T) {
	if renderRoundDigest(digestView{Round: 1}, 80, 2) != "" {
		t.Fatal("maxRows<3 should render nothing")
	}
}
