package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
