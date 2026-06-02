package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/pipeline"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func autoManifest(autonomy string) string {
	return `schema_version: 1
idea_slug: auto-demo
transport: local-dir
autonomy: ` + autonomy + `
participants: [codex, claude]
blocks:
  - {id: b1, kind: deliberation, output_artifact: A.md}
  - {id: b2, kind: deliberation, output_artifact: B.md}
`
}

// startAndFinalize starts the pipeline and pre-finalizes both block workspaces
// so the auto-loop control flow can be exercised without launching any agent.
func startAndFinalize(t *testing.T, autonomy string) string {
	t.Helper()
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, "parley-deck"), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(ws, "pipeline.yaml")
	writeFile(t, manifest, autoManifest(autonomy))

	var out, errOut bytes.Buffer
	if code := Run([]string{"pipeline", "start", "--dir", ws, manifest}, &out, &errOut); code != 0 {
		t.Fatalf("start exit=%d err=%s", code, errOut.String())
	}
	final := "---\nstatus: final\n---\n\ndone\n"
	writeFile(t, filepath.Join(ws, "parley-deck", "ideas", "auto-demo__b1", "FINAL.md"), final)
	writeFile(t, filepath.Join(ws, "parley-deck", "ideas", "auto-demo__b2", "FINAL.md"), final)
	return ws
}

func TestPipelineAutoWalksToDoneUnderAutoLeft(t *testing.T) {
	ws := startAndFinalize(t, "auto-left")
	var out, errOut bytes.Buffer
	code := Run([]string{"pipeline", "auto", "--dir", ws, "--yes", "auto-demo"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("auto exit=%d err=%s out=%s", code, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "complete") {
		t.Fatalf("expected completion, got:\n%s", out.String())
	}
	// Block 2's kickoff must have been seeded by the driver as it advanced.
	if _, err := os.Stat(filepath.Join(ws, "parley-deck", "ideas", "auto-demo__b2", "00-prompt.md")); err != nil {
		t.Fatalf("next block was not seeded: %v", err)
	}
}

func TestBlockCompleteRespectsBlockedReviewConsensus(t *testing.T) {
	ws := t.TempDir()
	deck := filepath.Join(ws, "parley-deck")
	slug := "impl-demo"
	block := pipeline.Block{ID: "build", Kind: pipeline.KindImplementation}
	rcDir := filepath.Join(pipeline.BlockWorkspace(deck, slug, block.ID), "review")
	complete := blockCompleteFunc(deck, slug)

	// 0 outstanding fixes but BLOCKED -> NOT complete (fail closed).
	writeFile(t, filepath.Join(rcDir, "consensus.md"), "---\nidea: x\noutstanding_agreed_fixes: 0\nblocked: true\n---\n")
	if done, err := complete(block); err != nil || done {
		t.Fatalf("blocked review consensus must not be complete (done=%v err=%v)", done, err)
	}

	// 0 outstanding fixes, not blocked -> complete.
	writeFile(t, filepath.Join(rcDir, "consensus.md"), "---\nidea: x\noutstanding_agreed_fixes: 0\nblocked: false\n---\n")
	if done, err := complete(block); err != nil || !done {
		t.Fatalf("unblocked zero-fix consensus should be complete (done=%v err=%v)", done, err)
	}
}

func TestActionBlockCompleteNeedsSucceededEffectNotJustPlan(t *testing.T) {
	ws := t.TempDir()
	deck := filepath.Join(ws, "parley-deck")
	slug, blockID := "dep-demo", "deploy"
	block := pipeline.Block{ID: blockID, Kind: pipeline.KindAction, OutputArtifact: "DEPLOYMENT.md"}
	complete := blockCompleteFunc(deck, slug)

	// A finalized PLAN alone must NOT make an action block complete (otherwise
	// the DAG/auto could advance past an unexecuted deploy).
	bw := pipeline.BlockWorkspace(deck, slug, blockID)
	writeFile(t, filepath.Join(bw, "DEPLOYMENT.md"), "---\nstatus: final\n---\n\nplan\n")
	writeFile(t, filepath.Join(bw, "FINAL.md"), "---\nstatus: final\n---\n\nplan\n")
	if done, err := complete(block); err != nil || done {
		t.Fatalf("action with finalized plan but no effect must NOT be complete (done=%v err=%v)", done, err)
	}

	// Record a succeeded effect for the block -> now complete.
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	eff := pipeline.NewEffect(slug, blockID, "vercel", "deploy.production", "app", pipeline.HashRequest([]byte("{}")), pipeline.RiskProduction, now)
	eff.Advance(pipeline.EffectSucceeded, "dpl_1", "done", now)
	if err := pipeline.SaveEffect(deck, eff); err != nil {
		t.Fatal(err)
	}
	if done, err := complete(block); err != nil || !done {
		t.Fatalf("action with a succeeded effect should be complete (done=%v err=%v)", done, err)
	}
}

func TestPipelineAutoStopsAtActionBlockNeedsHumanGate(t *testing.T) {
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, "parley-deck"), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(ws, "pipeline.yaml")
	writeFile(t, manifest, `schema_version: 1
idea_slug: act-demo
transport: local-dir
autonomy: auto-left
participants: [codex, claude]
blocks:
  - {id: spec, kind: deliberation, output_artifact: SPEC.md}
  - {id: deploy, kind: action, output_artifact: DEPLOYMENT.md}
`)
	var out, errOut bytes.Buffer
	if code := Run([]string{"pipeline", "start", "--dir", ws, manifest}, &out, &errOut); code != 0 {
		t.Fatalf("start exit=%d err=%s", code, errOut.String())
	}
	final := "---\nstatus: final\n---\n\ndone\n"
	writeFile(t, filepath.Join(ws, "parley-deck", "ideas", "act-demo__spec", "FINAL.md"), final)
	writeFile(t, filepath.Join(ws, "parley-deck", "ideas", "act-demo__deploy", "FINAL.md"), final)

	out.Reset()
	code := Run([]string{"pipeline", "auto", "--dir", ws, "--yes", "act-demo"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("auto exit=%d err=%s out=%s", code, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "needs_human_gate") || !strings.Contains(out.String(), "action block") {
		t.Fatalf("auto should stop at action block with needs_human_gate, got:\n%s", out.String())
	}
	// It advanced past the deliberation block to the action block, but must NOT
	// have completed the pipeline or executed anything.
	var statusOut bytes.Buffer
	Run([]string{"pipeline", "status", "--dir", ws, "act-demo"}, &statusOut, &errOut)
	if strings.Contains(statusOut.String(), "status=completed") {
		t.Fatalf("pipeline must not complete past an unexecuted action block:\n%s", statusOut.String())
	}
	if !strings.Contains(statusOut.String(), "current=deploy") {
		t.Fatalf("expected current=deploy, got:\n%s", statusOut.String())
	}
}

func TestPipelineAutoPausesAtSupervisedGate(t *testing.T) {
	ws := startAndFinalize(t, "supervised")
	var out, errOut bytes.Buffer
	code := Run([]string{"pipeline", "auto", "--dir", ws, "--yes", "auto-demo"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("auto exit=%d err=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "paused at gate") {
		t.Fatalf("expected supervised pause, got:\n%s", out.String())
	}
	// Status should reflect the open boundary gate.
	var statusOut bytes.Buffer
	Run([]string{"pipeline", "status", "--dir", ws, "auto-demo"}, &statusOut, &errOut)
	if !strings.Contains(statusOut.String(), "blocked_on_gate") {
		t.Fatalf("status not blocked_on_gate:\n%s", statusOut.String())
	}
}
