package pipeline

import (
	"testing"
	"time"
)

// Blocks are listed in REVERSE topological order to prove advanceDAG selects by
// dependency readiness, not list order (the codex MAJOR finding).
const dagReverseYAML = `
schema_version: 1
idea_slug: dag-order
transport: local-dir
execution: dag
autonomy: auto-left
participants: [codex, claude]
blocks:
  - {id: integrate, kind: deliberation}
  - {id: api, kind: deliberation}
  - {id: spec, kind: deliberation}
edges:
  - {from: spec, to: api}
  - {from: api, to: integrate}
`

func TestAdvanceDAGRespectsTopologicalOrder(t *testing.T) {
	m, err := Parse([]byte(dagReverseYAML))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	d := Driver{DeckDir: t.TempDir(), Manifest: m, Now: now}
	run := NewPipelineRun(m, now)
	run.CurrentBlock = "" // let the DAG executor pick the root

	doneSet := map[string]bool{}
	complete := func(b Block) (bool, error) { return doneSet[b.ID], nil }

	order := []string{}
	for i := 0; i < 12; i++ {
		res, err := d.Advance(&run, complete)
		if err != nil {
			t.Fatalf("advance %d: %v", i, err)
		}
		switch res.Action {
		case ActionSeededNext:
			order = append(order, res.Block)
			doneSet[res.Block] = true // simulate the block finishing before next advance
		case ActionDone:
			want := []string{"spec", "api", "integrate"}
			if len(order) != 3 || order[0] != want[0] || order[1] != want[1] || order[2] != want[2] {
				t.Fatalf("activation order = %v, want %v (topological, not list order)", order, want)
			}
			return
		case ActionAwaitGate, ActionRejected:
			t.Fatalf("unexpected %s under auto-left low-risk", res.Action)
		}
	}
	t.Fatalf("never completed; order so far = %v", order)
}
