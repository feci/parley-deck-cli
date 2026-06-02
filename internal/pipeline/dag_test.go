package pipeline

import "testing"

const dagYAML = `
schema_version: 1
idea_slug: dag-demo
transport: local-dir
execution: dag
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

func TestDAGManifestParses(t *testing.T) {
	m, err := Parse([]byte(dagYAML))
	if err != nil {
		t.Fatalf("valid DAG rejected: %v", err)
	}
	if !m.IsDAG() {
		t.Fatal("expected IsDAG")
	}
}

func TestDAGRejectsCycle(t *testing.T) {
	y := `
schema_version: 1
idea_slug: cyc
transport: local-dir
execution: dag
participants: [a, b]
blocks:
  - {id: x, kind: deliberation}
  - {id: y, kind: deliberation}
edges:
  - {from: x, to: y}
  - {from: y, to: x}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected cycle to be rejected")
	}
}

func TestDAGRejectsUnknownEndpoint(t *testing.T) {
	y := `
schema_version: 1
idea_slug: bad
transport: local-dir
execution: dag
participants: [a, b]
blocks:
  - {id: x, kind: deliberation}
edges:
  - {from: x, to: ghost}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected unknown endpoint to be rejected")
	}
}

func TestReadyBlocksFanOutAndJoin(t *testing.T) {
	m, err := Parse([]byte(dagYAML))
	if err != nil {
		t.Fatal(err)
	}
	allGatesOK := func(string) bool { return true }

	// Nothing done: only the root "spec" is ready.
	ready := m.ReadyBlocks(nil, allGatesOK)
	if len(ready) != 1 || ready[0] != "spec" {
		t.Fatalf("initial ready = %v, want [spec]", ready)
	}
	// spec done: api and ui both become ready (fan-out).
	ready = m.ReadyBlocks([]string{"spec"}, allGatesOK)
	if len(ready) != 2 {
		t.Fatalf("after spec ready = %v, want [api ui]", ready)
	}
	// api+ui done: integrate ready (join).
	ready = m.ReadyBlocks([]string{"spec", "api", "ui"}, allGatesOK)
	if len(ready) != 1 || ready[0] != "integrate" {
		t.Fatalf("join ready = %v, want [integrate]", ready)
	}
	// Only api done: integrate NOT ready (ui still pending).
	ready = m.ReadyBlocks([]string{"spec", "api"}, allGatesOK)
	for _, r := range ready {
		if r == "integrate" {
			t.Fatal("integrate must not be ready until both api and ui complete")
		}
	}
	if m.AllBlocksComplete([]string{"spec", "api", "ui"}) {
		t.Fatal("not all complete yet")
	}
	if !m.AllBlocksComplete([]string{"spec", "api", "ui", "integrate"}) {
		t.Fatal("should be all complete")
	}
}

func TestReadyBlocksRespectsUnapprovedGate(t *testing.T) {
	m, _ := Parse([]byte(dagYAML))
	// spec done but the spec->api gate not approved: api blocked, ui (gate ok) ready.
	ready := m.ReadyBlocks([]string{"spec"}, func(edge string) bool {
		return edge != EdgeID("spec", "api")
	})
	for _, r := range ready {
		if r == "api" {
			t.Fatal("api must be blocked by its unapproved gate")
		}
	}
}

func TestPerBlockTransportOverride(t *testing.T) {
	y := `
schema_version: 1
idea_slug: t
transport: local-dir
participants: [a, b]
blocks:
  - {id: b1, kind: deliberation}
  - {id: b2, kind: deliberation, transport: github-pr}
`
	m, err := Parse([]byte(y))
	if err != nil {
		t.Fatalf("per-block transport rejected: %v", err)
	}
	if m.EffectiveTransport(m.Blocks[0]) != "local-dir" {
		t.Fatalf("b1 transport = %q", m.EffectiveTransport(m.Blocks[0]))
	}
	if m.EffectiveTransport(m.Blocks[1]) != "github-pr" {
		t.Fatalf("b2 transport = %q", m.EffectiveTransport(m.Blocks[1]))
	}
}

func TestPerBlockTransportValidated(t *testing.T) {
	y := `
schema_version: 1
idea_slug: t
transport: local-dir
participants: [a, b]
blocks:
  - {id: b1, kind: deliberation, transport: smoke-signals}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("invalid per-block transport must be rejected")
	}
}

func TestDeciderAutoApprovesLowRiskOnly(t *testing.T) {
	// With a decider configured (hasDecider=true) under supervised autonomy,
	// ONLY low-risk gates auto-resolve (strictly narrower than auto-left).
	if !AutoApproveWithDecider(AutonomySupervised, RiskLow, true) {
		t.Fatal("decider should auto-approve low-risk")
	}
	if AutoApproveWithDecider(AutonomySupervised, RiskNormal, true) {
		t.Fatal("decider must NOT auto-approve normal-risk (low-risk only)")
	}
	if AutoApproveWithDecider(AutonomySupervised, RiskHigh, true) {
		t.Fatal("decider must NOT auto-approve high-risk")
	}
	// auto-left keeps the broader low/normal behavior.
	if !AutoApproveWithDecider(AutonomyAutoLeft, RiskNormal, false) {
		t.Fatal("auto-left should still auto-approve normal-risk")
	}
	if AutoApproveWithDecider(AutonomySupervised, RiskProduction, true) {
		t.Fatal("decider must NEVER auto-approve production")
	}
	// No decider, supervised: block-and-wait.
	if AutoApproveWithDecider(AutonomySupervised, RiskLow, false) {
		t.Fatal("no decider + supervised must block")
	}
}
