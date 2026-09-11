package app

import (
	"context"
	"os"
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
