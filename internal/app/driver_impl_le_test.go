package app

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

func hasEventType(evs []store.Event, typ string) bool {
	for _, e := range evs {
		if e.Type == typ {
			return true
		}
	}
	return false
}

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

// CF1: newDriverImplOps dedupes to distinct non-implementer reviewer IDs, so duplicate
// participant entries (e.g. [impl, rev, rev]) cannot inflate ReviewerCount past the LE-11
// `< 2` guard nor launch duplicate review goroutines onto the same artifact files.
func TestNewDriverImplOpsDedupesReviewers(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "") // no role metadata → implementer = participants[0]
	ops := newDriverImplOps(runner.Options{}, root, "demo", ideaDir, []string{"claude", "codex", "codex"}, io.Discard)
	o, ok := ops.(driverImplOps)
	if !ok {
		t.Fatalf("expected driverImplOps, got %T", ops)
	}
	if len(o.reviewers) != 1 || o.reviewers[0] != "codex" {
		t.Fatalf("duplicate reviewer IDs must collapse to one; got %v", o.reviewers)
	}
	// Two distinct reviewers + a duplicate → two reviewers, order preserved.
	ops2 := newDriverImplOps(runner.Options{}, root, "demo", ideaDir, []string{"claude", "codex", "agy", "codex"}, io.Discard)
	o2 := ops2.(driverImplOps)
	if got := strings.Join(o2.reviewers, ","); got != "codex,agy" {
		t.Fatalf("expected distinct reviewers [codex,agy], got %v", o2.reviewers)
	}
}

// CF6: GoalCheck fails open (advisory) without running any agent when the only available
// checker is the implementer itself (drafter == implementer) — it never runs the
// implementer as its own goal checker.
func TestGoalCheckNoIndependentChecker(t *testing.T) {
	o := newOpsFor(t.TempDir(), t.TempDir(), nil, "claude", nil)
	o.drafter = "claude" // == implementer
	ok, detail := o.GoalCheck(context.Background())
	if !ok {
		t.Fatalf("goal-check with no independent checker must fail open; got (%v, %q)", ok, detail)
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

	// F6: case-insensitive — "m1" vs "M1" are the same model.
	casing := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "M1"}},
	}
	o4 := newOpsFor(t.TempDir(), t.TempDir(), casing, "claude", []string{"codex"})
	if shared, _ := o4.reviewersShareImplementerModel(); !shared {
		t.Fatal("case-only model differences must still count as shared (F6)")
	}
}

func newOpsWithStore(root, ideaDir, runDir string, agentsList []agents.Discovery, out io.Writer) driverImplOps {
	return driverImplOps{
		base: runner.Options{Agents: agentsList, Store: store.New(runDir)}, root: root,
		ideaSlug: "demo", ideaDir: ideaDir, implementer: "claude", reviewers: []string{"codex"}, out: out,
	}
}

// F4/F8: the warn path emits a stdout warning AND an agent.model_diversity event.
func TestCheckModelDiversityWarnsAndEmitsEvent(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "") // no require_model_diversity → warn
	runDir := t.TempDir()
	var out bytes.Buffer
	same := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "m1"}},
	}
	o := newOpsWithStore(root, ideaDir, runDir, same, &out)
	if err := o.checkModelDiversity(); err != nil {
		t.Fatalf("warn path must not error: %v", err)
	}
	if !strings.Contains(out.String(), "WARNING model-diversity") {
		t.Fatalf("expected a warning, got: %q", out.String())
	}
	evs, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if !hasEventType(evs, "agent.model_diversity") {
		t.Fatal("expected an agent.model_diversity event in the store")
	}
}

// F8: require_model_diversity escalates (returns an error) on a same-model roster.
func TestCheckModelDiversityEscalatesWhenRequired(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "require_model_diversity: true\n")
	same := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "m1"}},
	}
	o := newOpsWithStore(root, ideaDir, t.TempDir(), same, io.Discard)
	if err := o.checkModelDiversity(); err == nil {
		t.Fatal("require_model_diversity: true must escalate on a same-model roster")
	}
}

// review-01 F4: track: fast makes model diversity a HARD gate even without
// require_model_diversity — a same-model single reviewer on fast must escalate.
func TestCheckModelDiversityHardGateOnFastTrack(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "track: fast\n") // no require_model_diversity, but fast forces it
	same := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "m1"}},
	}
	o := newOpsWithStore(root, ideaDir, t.TempDir(), same, io.Discard)
	if err := o.checkModelDiversity(); err == nil {
		t.Fatal("track: fast must escalate on a same-model roster (hard gate) even without require_model_diversity")
	}
}

// F8: a diverse roster is silent — no warning, no error.
func TestCheckModelDiversitySilentWhenDiverse(t *testing.T) {
	root := t.TempDir()
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, ideaDir, "")
	var out bytes.Buffer
	diff := []agents.Discovery{
		{Spec: agents.Spec{ID: "claude", Model: "m1"}},
		{Spec: agents.Spec{ID: "codex", Model: "m2"}},
	}
	o := newOpsWithStore(root, ideaDir, t.TempDir(), diff, &out)
	if err := o.checkModelDiversity(); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("a diverse roster must be silent, got: %q", out.String())
	}
}
