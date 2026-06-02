package pipeline

import (
	"sort"
	"testing"
	"time"
)

const diamondYAML = `
schema_version: 1
idea_slug: diamond
transport: local-dir
execution: dag
autonomy: auto-left
participants: [codex, claude]
blocks:
  - {id: spec, kind: deliberation}
  - {id: api, kind: deliberation}
  - {id: ui, kind: deliberation}
  - {id: integrate, kind: deliberation}
edges:
  - {from: spec, to: api}
  - {from: spec, to: ui}
  - {from: api, to: integrate}
  - {from: ui, to: integrate}
`

func TestComputeDAGStepParallelWaves(t *testing.T) {
	m, err := Parse([]byte(diamondYAML))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	d := Driver{DeckDir: t.TempDir(), Manifest: m, Now: now}

	step := func(completed ...string) DAGStep {
		run := PipelineRun{PipelineSlug: "diamond", CompletedBlocks: completed}
		s, err := d.ComputeDAGStep(&run)
		if err != nil {
			t.Fatalf("ComputeDAGStep: %v", err)
		}
		sort.Strings(s.Ready)
		return s
	}

	// Root only.
	if s := step(); len(s.Ready) != 1 || s.Ready[0] != "spec" {
		t.Fatalf("wave 0 ready = %v, want [spec]", s.Ready)
	}
	// After spec: api AND ui both ready in one wave (parallel fan-out).
	if s := step("spec"); len(s.Ready) != 2 || s.Ready[0] != "api" || s.Ready[1] != "ui" {
		t.Fatalf("wave 1 ready = %v, want [api ui]", s.Ready)
	}
	// Only api done: integrate NOT yet ready (join needs ui too).
	if s := step("spec", "api"); contains(s.Ready, "integrate") {
		t.Fatalf("integrate must not be ready before ui: %v", s.Ready)
	}
	// api+ui done: integrate ready.
	if s := step("spec", "api", "ui"); len(s.Ready) != 1 || s.Ready[0] != "integrate" {
		t.Fatalf("join wave ready = %v, want [integrate]", s.Ready)
	}
	// All done.
	if s := step("spec", "api", "ui", "integrate"); !s.Done {
		t.Fatalf("expected Done")
	}
}

func TestComputeDAGStepSupervisedAwaitsGate(t *testing.T) {
	m, _ := Parse([]byte(diamondYAML))
	m.Autonomy = AutonomySupervised // no auto-approval
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	d := Driver{DeckDir: t.TempDir(), Manifest: m, Now: now}
	run := PipelineRun{PipelineSlug: "diamond"}
	s, err := d.ComputeDAGStep(&run)
	if err != nil {
		t.Fatalf("ComputeDAGStep: %v", err)
	}
	if len(s.Ready) != 0 || len(s.AwaitGates) != 1 || s.AwaitGates[0] != EdgeID("ready", "spec") {
		t.Fatalf("supervised: ready=%v await=%v, want await [ready->spec]", s.Ready, s.AwaitGates)
	}
}

func contains(xs []string, x string) bool {
	for _, s := range xs {
		if s == x {
			return true
		}
	}
	return false
}
