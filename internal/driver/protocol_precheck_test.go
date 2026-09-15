package driver

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

func TestDriverProtocolPrecheckPrecedesFixupCharges(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX fixture process")
	}
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	ctx := context.Background()
	parts := []string{"builder", "reviewer"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\ntrack: deliberation\n")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	writeFinalValid(t, ideaDir)
	writeImplWithCycles(t, ideaDir, "implemented", 0)
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	for n, body := range map[string]string{".gitignore": ".parley-runtime/\nparley-deck/runs/\n", "source": "original\n", "parley-deck/meta/version.json": `{"protocolRole":"source"}`} {
		p := filepath.Join(root, n)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"commit", "-qm", "N1 driver fixture"}} {
		argv := append([]string{"-C", root, "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid"}, args...)
		if raw, err := exec.Command("git", argv...).CombinedOutput(); err != nil {
			t.Fatal(err, string(raw))
		}
	}
	b, err := budget.EnsureCycleBinding(ctx, root, "demo", budget.Fixup, 5, 0, runDir, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	policy, expected, err := trajectory.NewPolicy(ctx, root, "demo", "builder", []trajectory.Criterion{{Name: "material", Command: "true"}})
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Activate(ctx, root, expected, policy); err != nil {
		t.Fatal(err)
	}
	before, err := trajectory.Inspect(ctx, root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	beforeLedger, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", PromptMode: agents.PromptStdin, HeadlessArgs: []string{"-c", "touch .parley-runtime/n1-spawned; exit 7"}}, Found: true, Path: "/bin/sh"}
	opts := runner.Options{Root: root, RunID: "run1", Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: []string{"builder"}}, Agents: []agents.Discovery{agent}, Timeout: 5 * time.Second, Store: store.New(runDir)}
	fi := &cycleProcessImpl{fakeImpl: &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}, opts: opts}
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 5, MaxDriverSteps: 10, Impl: fi}, &fakeRunner{})
	if _, err = protocolpacket.Render(root, protocolcore.StoreAt(config.CentralHome()), protocolpacket.Request{Phase: 8, Track: "deliberation", IdeaSlug: "demo"}); err != nil {
		t.Fatal(err)
	}
	paths, err := filepath.Glob(filepath.Join(protocolpacket.RuntimeDir(root), "*.md"))
	if err != nil || len(paths) != 1 {
		t.Fatal(paths, err)
	}
	original, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(paths[0], []byte("tampered cached context"), 0600); err != nil {
		t.Fatal(err)
	}
	action, cursor, err := d.Advance(ctx)
	if err == nil || !strings.Contains(err.Error(), "protocol context refused") || action != ActionEscalated || cursor.FixupCyclesPublished != 0 || contains(fi.calls, "fixup") {
		t.Fatal("known refusal reached charged fixup", action, cursor, fi.calls, err)
	}
	after, err := trajectory.Inspect(ctx, root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	encode := func(v any) string {
		raw, e := json.Marshal(v)
		if e != nil {
			t.Fatal(e)
		}
		return string(raw)
	}
	if encode(before) != encode(after) || encode(beforeLedger) != encode(ledger) {
		t.Fatal("known refusal changed original trajectory or ledger")
	}
	step, err := budget.LoadStepBinding(ctx, root, "demo")
	if err != nil || step == nil {
		t.Fatal(step, err)
	}
	steps, err := step.Store.Inspect(ctx)
	if err != nil || len(steps.Entries) != 0 {
		t.Fatal("protocol precheck spent a driver step", steps, err)
	}
	terminals, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil || len(terminals) != 1 {
		t.Fatal("missing actual refusal terminal", terminals, err)
	}
	raw, err := os.ReadFile(terminals[0])
	if err != nil {
		t.Fatal(err)
	}
	var terminal telemetry.Record
	if err = json.Unmarshal(raw, &terminal); err != nil {
		t.Fatal(err)
	}
	if terminal.Metadata.Agent != "builder" || terminal.Metadata.RunID != "run1" || terminal.Metadata.Idea != "demo" || terminal.Metadata.Phase != "fixup" || terminal.StartedAt != nil || terminal.PID != nil || terminal.Outcome == nil || terminal.Outcome.FailureClass == nil || *terminal.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatal("wrong precheck refusal provenance", terminal)
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "n1-spawned")); !os.IsNotExist(err) {
		t.Fatal("known refusal spawned", err)
	}
	if err = os.WriteFile(paths[0], original, 0600); err != nil {
		t.Fatal(err)
	}
	_, cursor, err = d.Advance(ctx)
	if err == nil || !contains(fi.calls, "fixup") || cursor.FixupCyclesPublished != 1 {
		t.Fatal("restoration lost normal charged driver fixup", cursor, fi.calls, err)
	}
	after, err = trajectory.Inspect(ctx, root, "demo")
	if err != nil || len(after.Attempts) != 1 || after.Attempts[0].Launch == nil || after.Attempts[0].Terminal == nil {
		t.Fatal("restored driver lost actual charged process", after, err)
	}
	ledger, err = b.Store.Inspect(ctx)
	if err != nil || b.Count(ledger) != 1 {
		t.Fatal("restored driver did not charge exactly once", ledger, err)
	}
	steps, err = step.Store.Inspect(ctx)
	if err != nil || len(steps.Entries) != 1 {
		t.Fatal("restored driver step charge missing/duplicated", steps, err)
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "n1-spawned")); err != nil {
		t.Fatal("restored driver process did not execute", err)
	}
	current, err := os.ReadFile(terminals[0])
	if err != nil || string(current) != string(raw) {
		t.Fatal("restoration changed refusal history", err)
	}
}
