package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/runner"
)

func writePrompt(t *testing.T, ideaDir, frontmatter string) {
	t.Helper()
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nidea: demo\n" + frontmatter + "---\n\n## Problem\nx\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newOpsFor(root, ideaDir string, agentsList []agents.Discovery, implementer string, reviewers []string) driverImplOps {
	return driverImplOps{
		base: runner.Options{Agents: agentsList}, root: root, ideaSlug: "demo", ideaDir: ideaDir,
		implementer: implementer, reviewers: reviewers, drafter: "", out: io.Discard,
	}
}

// LE-4: an explicit checks: command is honored (pass + fail by exit code).
func TestRunChecksHonorsChecksCommand(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "checks: \"exit 0\"\n")
	o := newOpsFor(root, ideaDir, nil, "claude", []string{"codex"})
	if ok, detail := o.RunChecks(context.Background()); !ok {
		t.Fatalf("checks 'exit 0' should pass; detail=%s", detail)
	}
	writePrompt(t, ideaDir, "checks: \"exit 7\"\n")
	if ok, _ := o.RunChecks(context.Background()); ok {
		t.Fatal("checks 'exit 7' should fail")
	}
}

// LE-4: a code-writing (auto_implement) idea with no go.mod and no checks fails closed.
func TestRunChecksFailsClosedForCodeIdeaWithoutChecks(t *testing.T) {
	root := t.TempDir() // no go.mod
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "auto_implement: true\n")
	o := newOpsFor(root, ideaDir, nil, "claude", []string{"codex"})
	if ok, detail := o.RunChecks(context.Background()); ok {
		t.Fatalf("a code-writing idea with no checks must fail closed; got pass (%s)", detail)
	}
}

// LE-4: a design-only idea (no auto_implement, no go.mod, no checks) still passes.
func TestRunChecksDesignOnlyPasses(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "")
	o := newOpsFor(root, ideaDir, nil, "claude", []string{"codex"})
	if ok, _ := o.RunChecks(context.Background()); !ok {
		t.Fatal("a design-only idea with nothing to check must pass")
	}
}

// LE-3: model-diversity detection — all reviewers same model as implementer vs differing.
func TestReviewersShareImplementerModel(t *testing.T) {
	same := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "m1"}},
		{Spec: agents.Spec{ID: "agy", Model: "m1"}},
	}
	o := newOpsFor(t.TempDir(), t.TempDir(), same, "claude", []string{"codex", "agy"})
	if shared, model := o.reviewersShareImplementerModel(); !shared || model != "m1" {
		t.Fatalf("expected shared=true model=m1, got shared=%v model=%q", shared, model)
	}

	diff := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "m2"}},
	}
	o2 := newOpsFor(t.TempDir(), t.TempDir(), diff, "claude", []string{"codex"})
	if shared, _ := o2.reviewersShareImplementerModel(); shared {
		t.Fatal("differing models must not be flagged as shared")
	}

	// Unknown implementer model → never fires.
	o3 := newOpsFor(t.TempDir(), t.TempDir(), nil, "claude", []string{"codex"})
	if shared, _ := o3.reviewersShareImplementerModel(); shared {
		t.Fatal("unknown implementer model must not fire the guard")
	}
}
