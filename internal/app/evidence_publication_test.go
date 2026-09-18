package app

import (
	"context"
	"io"
	"os"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

func TestVerificationRequiresIgnoredRuntimeBeforeArtifacts(t *testing.T) {
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	for _, ignore := range []string{"", ".parley-runtime/evidence-verification/probe/\n"} {
		t.Run(strings.ReplaceAll(ignore, "/", "_"), func(t *testing.T) {
			root, dir := gateScratchRepo(t, twoCriterionContract())
			if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(ignore), 0644); err != nil {
				t.Fatal(err)
			}
			op := gateOps(root, dir)
			op.drafter = "reviewer"
			op.base = runner.Options{RunID: "runtime-prerequisite", Store: store.New(filepath.Join(root, ".parley-runtime/runs/prerequisite")), Agents: []agents.Discovery{{Spec: agents.Spec{ID: "reviewer"}, Found: true, Path: "/usr/bin/false"}}}
			if ok, detail := op.runChecksContract(context.Background(), mustContract(t, dir)); !ok {
				t.Fatal(detail)
			}
			ok, detail := op.VerifyCompletionEvidence(context.Background())
			if ok || !strings.Contains(detail, "runtime must be ignored before execution") {
				t.Fatalf("wrong prerequisite result: %t %s", ok, detail)
			}
			if _, err := os.Lstat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
				t.Fatalf("refusal created runtime artifacts: %v", err)
			}
		})
	}
}

func TestTripleDeletedScopeCannotClose(t *testing.T) {
	for _, throughDriver := range []bool{false, true} {
		name := "direct"
		if throughDriver {
			name = "advance"
		}
		t.Run(name, func(t *testing.T) {
			root, ideaDir := gateScratchRepo(t, twoCriterionContract())
			op := gateOps(root, ideaDir)
			if ok, detail := op.runChecksContract(context.Background(), mustContract(t, ideaDir)); !ok {
				t.Fatal(detail)
			}
			runDir := filepath.Join(root, ".parley-runtime", "prior-run")
			op.drafter = "reviewer"
			op.base = runner.Options{RunID: "prior-run", Store: store.New(runDir)}
			pin, err := driver.ObserveChecksContract(ideaDir, "")
			if err != nil {
				t.Fatal(err)
			}
			c := driver.Rebuild(ideaDir, 4)
			c.ChecksContractSHA256 = pin
			if err := c.Save(filepath.Join(runDir, "driver.json")); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("prior reviewed fixture"), 0644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte("---\nidea: idea-x\n---\n"), 0644); err != nil {
				t.Fatal(err)
			}
			for _, path := range []string{evidence.ReportPath(ideaDir), filepath.Join(runDir, "driver.json")} {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
			if throughDriver {
				d := driver.New(driver.Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "idea-x", RunDir: runDir, Participants: []string{"implementer", "reviewer"}, Auto: true, Events: store.New(runDir), Impl: evidenceCloseFixtureOps{driverImplOps: op}, Out: io.Discard}, nil)
				action, _, err := d.Advance(context.Background())
				if err == nil || action != driver.ActionEscalated {
					t.Fatalf("triple deletion advanced: %s %v", action, err)
				}
			} else if err := op.Complete(context.Background()); err == nil {
				t.Fatal("triple deletion bypassed all completion evidence")
			}
			data, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), "status: complete") {
				t.Fatal("deleted contract history closed without evidence")
			}
		})
	}
}
